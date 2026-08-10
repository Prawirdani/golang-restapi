---
name: gorest-handlers
description: "HTTP handler conventions for github.com/prawirdani/golang-restapi — httpx.Context handler signature, httpx.Handler wrapping, BindValidate, response envelope httpx.Body, status codes, cookies, multipart handling, and ETag behavior. Use when writing or reviewing any handler in internal/transport/http/handler/."
user-invocable: true
license: MIT
compatibility: Designed for AI coding agents working in the golang-restapi repository.
metadata:
  author: prawirdani
  version: "1.0.0"
  module: github.com/prawirdani/golang-restapi
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Agent
---

**Persona:** You are a strict HTTP layer reviewer. Handlers are thin: parse → validate → call service → return envelope. Any business logic, retry logic, or DB access in a handler is a defect. Handlers never write error responses themselves.

## Rules

1. **Signature is fixed:** `func (h *<Name>Handler) <Op>(c *httpx.Context) error`. Register with `httpx.Handler(fn)` in `cmd/api/server.go`. Middleware is applied via `httpx.Middleware(...)` or `chi` route groups — never inside the handler body.

2. **Flow order is:**
   ```go
   ctx := c.Context()          // 1. snapshot request context
   claims, err := auth.GetAccessTokenCtx(ctx) // 2. auth claims when needed
   var reqBody user.UpdateUserInput
   if err := c.BindValidate(&reqBody); err != nil { // 3. bind + validate
       return err
   }
   if err := h.userService.UpdateUser(ctx, claims.UserID, reqBody); err != nil { // 4. delegate
       return err
   }
   return c.JSON(&httpx.Body{Message: "user updated!"}) // 5. envelope
   ```

3. **Bind with `c.BindValidate`** (enforces `Content-Type: application/json`, decodes, runs `pkg/validator`). Do NOT use plain `c.Bind` for request bodies, and never write manual struct validation in handlers.

4. **Respond only through `httpx.Body`** — `&httpx.Body{Data: ..., Message: ...}` (pointer; `MarshalJSON` keeps `message` present). Success codes: `c.Status(http.StatusCreated).JSON(...)` for 201, plain `c.JSON(...)` for 200.

5. **Never write error responses manually.** Return the error; `httpx.Handler` normalizes it (`internal/transport/http/context.go:216-226`). Handlers may `log.ErrorCtx(ctx, "...", err)` before returning for observability — but do not log-and-return in tests of pure paths without reason.

6. **Auth claims:** read via `auth.GetAccessTokenCtx(ctx)`; user identity comes from claims (UserID, SessionID), never from the request body.

7. **Multipart:** `c.EnsureMultipartForm()` → `defer c.CleanupMultipart()` → `c.FormFile(httpx.ImageFormKey)` → `httpx.NewParsedFile(fh)` + `defer file.Close()` → `httpx.ValidateFile(...)` with explicit `ValidationRules` (max size + allowed MIME).

8. **ETag is automatic** for GET/HEAD 2xx responses (`c.JSON` derives a strong ETag, honors `If-None-Match` with 304). Don't set `ETag`/`Cache-Control` manually.

## Checklist (Review mode)

- [ ] Handler is one parse/validate/delegate/respond pass — no business logic, no DB calls
- [ ] `BindValidate` used for JSON bodies; no manual validation
- [ ] Responses use `&httpx.Body{...}`; errors returned, never `http.Error`/`c.w.Write`
- [ ] Claims-based identity; `log.ErrorCtx` used where debugging matters
- [ ] Multipart lifecycle: EnsureMultipartForm → CleanupMultipart → ValidateFile
- [ ] Routes registered in `server.go` with `httpx.Handler(fn)`

## References

- `internal/transport/http/context.go` — Context API (JSON/Bind/BindValidate/Status/FormFile)
- `internal/transport/http/handler/auth_handler.go` — full flow incl. cookies (`setTokenCookies`)
- `internal/transport/http/handler/user_handler.go` — update + multipart flows
- `cmd/api/server.go` — route registration and middleware composition
