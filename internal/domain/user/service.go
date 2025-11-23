// Package user provides the domain model and business logic for managing users in system.
package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/prawirdani/golang-restapi/internal/infrastructure/repository"
	"github.com/prawirdani/golang-restapi/internal/infrastructure/storage"
	"github.com/prawirdani/golang-restapi/pkg/log"
)

type Service struct {
	transactor   repository.Transactor
	userRepo     Repository
	imageStorage storage.Storage
}

func NewService(
	transactor repository.Transactor,
	userRepo Repository,
	imageStorage storage.Storage,
) *Service {
	return &Service{
		transactor:   transactor,
		userRepo:     userRepo,
		imageStorage: imageStorage,
	}
}

func (s *Service) GetUserByID(ctx context.Context, userID string) (*User, error) {
	u, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if err := s.assignProfileImageURL(ctx, u); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	u, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if err := s.assignProfileImageURL(ctx, u); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *Service) ChangeProfilePicture(
	ctx context.Context,
	userID string,
	file storage.File,
) error {
	return s.transactor.Transact(ctx, func(ctx context.Context) error {
		u, err := s.userRepo.GetByID(ctx, userID)
		if err != nil {
			return err
		}

		//  Prev image name + storage path for cleanup
		var prevImage string
		if u.ProfileImage.Valid() {
			prevImage = s.buildProfileImagePath(u.ProfileImage.Get())
		}

		//  Set New Image name using UUID
		if err := file.SetName(uuid.NewString()); err != nil {
			log.ErrorCtx(ctx, "Failed to set profile image file name", err)
			return err
		}

		newImageName := file.Name()
		newImagePath := s.buildProfileImagePath(newImageName)

		//  Store new image to storage
		if err := s.imageStorage.Put(ctx, newImagePath, file, file.ContentType()); err != nil {
			log.ErrorCtx(ctx, "Failed store new profile image", err)
			return err
		}

		//  Update user image_profile field and save to db
		u.ProfileImage.Set(newImageName, false)
		if err := s.userRepo.Update(ctx, u); err != nil {
			// If Fail, Rollback & Delete Latest Image from storage
			if err := s.imageStorage.Delete(ctx, newImagePath); err != nil {
				// Non-Fatal
				log.WarnCtx(ctx, "Failed cleanup new profile image", "error", err.Error())
			}
			return err
		}

		// -- Cleanup old image (Non Fatal: Should not rollback if error)
		go func(prevImage string) {
			if prevImage != "" {
				if err := s.imageStorage.Delete(context.Background(), prevImage); err != nil {
					log.WarnCtx(ctx, "Failed cleanup old profile image", "error", err.Error())
				}
			}
		}(prevImage)

		return nil
	})
}

// imageName + ext
func (s *Service) buildProfileImagePath(imageName string) string {
	return fmt.Sprintf("profiles/%s", imageName)
}

// WARNING: Do not use inside for loop if using private bucket / presigned url
func (s *Service) assignProfileImageURL(ctx context.Context, u *User) error {
	// If null
	if !u.ProfileImage.Valid() {
		return nil
	}

	profileURL, err := s.imageStorage.GetURL(ctx, s.buildProfileImagePath(u.ProfileImage.Get()), 0)
	if err != nil {
		return err
	}
	u.ProfileImage.Set(profileURL, false)

	return nil
}
