package http

import (
	"net/http"

	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/internal/model"
	"github.com/prawirdani/golang-restapi/internal/usecase"
	"github.com/prawirdani/golang-restapi/pkg/httputil"
	"github.com/prawirdani/golang-restapi/pkg/utils"
)

type AuthHandler struct {
	authUC usecase.AuthUseCase
	cfg    *config.Config
}

func NewAuthHandler(cfg *config.Config, us usecase.AuthUseCase) AuthHandler {
	return AuthHandler{
		authUC: us,
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

	if err := h.authUC.Register(r.Context(), reqBody); err != nil {
		return err
	}

	return response(w, status(201), message("Registration successful."))
}

func (h AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) error {
	var reqBody model.LoginRequest
	if err := httputil.BindJSON(r, &reqBody); err != nil {
		return err
	}

	if err := reqBody.ValidateRequest(); err != nil {
		return err
	}

	tokens, err := h.authUC.Login(r.Context(), reqBody)
	if err != nil {
		return err
	}

	d := make(map[string]string)

	for i := 0; i < len(tokens); i++ {
		tokenTypeName := "accessToken"
		if tokens[i].Claims.TokenType == utils.RefreshToken {
			tokenTypeName = "refreshToken"
		}
		d[tokenTypeName] = tokens[i].String()
		tokens[i].SetCookie(w)
	}

	return response(w, data(d), message("Login successful."))
}

func (h AuthHandler) CurrentUser(w http.ResponseWriter, r *http.Request) error {
	tokenClaims := httputil.GetAuthCtx(r.Context())

	return response(w, data(tokenClaims))
}

func (h AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) error {
	claims, err := utils.ParseJWT(r, &h.cfg.Token, utils.RefreshToken)
	if err != nil {
		return err
	}

	userPayload, ok := claims["user"].(map[string]interface{})
	if !ok {
		return httputil.ErrBadRequest("Missing user payload in token.")
	}

	userID := userPayload["id"].(string)
	newAccessToken, err := h.authUC.RefreshToken(r.Context(), userID)
	if err != nil {
		return err
	}

	d := map[string]string{
		"accessToken": newAccessToken.String(),
	}
	newAccessToken.SetCookie(w)

	return response(w, data(d), message("Token refreshed."))
}
