package handler

import (
	"net/http"

	"github.com/prawirdani/golang-restapi/internal/service"
	"github.com/prawirdani/golang-restapi/internal/transport/http/response"
	"github.com/prawirdani/golang-restapi/pkg/log"
	"github.com/prawirdani/golang-restapi/pkg/uploader"
)

type UserHandler struct {
	userService     *service.UserService
	imageFileParser *uploader.Parser
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService:     userService,
		imageFileParser: uploader.New(uploader.ImageConfig),
	}
}

func (h *UserHandler) ChangeProfilePictureHandler(w http.ResponseWriter, r *http.Request) error {
	r.Body = http.MaxBytesReader(w, r.Body, 2<<20)

	file, err := h.imageFileParser.ParseSingleFile(r, "image")
	if err != nil {
		log.Warn("failed to parse profile image file", "error", err.Error())
		return err
	}

	if err := h.userService.ChangeProfilePicture(r.Context(), file); err != nil {
		return err
	}

	return response.Send(w, r, response.WithMessage("Profile picture updated successfully!"))
}
