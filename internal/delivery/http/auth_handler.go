package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prawirdani/golang-restapi/internal/middleware"
	"github.com/prawirdani/golang-restapi/internal/model"
	"github.com/prawirdani/golang-restapi/internal/usecase"
	"github.com/prawirdani/golang-restapi/pkg/httputil"
)

type AuthHandler struct {
	middleware middleware.MiddlewareManager
	userUC     usecase.AuthUseCase
}

func NewAuthHandler(mw middleware.MiddlewareManager, us usecase.AuthUseCase) AuthHandler {
	return AuthHandler{
		userUC:     us,
		middleware: mw,
	}
}

func (h AuthHandler) Routes(r chi.Router) {
	handlerFn := httputil.HandlerWrapper
	r.Post("/register", handlerFn(h.HandleRegister))
	r.Post("/login", handlerFn(h.HandleLogin))
	r.With(h.middleware.Authenticate).Get("/current", handlerFn(h.CurrentUser))
}

func (h AuthHandler) URLPattern() string {
	return "/auth"
}

func (h AuthHandler) HandleRegister(w http.ResponseWriter, r *http.Request) error {
	var reqBody model.RegisterRequestPayload

	if err := httputil.BindJSON(r, &reqBody); err != nil {
		return err
	}

	if err := reqBody.ValidateRequest(); err != nil {
		return err
	}

	if err := h.userUC.Register(r.Context(), reqBody); err != nil {
		return err
	}
	return httputil.SendJSON(w, http.StatusCreated, nil)
}

func (h AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) error {
	var reqBody model.LoginRequestPayload
	if err := httputil.BindJSON(r, &reqBody); err != nil {
		return err
	}

	if err := reqBody.ValidateRequest(); err != nil {
		return err
	}

	tokenString, err := h.userUC.Login(r.Context(), reqBody)
	if err != nil {
		return err
	}

	response := model.TokenResponse{
		Token: tokenString,
	}

	return httputil.SendJSON(w, http.StatusOK, response)
}

func (h AuthHandler) CurrentUser(w http.ResponseWriter, r *http.Request) error {
	tokenClaims := h.middleware.GetAuthCtx(r.Context())
	response := model.TokenInfoResponse{
		TokenInfo: tokenClaims,
	}
	return httputil.SendJSON(w, http.StatusOK, response)
}
