package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/internal/entity/user"
	"github.com/prawirdani/golang-restapi/internal/infra/repository"
	"github.com/prawirdani/golang-restapi/internal/infra/storage"
	"github.com/prawirdani/golang-restapi/pkg/common"
	"github.com/prawirdani/golang-restapi/pkg/contextx"
	"github.com/prawirdani/golang-restapi/pkg/logging"
)

type UserService struct {
	logger       logging.Logger
	cfg          *config.Config
	tx           repository.Transactor
	userRepo     user.Repository
	imageStorage storage.Storage
}

func NewUserService(
	cfg *config.Config,
	logger logging.Logger,
	transactor repository.Transactor,
	userRepo user.Repository,
	imageStorage storage.Storage,
) *UserService {
	return &UserService{
		logger:       logger,
		cfg:          cfg,
		tx:           transactor,
		userRepo:     userRepo,
		imageStorage: imageStorage,
	}
}

func (s *UserService) GetUserByID(ctx context.Context, userID string) (user.User, error) {
	u, err := s.userRepo.GetUserBy(ctx, "id", userID)
	if err != nil {
		return user.User{}, err
	}

	profileURL, err := s.imageStorage.GetURL(
		ctx,
		s.buildProfileImagePath(u.ProfileImage),
		1*time.Hour,
	)
	if err != nil {
		return user.User{}, err
	}
	u.ProfileImageURL = profileURL

	return u, nil
}

func (s *UserService) ChangeProfilePicture(ctx context.Context, file common.File) error {
	// 1. Retrieve User Data
	tokenPayload, err := contextx.GetAuthCtx(ctx)
	if err != nil {
		return err
	}
	return s.tx.Transact(ctx, func(ctx context.Context) error {
		u, err := s.userRepo.GetUserBy(ctx, "id", tokenPayload["id"].(string))
		if err != nil {
			return err
		}

		// 2. Prev image name + storage path
		prevImage := s.buildProfileImagePath(u.ProfileImage)

		// 3. Set New Image name using UUID
		if err := file.SetName(uuid.NewString()); err != nil {
			s.logger.Error(logging.Service, "UserService.file.SetName", err.Error())
			return err
		}

		newImage := file.Name()
		newImagePath := s.buildProfileImagePath(newImage)

		// 4. Store new image to storage
		if err := s.imageStorage.Put(ctx, newImagePath, file, file.ContentType()); err != nil {
			s.logger.Error(logging.Service, "UserService.ImageStorage.Put", err.Error())
			return err
		}

		// 5. Update user image_profile field and save to db
		u.ProfileImage = newImage
		if err := s.userRepo.UpdateUser(ctx, u); err != nil {
			// If Fail, Rollback & Delete Latest Image from storage
			if err := s.imageStorage.Delete(ctx, newImagePath); err != nil {
				// Non-Fatal
				s.logger.Error(logging.Service, "UserService.ChangeProfilePicture", err.Error())
			}
			return err
		}

		// -- Cleanup old image (Non Fatal: Should not rollback if error)
		go func(prevImage string) {
			if prevImage != s.buildProfileImagePath(user.DEFAULT_PROFILE_IMG) {
				if err := s.imageStorage.Delete(context.Background(), prevImage); err != nil {
					s.logger.Error(
						logging.Service,
						"UserService.ChangeProfilePicture.CleanupOldImage",
						err.Error(),
					)
				}
			}
		}(prevImage)

		return nil
	})
}

// imageName + ext
func (s *UserService) buildProfileImagePath(imageName string) string {
	return fmt.Sprintf("profiles/%s", imageName)
}
