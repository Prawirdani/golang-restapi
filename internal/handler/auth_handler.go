package handler

import (
	"net/http"
	"time"

	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/internal/auth"
	"github.com/prawirdani/golang-restapi/internal/model"
	"github.com/prawirdani/golang-restapi/internal/service"
	req "github.com/prawirdani/golang-restapi/pkg/request"
	res "github.com/prawirdani/golang-restapi/pkg/response"
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
	if err := req.BindValidate(r, &reqBody); err != nil {
		return err
	}

	if err := h.authUC.Register(r.Context(), reqBody); err != nil {
		return err
	}

	return res.Send(w, res.WithStatus(201), res.WithMessage("Registration successful."))
}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) error {
	var reqBody model.LoginRequest
	if err := req.BindValidate(r, &reqBody); err != nil {
		return err
	}

	at, rt, err := h.authUC.Login(r.Context(), reqBody)
	if err != nil {
		return err
	}

	d := map[string]string{
		auth.AccessToken.Label():  at,
		auth.RefreshToken.Label(): rt,
	}

	h.setTokenCookies(w, auth.AccessToken, at)
	h.setTokenCookies(w, auth.RefreshToken, rt)

	return res.Send(w, res.WithData(d), res.WithMessage("Login successful."))
}

func (h *AuthHandler) CurrentUser(w http.ResponseWriter, r *http.Request) error {
	user, err := h.authUC.IdentifyUser(r.Context())
	if err != nil {
		return err
	}

	return res.Send(w, res.WithData(user))
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) error {
	newAccessToken, err := h.authUC.RefreshToken(r.Context())
	if err != nil {
		return err
	}

	d := map[string]string{
		auth.AccessToken.Label(): newAccessToken,
	}

	h.setTokenCookies(w, auth.AccessToken, newAccessToken)

	return res.Send(w, res.WithData(d), res.WithMessage("Token refreshed."))
}

func (h *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) error {
	atCookie := &http.Cookie{
		Name:     auth.AccessToken.Label(),
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: h.cfg.IsProduction(),
		Secure:   h.cfg.IsProduction(),
		Path:     "/",
	}

	rtCookie := *atCookie
	rtCookie.Name = auth.RefreshToken.Label()

	http.SetCookie(w, atCookie)
	http.SetCookie(w, &rtCookie)

	return res.Send(w, res.WithMessage("Logout successful."))
}

func (h *AuthHandler) setTokenCookies(
	w http.ResponseWriter,
	tokenType auth.TokenType,
	tokenString string,
) {
	currTime := time.Now()

	expiry := currTime.Add(h.cfg.Token.AccessTokenExpiry)
	if tokenType == auth.RefreshToken {
		expiry = currTime.Add(h.cfg.Token.RefreshTokenExpiry)
	}

	ck := &http.Cookie{
		Name:     tokenType.Label(),
		Value:    tokenString,
		Expires:  expiry,
		HttpOnly: h.cfg.IsProduction(),
		Secure:   h.cfg.IsProduction(),
		Path:     "/",
	}
	http.SetCookie(w, ck)
}
