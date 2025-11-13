package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/prawirdani/golang-restapi/pkg/log"

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

func (h *AuthHandler) RegisterHandler(c *Context) error {
	var reqBody model.CreateUserInput
	if err := c.BindValidate(&reqBody); err != nil {
		log.ErrorCtx(c.Context(), BindValidateWarnLog, err)
		return err
	}

	if err := h.authService.Register(c.Context(), reqBody); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, Body{
		Message: "Registration successful",
	})
}

func (h *AuthHandler) LoginHandler(c *Context) error {
	var reqBody model.LoginInput
	if err := c.BindValidate(&reqBody); err != nil {
		log.ErrorCtx(c.Context(), BindValidateWarnLog, err)
		return err
	}
	reqBody.UserAgent = c.Get("User-Agent")

	accessToken, refreshToken, err := h.authService.Login(c.Context(), reqBody)
	if err != nil {
		return err
	}

	tp := model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	c.SetCookie(h.createTokenCookie(accessToken, auth.ACCESS_TOKEN))
	c.SetCookie(h.createTokenCookie(refreshToken, auth.REFRESH_TOKEN))

	return c.JSON(200, Body{
		Data: tp,
	})
}

func (h *AuthHandler) CurrentUserHandler(c *Context) error {
	user, err := h.authService.IdentifyUser(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, Body{
		Data: user,
	})
}

func (h *AuthHandler) RefreshTokenHandler(c *Context) error {
	var refreshToken string

	if cookie, err := c.GetCookie(auth.REFRESH_TOKEN); err == nil {
		refreshToken = cookie.Value
	}

	// If token doesn't exist in cookie, retrieve from Authorization header
	if refreshToken == "" {
		authHeader := c.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			refreshToken = authHeader[len("Bearer "):]
		}
	}

	// If token is still empty, return an error
	if refreshToken == "" {
		return auth.ErrTokenNotProvided
	}

	newAccessToken, err := h.authService.RefreshAccessToken(c.Context(), refreshToken)
	if err != nil {
		return err
	}

	d := map[string]string{
		auth.ACCESS_TOKEN: newAccessToken,
	}

	c.SetCookie(h.createTokenCookie(newAccessToken, auth.ACCESS_TOKEN))

	return c.JSON(http.StatusOK, Body{
		Data:    d,
		Message: "Token refreshed",
	})
}

func (h *AuthHandler) LogoutHandler(c *Context) error {
	// Retrieve the refresh token from the request cookie
	var refreshToken string
	if cookie, err := c.GetCookie(auth.REFRESH_TOKEN); err == nil {
		refreshToken = cookie.Value
	}

	_ = h.authService.Logout(c.Context(), refreshToken)
	h.removeTokenCookies(c)

	return c.JSON(http.StatusOK, Body{
		Message: "Logged out",
	})
}

func (h *AuthHandler) ForgotPasswordHandler(c *Context) error {
	var reqBody model.ForgotPasswordInput
	if err := c.BindValidate(&reqBody); err != nil {
		log.ErrorCtx(c.Context(), BindValidateWarnLog, err)
		return err
	}

	if err := h.authService.ForgotPassword(c.Context(), reqBody); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, Body{
		Message: "Password recovery email have been sent!",
	})
}

func (h *AuthHandler) GetResetPasswordTokenHandler(c *Context) error {
	token := c.Param("token")

	tokenObj, err := h.authService.GetResetPasswordToken(c.Context(), token)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, Body{
		Data: tokenObj,
	})
}

func (h *AuthHandler) ResetPasswordHandler(c *Context) error {
	var reqBody model.ResetPasswordInput
	if err := c.BindValidate(&reqBody); err != nil {
		log.ErrorCtx(c.Context(), BindValidateWarnLog, err)
		return err
	}

	if err := h.authService.ResetPassword(c.Context(), reqBody); err != nil {
		return err
	}

	return c.JSON(200, Body{
		Message: "Password has been reset successfuly",
	})
}

func (h *AuthHandler) ChangePasswordHandler(c *Context) error {
	var reqBody model.ChangePasswordInput
	if err := c.BindValidate(&reqBody); err != nil {
		log.ErrorCtx(c.Context(), BindValidateWarnLog, err)
		return err
	}

	if err := h.authService.ChangePassword(c.Context(), reqBody); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, Body{
		Message: "Password has been reset successfuly!",
	})
}

func (h *AuthHandler) createTokenCookie(
	token string,
	label string,
) *http.Cookie {
	expiry := h.cfg.Token.AccessTokenExpiry
	if label == auth.REFRESH_TOKEN {
		expiry = h.cfg.Token.RefreshTokenExpiry
	}

	currTime := time.Now()

	return &http.Cookie{
		Name:     label,
		Value:    token,
		Expires:  currTime.Add(expiry),
		HttpOnly: h.cfg.IsProduction(),
		Secure:   h.cfg.IsProduction(),
		Path:     "/",
	}
}

func (h *AuthHandler) removeTokenCookies(c *Context) {
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

	c.SetCookie(atCookie)
	c.SetCookie(&rtCookie)
}
