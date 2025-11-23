package user_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/prawirdani/golang-restapi/internal/domain/user"
	"github.com/prawirdani/golang-restapi/internal/testing/mocks"
	"github.com/prawirdani/golang-restapi/pkg/nullable"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewUserService(t *testing.T) {
	mockTransactor := mocks.NewTransactor(t)
	mockUserRepo := mocks.NewUserRepository(t)
	mockImageStorage := mocks.NewStorage(t)

	service := user.NewService(mockTransactor, mockUserRepo, mockImageStorage)

	require.NotNil(t, service)
}

func TestUserService_GetUserByID(t *testing.T) {
	ctx := context.Background()

	t.Run("Success user without profile image", func(t *testing.T) {
		mockTransactor := mocks.NewTransactor(t)
		mockUserRepo := mocks.NewUserRepository(t)
		mockImageStorage := mocks.NewStorage(t)

		service := user.NewService(mockTransactor, mockUserRepo, mockImageStorage)

		userID := uuid.New().String()
		expectedUser := &user.User{
			ID:           uuid.New(),
			Name:         "John Doe",
			Email:        "john@example.com",
			Password:     "hashedpassword",
			Phone:        nullable.New("123456789", false),
			ProfileImage: nullable.New("", false),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		mockUserRepo.EXPECT().GetByID(ctx, userID).Return(expectedUser, nil)

		result, err := service.GetUserByID(ctx, userID)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, result)
	})

	t.Run("Success user with profile image", func(t *testing.T) {
		mockTransactor := mocks.NewTransactor(t)
		mockUserRepo := mocks.NewUserRepository(t)
		mockImageStorage := mocks.NewStorage(t)

		service := user.NewService(mockTransactor, mockUserRepo, mockImageStorage)

		userID := uuid.New().String()
		expectedUser := &user.User{
			ID:           uuid.New(),
			Name:         "John Doe",
			Email:        "john@example.com",
			Password:     "hashedpassword",
			Phone:        nullable.New("123456789", false),
			ProfileImage: nullable.New("profile.jpg", false),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		mockUserRepo.EXPECT().GetByID(ctx, userID).Return(expectedUser, nil)
		mockImageStorage.EXPECT().
			GetURL(ctx, "profiles/profile.jpg", time.Duration(0)).
			Return("https://example.com/profiles/profile.jpg", nil)

		result, err := service.GetUserByID(ctx, userID)
		assert.NoError(t, err)
		assert.Equal(t, "https://example.com/profiles/profile.jpg", result.ProfileImage.Get())
	})

	t.Run("Repository error", func(t *testing.T) {
		mockTransactor := mocks.NewTransactor(t)
		mockUserRepo := mocks.NewUserRepository(t)
		mockImageStorage := mocks.NewStorage(t)

		service := user.NewService(mockTransactor, mockUserRepo, mockImageStorage)

		userID := uuid.New()
		repoError := user.ErrNotFound

		mockUserRepo.EXPECT().GetByID(ctx, userID.String()).Return(nil, repoError)

		result, err := service.GetUserByID(ctx, userID.String())
		assert.Error(t, err)
		assert.Equal(t, repoError, err)
		assert.Nil(t, result)
	})

	t.Run("Storage error when getting profile URL", func(t *testing.T) {
		mockTransactor := mocks.NewTransactor(t)
		mockUserRepo := mocks.NewUserRepository(t)
		mockImageStorage := mocks.NewStorage(t)

		service := user.NewService(mockTransactor, mockUserRepo, mockImageStorage)

		userID := uuid.New().String()
		expectedUser := &user.User{
			ID:           uuid.New(),
			Name:         "John Doe",
			Email:        "john@example.com",
			Password:     "hashedpassword",
			Phone:        nullable.New("123456789", false),
			ProfileImage: nullable.New("profile.jpg", false),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		storageError := errors.New("storage error")
		mockUserRepo.EXPECT().GetByID(ctx, userID).Return(expectedUser, nil)
		mockImageStorage.EXPECT().GetURL(ctx, "profiles/profile.jpg", time.Duration(0)).Return("", storageError)

		result, err := service.GetUserByID(ctx, userID)
		assert.Error(t, err)
		assert.Equal(t, storageError, err)
		assert.Nil(t, result)
	})
}

func TestUserService_GetUserByEmail(t *testing.T) {
	ctx := context.Background()

	t.Run("Success user without profile image", func(t *testing.T) {
		mockTransactor := mocks.NewTransactor(t)
		mockUserRepo := mocks.NewUserRepository(t)
		mockImageStorage := mocks.NewStorage(t)

		service := user.NewService(mockTransactor, mockUserRepo, mockImageStorage)

		email := "john@example.com"
		expectedUser := &user.User{
			ID:           uuid.New(),
			Name:         "John Doe",
			Email:        email,
			Password:     "hashedpassword",
			Phone:        nullable.New("123456789", false),
			ProfileImage: nullable.New("", false),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		mockUserRepo.EXPECT().GetByEmail(ctx, email).Return(expectedUser, nil)

		result, err := service.GetUserByEmail(ctx, email)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, result)
	})

	t.Run("Success user with profile image", func(t *testing.T) {
		mockTransactor := mocks.NewTransactor(t)
		mockUserRepo := mocks.NewUserRepository(t)
		mockImageStorage := mocks.NewStorage(t)

		service := user.NewService(mockTransactor, mockUserRepo, mockImageStorage)

		email := "john@example.com"
		expectedUser := &user.User{
			ID:           uuid.New(),
			Name:         "John Doe",
			Email:        email,
			Password:     "hashedpassword",
			Phone:        nullable.New("123456789", false),
			ProfileImage: nullable.New("profile.jpg", false),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		mockUserRepo.EXPECT().GetByEmail(ctx, email).Return(expectedUser, nil)
		mockImageStorage.EXPECT().
			GetURL(ctx, "profiles/profile.jpg", time.Duration(0)).
			Return("https://example.com/profiles/profile.jpg", nil)

		result, err := service.GetUserByEmail(ctx, email)
		assert.NoError(t, err)
		assert.Equal(t, "https://example.com/profiles/profile.jpg", result.ProfileImage.Get())
	})

	t.Run("Repository error", func(t *testing.T) {
		mockTransactor := mocks.NewTransactor(t)
		mockUserRepo := mocks.NewUserRepository(t)
		mockImageStorage := mocks.NewStorage(t)

		service := user.NewService(mockTransactor, mockUserRepo, mockImageStorage)

		email := "nonexistent@example.com"
		repoError := user.ErrNotFound

		mockUserRepo.EXPECT().GetByEmail(ctx, email).Return(nil, repoError)

		result, err := service.GetUserByEmail(ctx, email)
		assert.Error(t, err)
		assert.Equal(t, repoError, err)
		assert.Nil(t, result)
	})

	t.Run("Storage error when getting profile URL", func(t *testing.T) {
		mockTransactor := mocks.NewTransactor(t)
		mockUserRepo := mocks.NewUserRepository(t)
		mockImageStorage := mocks.NewStorage(t)

		service := user.NewService(mockTransactor, mockUserRepo, mockImageStorage)

		email := "john@example.com"
		expectedUser := &user.User{
			ID:           uuid.New(),
			Name:         "John Doe",
			Email:        email,
			Password:     "hashedpassword",
			Phone:        nullable.New("123456789", false),
			ProfileImage: nullable.New("profile.jpg", false),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		storageError := errors.New("storage error")
		mockUserRepo.EXPECT().GetByEmail(ctx, email).Return(expectedUser, nil)
		mockImageStorage.EXPECT().GetURL(ctx, "profiles/profile.jpg", time.Duration(0)).Return("", storageError)

		result, err := service.GetUserByEmail(ctx, email)
		assert.Error(t, err)
		assert.Equal(t, storageError, err)
		assert.Nil(t, result)
	})
}

func TestUserService_ChangeProfilePicture(t *testing.T) {
	ctx := context.Background()

	t.Run("Success user without existing profile image", func(t *testing.T) {
		mockTransactor := mocks.NewTransactor(t)
		mockUserRepo := mocks.NewUserRepository(t)
		mockImageStorage := mocks.NewStorage(t)
		mockFile := mocks.NewFile(t)

		service := user.NewService(mockTransactor, mockUserRepo, mockImageStorage)

		userID := uuid.New().String()
		newFileName := "new-profile.jpg"
		newImagePath := "profiles/" + newFileName

		existingUser := &user.User{
			ID:           uuid.New(),
			Name:         "John Doe",
			Email:        "john@example.com",
			Password:     "hashedpassword",
			Phone:        nullable.New("123456789", false),
			ProfileImage: nullable.New("", false),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		mockTransactor.EXPECT().
			Transact(ctx, mock.AnythingOfType("func(context.Context) error")).
			Run(func(ctx context.Context, fn func(context.Context) error) {
				err := fn(ctx)
				assert.NoError(t, err)
			}).
			Return(nil)

		mockUserRepo.EXPECT().GetByID(ctx, userID).Return(existingUser, nil)
		mockFile.EXPECT().SetName(mock.AnythingOfType("string")).Return(nil)
		mockFile.EXPECT().Name().Return(newFileName)
		mockFile.EXPECT().ContentType().Return("image/jpeg")
		mockImageStorage.EXPECT().Put(ctx, newImagePath, mockFile, "image/jpeg").Return(nil)
		mockUserRepo.EXPECT().Update(ctx, mock.MatchedBy(func(u *user.User) bool {
			return u.ProfileImage.Get() == newFileName
		})).Return(nil)

		err := service.ChangeProfilePicture(ctx, userID, mockFile)
		assert.NoError(t, err)
	})

	t.Run("Error user not found", func(t *testing.T) {
		mockTransactor := mocks.NewTransactor(t)
		mockUserRepo := mocks.NewUserRepository(t)
		mockImageStorage := mocks.NewStorage(t)
		mockFile := mocks.NewFile(t)

		service := user.NewService(mockTransactor, mockUserRepo, mockImageStorage)

		userID := uuid.New().String()
		repoError := user.ErrNotFound

		mockTransactor.EXPECT().Transact(ctx, mock.AnythingOfType("func(context.Context) error")).Return(repoError)

		err := service.ChangeProfilePicture(ctx, userID, mockFile)
		assert.Error(t, err)
		assert.Equal(t, repoError, err)
	})

	t.Run("Error failed to set file name", func(t *testing.T) {
		mockTransactor := mocks.NewTransactor(t)
		mockUserRepo := mocks.NewUserRepository(t)
		mockImageStorage := mocks.NewStorage(t)
		mockFile := mocks.NewFile(t)

		service := user.NewService(mockTransactor, mockUserRepo, mockImageStorage)

		userID := uuid.New().String()
		fileError := errors.New("file error")

		existingUser := &user.User{
			ID:           uuid.New(),
			Name:         "John Doe",
			Email:        "john@example.com",
			Password:     "hashedpassword",
			Phone:        nullable.New("123456789", false),
			ProfileImage: nullable.New("", false),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		mockTransactor.EXPECT().
			Transact(ctx, mock.AnythingOfType("func(context.Context) error")).
			Run(func(ctx context.Context, fn func(context.Context) error) {
				err := fn(ctx)
				assert.Error(t, err)
				assert.Equal(t, fileError, err)
			}).
			Return(fileError)

		mockUserRepo.EXPECT().GetByID(ctx, userID).Return(existingUser, nil)
		mockFile.EXPECT().SetName(mock.AnythingOfType("string")).Return(fileError)

		err := service.ChangeProfilePicture(ctx, userID, mockFile)
		assert.Error(t, err)
		assert.Equal(t, fileError, err)
	})

	t.Run("Error failed to store image", func(t *testing.T) {
		mockTransactor := mocks.NewTransactor(t)
		mockUserRepo := mocks.NewUserRepository(t)
		mockImageStorage := mocks.NewStorage(t)
		mockFile := mocks.NewFile(t)

		service := user.NewService(mockTransactor, mockUserRepo, mockImageStorage)

		userID := uuid.New().String()
		newFileName := "new-profile.jpg"
		newImagePath := "profiles/" + newFileName
		storageError := errors.New("storage error")

		existingUser := &user.User{
			ID:           uuid.New(),
			Name:         "John Doe",
			Email:        "john@example.com",
			Password:     "hashedpassword",
			Phone:        nullable.New("123456789", false),
			ProfileImage: nullable.New("", false),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		mockTransactor.EXPECT().
			Transact(ctx, mock.AnythingOfType("func(context.Context) error")).
			Run(func(ctx context.Context, fn func(context.Context) error) {
				err := fn(ctx)
				assert.Error(t, err)
				assert.Equal(t, storageError, err)
			}).
			Return(storageError)

		mockUserRepo.EXPECT().GetByID(ctx, userID).Return(existingUser, nil)
		mockFile.EXPECT().SetName(mock.AnythingOfType("string")).Return(nil)
		mockFile.EXPECT().Name().Return(newFileName)
		mockFile.EXPECT().ContentType().Return("image/jpeg")
		mockImageStorage.EXPECT().Put(ctx, newImagePath, mockFile, "image/jpeg").Return(storageError)

		err := service.ChangeProfilePicture(ctx, userID, mockFile)
		assert.Error(t, err)
		assert.Equal(t, storageError, err)
	})

	t.Run("Error failed to update user", func(t *testing.T) {
		mockTransactor := mocks.NewTransactor(t)
		mockUserRepo := mocks.NewUserRepository(t)
		mockImageStorage := mocks.NewStorage(t)
		mockFile := mocks.NewFile(t)

		service := user.NewService(mockTransactor, mockUserRepo, mockImageStorage)

		userID := uuid.New().String()
		newFileName := "new-profile.jpg"
		newImagePath := "profiles/" + newFileName
		updateError := errors.New("update error")

		existingUser := &user.User{
			ID:           uuid.New(),
			Name:         "John Doe",
			Email:        "john@example.com",
			Password:     "hashedpassword",
			Phone:        nullable.New("123456789", false),
			ProfileImage: nullable.New("", false),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		mockTransactor.EXPECT().
			Transact(ctx, mock.AnythingOfType("func(context.Context) error")).
			Run(func(ctx context.Context, fn func(context.Context) error) {
				err := fn(ctx)
				assert.Error(t, err)
				assert.Equal(t, updateError, err)
			}).
			Return(updateError)

		mockUserRepo.EXPECT().GetByID(ctx, userID).Return(existingUser, nil)
		mockFile.EXPECT().SetName(mock.AnythingOfType("string")).Return(nil)
		mockFile.EXPECT().Name().Return(newFileName)
		mockFile.EXPECT().ContentType().Return("image/jpeg")
		mockImageStorage.EXPECT().Put(ctx, newImagePath, mockFile, "image/jpeg").Return(nil)
		mockUserRepo.EXPECT().Update(ctx, mock.MatchedBy(func(u *user.User) bool {
			return u.ProfileImage.Get() == newFileName
		})).Return(updateError)
		mockImageStorage.EXPECT().Delete(ctx, newImagePath).Return(nil)

		err := service.ChangeProfilePicture(ctx, userID, mockFile)
		assert.Error(t, err)
		assert.Equal(t, updateError, err)
	})

	t.Run("Error transaction fails", func(t *testing.T) {
		mockTransactor := mocks.NewTransactor(t)
		mockUserRepo := mocks.NewUserRepository(t)
		mockImageStorage := mocks.NewStorage(t)
		mockFile := mocks.NewFile(t)

		service := user.NewService(mockTransactor, mockUserRepo, mockImageStorage)

		userID := uuid.New().String()
		transactError := errors.New("transaction error")

		mockTransactor.EXPECT().Transact(ctx, mock.AnythingOfType("func(context.Context) error")).Return(transactError)

		err := service.ChangeProfilePicture(ctx, userID, mockFile)
		assert.Error(t, err)
		assert.Equal(t, transactError, err)
	})
}
