package http

import (
	"net/http"

	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/internal/model"
	"github.com/prawirdani/golang-restapi/internal/usecase"
	"github.com/prawirdani/golang-restapi/pkg/httputil"
)

type AuthHandler struct {
	userUC usecase.AuthUseCase
	cfg    config.Config
}

func NewAuthHandler(cfg config.Config, us usecase.AuthUseCase) AuthHandler {
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
	return httputil.SendJSON(w, http.StatusCreated, nil)
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

	response := model.TokenResponse{
		Token: tokenString,
	}

	return httputil.SendJSON(w, http.StatusOK, response)
}

func (h AuthHandler) CurrentUser(w http.ResponseWriter, r *http.Request) error {
	tokenClaims := httputil.GetAuthCtx(r.Context())
	response := model.TokenInfoResponse{
		TokenInfo: tokenClaims,
	}
	return httputil.SendJSON(w, http.StatusOK, response)
}
