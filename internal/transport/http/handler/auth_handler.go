package handler

import (
	"errors"
	"net/http"
	"strings"
	"time"

	req "github.com/prawirdani/golang-restapi/internal/transport/http/request"
	res "github.com/prawirdani/golang-restapi/internal/transport/http/response"

	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/internal/auth"
	"github.com/prawirdani/golang-restapi/internal/model"
	"github.com/prawirdani/golang-restapi/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
	cfg         *config.Config
}

func NewAuthHandler(cfg *config.Config, us *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: us,
		cfg:         cfg,
	}
}

func (h *AuthHandler) HandleRegister(w http.ResponseWriter, r *http.Request) error {
	var reqBody model.CreateUserInput
	if err := req.BindValidate(r, &reqBody); err != nil {
		return err
	}

	if err := h.authService.Register(r.Context(), reqBody); err != nil {
		return err
	}

	return res.Send(
		w,
		r,
		res.WithStatus(201),
		res.WithMessage("Registration successful."),
	)
}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) error {
	var reqBody model.LoginRequest
	if err := req.BindValidate(r, &reqBody); err != nil {
		return err
	}

	uAgent := r.Header.Get("User-Agent")
	reqBody.UserAgent = uAgent

	accessToken, refreshToken, err := h.authService.Login(r.Context(), reqBody)
	if err != nil {
		return err
	}

	d := map[string]string{
		auth.ACCESS_TOKEN:  accessToken,
		auth.REFRESH_TOKEN: refreshToken,
	}

	if err := h.setTokenCookie(w, accessToken, auth.ACCESS_TOKEN); err != nil {
		return err
	}

	if err := h.setTokenCookie(w, refreshToken, auth.REFRESH_TOKEN); err != nil {
		return err
	}

	return res.Send(w, r, res.WithData(d), res.WithMessage("Login successful."))
}

func (h *AuthHandler) CurrentUser(w http.ResponseWriter, r *http.Request) error {
	user, err := h.authService.IdentifyUser(r.Context())
	if err != nil {
		return err
	}

	return res.Send(w, r, res.WithData(&user))
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) error {
	var refreshToken string

	if cookie, err := r.Cookie(auth.REFRESH_TOKEN); err == nil {
		refreshToken = cookie.Value
	}

	// If token doesn't exist in cookie, retrieve from Authorization header
	if refreshToken == "" {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			refreshToken = authHeader[len("Bearer "):]
		}
	}

	// If token is still empty, return an error
	if refreshToken == "" {
		return auth.ErrMissingToken
	}

	newAccessToken, err := h.authService.RefreshAccessToken(r.Context(), refreshToken)
	if err != nil {
		return err
	}

	d := map[string]string{
		auth.ACCESS_TOKEN: newAccessToken,
	}

	if err := h.setTokenCookie(w, newAccessToken, auth.ACCESS_TOKEN); err != nil {
		return err
	}

	return res.Send(
		w,
		r,
		res.WithData(d),
		res.WithMessage("Token refreshed."),
	)
}

func (h *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) error {
	// Retrieve the refresh token from the request cookie
	var refreshToken string
	if cookie, err := r.Cookie(auth.REFRESH_TOKEN); err == nil {
		refreshToken = cookie.Value
	}

	_ = h.authService.Logout(r.Context(), refreshToken)
	h.removeTokenCookies(w)

	return res.Send(w, r, res.WithMessage("Logout successful."))
}

func (h *AuthHandler) ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) error {
	var reqBody model.ForgotPasswordInput
	if err := req.BindValidate(r, &reqBody); err != nil {
		return err
	}

	if err := h.authService.ForgotPassword(r.Context(), reqBody); err != nil {
		return err
	}

	return res.Send(w, r, res.WithMessage("Password recovery email have been sent!"))
}

func (h *AuthHandler) setTokenCookie(
	w http.ResponseWriter,
	token string,
	label string,
) error {
	if label != auth.ACCESS_TOKEN && label != auth.REFRESH_TOKEN {
		return errors.New("invalid token label")
	}

	expiry := h.cfg.Token.AccessTokenExpiry
	if label == auth.REFRESH_TOKEN {
		expiry = h.cfg.Token.RefreshTokenExpiry
	}

	currTime := time.Now()

	ck := &http.Cookie{
		Name:     label,
		Value:    token,
		Expires:  currTime.Add(expiry),
		HttpOnly: h.cfg.IsProduction(),
		Secure:   h.cfg.IsProduction(),
		Path:     "/",
	}

	http.SetCookie(w, ck)
	return nil
}

func (h *AuthHandler) removeTokenCookies(w http.ResponseWriter) {
	atCookie := &http.Cookie{
		Name:     auth.ACCESS_TOKEN,
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: h.cfg.IsProduction(),
		Secure:   h.cfg.IsProduction(),
		Path:     "/",
	}

	rtCookie := *atCookie
	rtCookie.Name = auth.REFRESH_TOKEN

	http.SetCookie(w, atCookie)
	http.SetCookie(w, &rtCookie)
}
