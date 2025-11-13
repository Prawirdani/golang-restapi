package handler

import (
	"fmt"
	"net/http"

	"github.com/prawirdani/golang-restapi/internal/auth"
	"github.com/prawirdani/golang-restapi/internal/service"
	httperr "github.com/prawirdani/golang-restapi/internal/transport/http/error"
	"github.com/prawirdani/golang-restapi/internal/transport/http/uploader"
	"github.com/prawirdani/golang-restapi/pkg/log"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

const ImageFormKey = "image"

func (h *UserHandler) ChangeProfilePictureHandler(c *Context) error {
	fh, err := c.FormFile(ImageFormKey)
	if err != nil {
		if isMissingFileError(err) {
			return httperr.New(
				http.StatusBadRequest,
				fmt.Sprintf("missing required file '%s'", ImageFormKey),
				nil,
			)
		}
		log.ErrorCtx(c.Context(), "Failed to parse profile image form file", err)
		return err
	}

	file := uploader.NewParsedFile(fh)
	defer file.Close()

	if err := uploader.ValidateFile(c.Context(), file, uploader.ValidationRules{
		MaxSize:      1 << 20, // 2MB,
		AllowedMIMEs: uploader.ImageMIMEs,
	}); err != nil {
		return err
	}

	claims, err := auth.GetAccessTokenCtx(c.Context())
	if err != nil {
		return err
	}

	if err := h.userService.ChangeProfilePicture(c.Context(), claims.UserID, file); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &Body{
		Message: "Profile picture updated!",
	})
}
