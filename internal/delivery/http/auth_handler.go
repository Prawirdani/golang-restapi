package http

import (
	"net/http"
	"time"

	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/internal/model"
	"github.com/prawirdani/golang-restapi/internal/usecase"
	"github.com/prawirdani/golang-restapi/pkg/httputil"
	"github.com/prawirdani/golang-restapi/pkg/with"
)

type AuthHandler struct {
	userUC usecase.AuthUseCase
	cfg    *config.Config
}

func NewAuthHandler(cfg *config.Config, us usecase.AuthUseCase) AuthHandler {
	return AuthHandler{
		userUC: us,
		cfg:    cfg,
	}
}

func (h AuthHandler) HandleRegister(w http.ResponseWriter, r *http.Request) error {
	var reqBody model.RegisterRequest

	if err := httputil.BindJSON(r, &reqBody); err != nil {
		return err
	}

	if err := reqBody.ValidateRequest(); err != nil {
		return err
	}

	if err := h.userUC.Register(r.Context(), reqBody); err != nil {
		return err
	}

	return httputil.Response(w, with.Status(201), with.Message("Registration successful."))
}

func (h AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) error {
	var reqBody model.LoginRequest
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

	tokenCookie := &http.Cookie{
		Name:     h.cfg.Token.AccessCookieName,
		Value:    tokenString,
		Path:     "/",
		Expires:  time.Now().Add(time.Duration(h.cfg.Token.Expiry * int(time.Hour))),
		HttpOnly: h.cfg.IsProduction(),
	}

	http.SetCookie(w, tokenCookie)

	data := map[string]string{
		"token": tokenString,
	}
	return httputil.Response(w, with.Data(data), with.Message("Login successful."))
}

func (h AuthHandler) CurrentUser(w http.ResponseWriter, r *http.Request) error {
	tokenClaims := httputil.GetAuthCtx(r.Context())

	return httputil.Response(w, with.Data(tokenClaims))
}
