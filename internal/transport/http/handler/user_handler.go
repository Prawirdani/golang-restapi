package handler

import (
	"net/http"

	"github.com/prawirdani/golang-restapi/internal/service"
	"github.com/prawirdani/golang-restapi/internal/transport/http/response"
	"github.com/prawirdani/golang-restapi/pkg/logging"
	"github.com/prawirdani/golang-restapi/pkg/uploader"
)

type UserHandler struct {
	logger          logging.Logger
	userService     *service.UserService
	imageFileParser *uploader.Parser
}

func NewUserHandler(logger logging.Logger, userService *service.UserService) *UserHandler {
	return &UserHandler{
		logger:          logger,
		userService:     userService,
		imageFileParser: uploader.New(uploader.ImageConfig),
	}
}

func (h *UserHandler) ChangeProfilePictureHandler(w http.ResponseWriter, r *http.Request) error {
	r.Body = http.MaxBytesReader(w, r.Body, 2<<20)

	file, err := h.imageFileParser.ParseSingleFile(r, "image")
	if err != nil {
		h.logger.Error(
			logging.TransportHTTP,
			"UserHandler.ChangeProfilePictureHandler",
			err.Error(),
		)
		return err
	}

	if err := h.userService.ChangeProfilePicture(r.Context(), file); err != nil {
		return err
	}

	return response.Send(w, r, response.WithMessage("Profile picture updated successfully!"))
}
