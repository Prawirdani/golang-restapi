package handler

import (
	"net/http"
	"time"

	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/internal/model"
	"github.com/prawirdani/golang-restapi/internal/service"
	"github.com/prawirdani/golang-restapi/pkg/common"
	"github.com/prawirdani/golang-restapi/pkg/httputil"
	"github.com/prawirdani/golang-restapi/pkg/validator"
)

type AuthHandler struct {
	authUC *service.AuthService
	cfg    *config.Config
}

func NewAuthHandler(cfg *config.Config, us *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authUC: us,
		cfg:    cfg,
	}
}

func (h *AuthHandler) HandleRegister(w http.ResponseWriter, r *http.Request) error {
	var reqBody model.RegisterRequest

	if err := httputil.BindJSON(r, &reqBody); err != nil {
		return err
	}

	if err := validator.Struct(reqBody); err != nil {
		return err
	}

	if err := h.authUC.Register(r.Context(), reqBody); err != nil {
		return err
	}

	return response(w, status(201), message("Registration successful."))
}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) error {
	var reqBody model.LoginRequest
	if err := httputil.BindJSON(r, &reqBody); err != nil {
		return err
	}

	if err := validator.Struct(reqBody); err != nil {
		return err
	}

	at, rt, err := h.authUC.Login(r.Context(), reqBody)
	if err != nil {
		return err
	}

	d := map[string]string{
		"accessToken":  at,
		"refreshToken": rt,
	}

	h.setTokenCookies(w, common.AccessToken, at)
	h.setTokenCookies(w, common.RefreshToken, rt)

	return response(w, data(d), message("Login successful."))
}

func (h *AuthHandler) CurrentUser(w http.ResponseWriter, r *http.Request) error {
	payload, err := httputil.GetAuthCtx[model.AccessTokenPayload](r.Context())
	if err != nil {
		return err
	}

	return response(w, data(payload.User))
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) error {
	payload, err := httputil.GetAuthCtx[model.RefreshTokenPayload](r.Context())
	if err != nil {
		return err
	}

	newAccessToken, err := h.authUC.RefreshToken(r.Context(), payload.User.ID)
	if err != nil {
		return err
	}

	d := map[string]string{
		"accessToken": newAccessToken,
	}

	h.setTokenCookies(w, common.AccessToken, newAccessToken)

	return response(w, data(d), message("Token refreshed."))
}

func (h *AuthHandler) setTokenCookies(w http.ResponseWriter, tokenType common.TokenType, tokenString string) {
	expiry := time.Now().Add(h.cfg.Token.AccessTokenExpiry)
	if tokenType == common.RefreshToken {
		expiry = time.Now().Add(h.cfg.Token.RefreshTokenExpiry)
	}

	ck := &http.Cookie{
		Name:     tokenType.Label(),
		Value:    tokenString,
		Expires:  expiry,
		HttpOnly: h.cfg.IsProduction(),
		Secure:   h.cfg.IsProduction(),
	}
	http.SetCookie(w, ck)
}
