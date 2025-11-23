package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/internal/domain/auth"
	"github.com/prawirdani/golang-restapi/internal/domain/user"
	"github.com/prawirdani/golang-restapi/internal/testing/mocks"
)

func TestService_Register(t *testing.T) {
	ctx := context.Background()
	cfg := config.Auth{
		JwtSecret:                 "test-secret",
		JwtTTL:                    time.Hour,
		SessionTTL:                24 * time.Hour,
		ResetPasswordTTL:          time.Hour,
		ResetPasswordFormEndpoint: "http://localhost:3000/reset-password",
	}

	t.Run("Success", func(t *testing.T) {
		// Setup
		mockTransactor := mocks.NewTransactor(t)
		mockUserRepo := mocks.NewUserRepository(t)
		mockAuthRepo := mocks.NewAuthRepository(t)
		mockPublisher := mocks.NewAuthMessagePublisher(t)

		service := auth.NewService(cfg, mockTransactor, mockUserRepo, mockAuthRepo, mockPublisher)

		input := auth.RegisterInput{
			Name:           "John Doe",
			Email:          "john@example.com",
			Phone:          "1234567890",
			Password:       "password123",
			RepeatPassword: "password123",
		}

		// Mock expectations
		mockUserRepo.EXPECT().GetByEmail(ctx, input.Email).Return(nil, user.ErrNotFound)
		mockUserRepo.EXPECT().Store(ctx, mock.AnythingOfType("*user.User")).Return(nil)

		// Execute
		err := service.Register(ctx, input)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("EmailAlreadyExists", func(t *testing.T) {
		// Setup
		mockTransactor := mocks.NewTransactor(t)
		mockUserRepo := mocks.NewUserRepository(t)
		mockAuthRepo := mocks.NewAuthRepository(t)
		mockPublisher := mocks.NewAuthMessagePublisher(t)

		service := auth.NewService(cfg, mockTransactor, mockUserRepo, mockAuthRepo, mockPublisher)

		input := auth.RegisterInput{
			Name:           "John Doe",
			Email:          "john@example.com",
			Password:       "password123",
			RepeatPassword: "password123",
		}

		existingUser := &user.User{
			ID:    uuid.New(),
			Name:  "Existing User",
			Email: input.Email,
		}

		// Mock expectations
		mockUserRepo.EXPECT().GetByEmail(ctx, input.Email).Return(existingUser, nil)

		// Execute
		err := service.Register(ctx, input)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, user.ErrEmailExists, err)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		// Setup
		mockTransactor := mocks.NewTransactor(t)
		mockUserRepo := mocks.NewUserRepository(t)
		mockAuthRepo := mocks.NewAuthRepository(t)
		mockPublisher := mocks.NewAuthMessagePublisher(t)

		service := auth.NewService(cfg, mockTransactor, mockUserRepo, mockAuthRepo, mockPublisher)

		input := auth.RegisterInput{
			Name:           "John Doe",
			Email:          "john@example.com",
			Password:       "password123",
			RepeatPassword: "password123",
		}

		repoErr := errors.New("database error")

		// Mock expectations
		mockUserRepo.EXPECT().GetByEmail(ctx, input.Email).Return(nil, repoErr)

		// Execute
		err := service.Register(ctx, input)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, repoErr, err)
	})
}

func TestService_Login(t *testing.T) {
	ctx := context.Background()
	cfg := config.Auth{
		JwtSecret:                 "test-secret",
		JwtTTL:                    time.Hour,
		SessionTTL:                24 * time.Hour,
		ResetPasswordTTL:          time.Hour,
		ResetPasswordFormEndpoint: "http://localhost:3000/reset-password",
	}

	t.Run("Success", func(t *testing.T) {
		// Setup
		mockTransactor := mocks.NewTransactor(t)
		mockUserRepo := mocks.NewUserRepository(t)
		mockAuthRepo := mocks.NewAuthRepository(t)
		mockPublisher := mocks.NewAuthMessagePublisher(t)

		service := auth.NewService(cfg, mockTransactor, mockUserRepo, mockAuthRepo, mockPublisher)

		input := auth.LoginInput{
			Email:     "john@example.com",
			Password:  "password123",
			UserAgent: "test-agent",
		}

		hashedPassword, err := auth.HashPassword("password123")
		require.NoError(t, err)

		testUser := &user.User{
			ID:       uuid.New(),
			Name:     "John Doe",
			Email:    input.Email,
			Password: string(hashedPassword),
		}

		// Mock expectations
		mockUserRepo.EXPECT().GetByEmail(ctx, input.Email).Return(testUser, nil)
		mockAuthRepo.EXPECT().StoreSession(ctx, mock.AnythingOfType("*auth.Session")).Return(nil)

		// Execute
		accessToken, sessionID, err := service.Login(ctx, input)

		// Assert
		assert.NoError(t, err)
		assert.NotEmpty(t, accessToken)
		assert.NotEmpty(t, sessionID)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		// Setup
		mockTransactor := mocks.NewTransactor(t)
		mockUserRepo := mocks.NewUserRepository(t)
		mockAuthRepo := mocks.NewAuthRepository(t)
		mockPublisher := mocks.NewAuthMessagePublisher(t)

		service := auth.NewService(cfg, mockTransactor, mockUserRepo, mockAuthRepo, mockPublisher)

		input := auth.LoginInput{
			Email:    "nonexistent@example.com",
			Password: "password123",
		}

		// Mock expectations
		mockUserRepo.EXPECT().GetByEmail(ctx, input.Email).Return(nil, user.ErrNotFound)

		// Execute
		accessToken, sessionID, err := service.Login(ctx, input)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, auth.ErrWrongCredentials, err)
		assert.Empty(t, accessToken)
		assert.Empty(t, sessionID)
	})

	t.Run("WrongPassword", func(t *testing.T) {
		// Setup
		mockTransactor := mocks.NewTransactor(t)
		mockUserRepo := mocks.NewUserRepository(t)
		mockAuthRepo := mocks.NewAuthRepository(t)
		mockPublisher := mocks.NewAuthMessagePublisher(t)

		service := auth.NewService(cfg, mockTransactor, mockUserRepo, mockAuthRepo, mockPublisher)

		input := auth.LoginInput{
			Email:    "john@example.com",
			Password: "wrongpassword",
		}

		hashedPassword, err := auth.HashPassword("password123")
		require.NoError(t, err)

		testUser := &user.User{
			ID:       uuid.New(),
			Name:     "John Doe",
			Email:    input.Email,
			Password: string(hashedPassword),
		}

		// Mock expectations
		mockUserRepo.EXPECT().GetByEmail(ctx, input.Email).Return(testUser, nil)

		// Execute
		accessToken, sessionID, err := service.Login(ctx, input)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, accessToken)
		assert.Empty(t, sessionID)
	})
}

func TestService_RefreshAccessToken(t *testing.T) {
	ctx := context.Background()
	cfg := config.Auth{
		JwtSecret:                 "test-secret",
		JwtTTL:                    time.Hour,
		SessionTTL:                24 * time.Hour,
		ResetPasswordTTL:          time.Hour,
		ResetPasswordFormEndpoint: "http://localhost:3000/reset-password",
	}

	t.Run("Success", func(t *testing.T) {
		// Setup
		mockTransactor := mocks.NewTransactor(t)
		mockUserRepo := mocks.NewUserRepository(t)
		mockAuthRepo := mocks.NewAuthRepository(t)
		mockPublisher := mocks.NewAuthMessagePublisher(t)

		service := auth.NewService(cfg, mockTransactor, mockUserRepo, mockAuthRepo, mockPublisher)

		sessionID := uuid.New().String()
		userID := uuid.New()

		session, err := auth.NewSession(userID, "test-agent", cfg.SessionTTL)
		require.NoError(t, err)

		testUser := &user.User{
			ID:    userID,
			Name:  "John Doe",
			Email: "john@example.com",
		}

		// Mock expectations
		mockAuthRepo.EXPECT().GetSession(ctx, sessionID).Return(session, nil)
		mockUserRepo.EXPECT().GetByID(ctx, userID.String()).Return(testUser, nil)

		// Execute
		newAccessToken, err := service.RefreshAccessToken(ctx, sessionID)

		// Assert
		assert.NoError(t, err)
		assert.NotEmpty(t, newAccessToken)
	})

	t.Run("SessionExpired", func(t *testing.T) {
		// Setup
		mockTransactor := mocks.NewTransactor(t)
		mockUserRepo := mocks.NewUserRepository(t)
		mockAuthRepo := mocks.NewAuthRepository(t)
		mockPublisher := mocks.NewAuthMessagePublisher(t)

		service := auth.NewService(cfg, mockTransactor, mockUserRepo, mockAuthRepo, mockPublisher)

		sessionID := uuid.New().String()
		userID := uuid.New()

		// Create valid session and manually set it as expired
		session, err := auth.NewSession(userID, "test-agent", cfg.SessionTTL)
		require.NoError(t, err)
		session.ExpiresAt = time.Now().Add(-time.Hour) // Set to past

		// Mock expectations
		mockAuthRepo.EXPECT().GetSession(ctx, sessionID).Return(session, nil)

		// Execute
		newAccessToken, err := service.RefreshAccessToken(ctx, sessionID)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, auth.ErrSessionExpired, err)
		assert.Empty(t, newAccessToken)
	})
}

func TestService_Logout(t *testing.T) {
	ctx := context.Background()
	cfg := config.Auth{
		JwtSecret:                 "test-secret",
		JwtTTL:                    time.Hour,
		SessionTTL:                24 * time.Hour,
		ResetPasswordTTL:          time.Hour,
		ResetPasswordFormEndpoint: "http://localhost:3000/reset-password",
	}

	t.Run("Success", func(t *testing.T) {
		// Setup
		mockTransactor := mocks.NewTransactor(t)
		mockUserRepo := mocks.NewUserRepository(t)
		mockAuthRepo := mocks.NewAuthRepository(t)
		mockPublisher := mocks.NewAuthMessagePublisher(t)

		service := auth.NewService(cfg, mockTransactor, mockUserRepo, mockAuthRepo, mockPublisher)

		sessionID := uuid.New().String()
		userID := uuid.New()

		session, err := auth.NewSession(userID, "test-agent", cfg.SessionTTL)
		require.NoError(t, err)

		// Mock expectations
		mockAuthRepo.EXPECT().GetSession(ctx, sessionID).Return(session, nil)
		mockAuthRepo.EXPECT().UpdateSession(ctx, mock.AnythingOfType("*auth.Session")).Return(nil)

		// Execute
		err = service.Logout(ctx, sessionID)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("SessionAlreadyExpired", func(t *testing.T) {
		// Setup
		mockTransactor := mocks.NewTransactor(t)
		mockUserRepo := mocks.NewUserRepository(t)
		mockAuthRepo := mocks.NewAuthRepository(t)
		mockPublisher := mocks.NewAuthMessagePublisher(t)

		service := auth.NewService(cfg, mockTransactor, mockUserRepo, mockAuthRepo, mockPublisher)

		sessionID := uuid.New().String()
		userID := uuid.New()

		// Create valid session and manually set it as expired
		session, err := auth.NewSession(userID, "test-agent", cfg.SessionTTL)
		require.NoError(t, err)
		session.ExpiresAt = time.Now().Add(-time.Hour) // Set to past

		// Mock expectations
		mockAuthRepo.EXPECT().GetSession(ctx, sessionID).Return(session, nil)

		// Execute
		err = service.Logout(ctx, sessionID)

		// Assert
		assert.NoError(t, err) // Should not error even if session is expired
	})
}

func TestService_ForgotPassword(t *testing.T) {
	ctx := context.Background()
	cfg := config.Auth{
		JwtSecret:                 "test-secret",
		JwtTTL:                    time.Hour,
		SessionTTL:                24 * time.Hour,
		ResetPasswordTTL:          time.Hour,
		ResetPasswordFormEndpoint: "http://localhost:3000/reset-password",
	}

	t.Run("Success", func(t *testing.T) {
		// Setup
		mockTransactor := mocks.NewTransactor(t)
		mockUserRepo := mocks.NewUserRepository(t)
		mockAuthRepo := mocks.NewAuthRepository(t)
		mockPublisher := mocks.NewAuthMessagePublisher(t)

		service := auth.NewService(cfg, mockTransactor, mockUserRepo, mockAuthRepo, mockPublisher)

		input := auth.ForgotPasswordInput{
			Email: "john@example.com",
		}

		testUser := &user.User{
			ID:    uuid.New(),
			Name:  "John Doe",
			Email: input.Email,
		}

		// Mock expectations
		mockTransactor.EXPECT().Transact(ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil).Run(func(ctx context.Context, fn func(context.Context) error) {
			// Inside the transaction, we need to set up the expectations
			mockUserRepo.EXPECT().GetByEmail(ctx, input.Email).Return(testUser, nil)
			mockAuthRepo.EXPECT().StoreResetPasswordToken(ctx, mock.AnythingOfType("*auth.ResetPasswordToken")).Return(nil)
			mockPublisher.EXPECT().SendResetPasswordEmail(ctx, mock.AnythingOfType("auth.ResetPasswordEmailMessage")).Return(nil)

			err := fn(ctx)
			assert.NoError(t, err)
		})

		// Execute
		err := service.ForgotPassword(ctx, input)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		// Setup
		mockTransactor := mocks.NewTransactor(t)
		mockUserRepo := mocks.NewUserRepository(t)
		mockAuthRepo := mocks.NewAuthRepository(t)
		mockPublisher := mocks.NewAuthMessagePublisher(t)

		service := auth.NewService(cfg, mockTransactor, mockUserRepo, mockAuthRepo, mockPublisher)

		input := auth.ForgotPasswordInput{
			Email: "nonexistent@example.com",
		}

		// Mock expectations
		mockTransactor.EXPECT().Transact(ctx, mock.AnythingOfType("func(context.Context) error")).Run(func(ctx context.Context, fn func(context.Context) error) {
			mockUserRepo.EXPECT().GetByEmail(ctx, input.Email).Return(nil, user.ErrNotFound)

			err := fn(ctx)
			assert.Error(t, err)
			assert.Equal(t, user.ErrEmailNotVerified, err)
		}).Return(user.ErrEmailNotVerified)

		// Execute
		err := service.ForgotPassword(ctx, input)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, user.ErrEmailNotVerified, err)
	})
}

func TestService_ResetPassword(t *testing.T) {
	ctx := context.Background()
	cfg := config.Auth{
		JwtSecret:                 "test-secret",
		JwtTTL:                    time.Hour,
		SessionTTL:                24 * time.Hour,
		ResetPasswordTTL:          time.Hour,
		ResetPasswordFormEndpoint: "http://localhost:3000/reset-password",
	}

	t.Run("Success", func(t *testing.T) {
		// Setup
		mockTransactor := mocks.NewTransactor(t)
		mockUserRepo := mocks.NewUserRepository(t)
		mockAuthRepo := mocks.NewAuthRepository(t)
		mockPublisher := mocks.NewAuthMessagePublisher(t)

		service := auth.NewService(cfg, mockTransactor, mockUserRepo, mockAuthRepo, mockPublisher)

		userID := uuid.New()
		token, err := auth.NewResetPasswordToken(userID, cfg.ResetPasswordTTL)
		require.NoError(t, err)

		input := auth.ResetPasswordInput{
			Token:             token.Value,
			NewPassword:       "newpassword123",
			RepeatNewPassword: "newpassword123",
		}

		testUser := &user.User{
			ID:    userID,
			Name:  "John Doe",
			Email: "john@example.com",
		}

		// Mock expectations
		mockTransactor.EXPECT().Transact(ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil).Run(func(ctx context.Context, fn func(context.Context) error) {
			mockAuthRepo.EXPECT().GetResetPasswordToken(ctx, input.Token).Return(token, nil)
			mockUserRepo.EXPECT().GetByID(ctx, userID.String()).Return(testUser, nil)
			mockAuthRepo.EXPECT().UpdateResetPasswordToken(ctx, mock.AnythingOfType("*auth.ResetPasswordToken")).Return(nil)
			mockUserRepo.EXPECT().Update(ctx, mock.AnythingOfType("*user.User")).Return(nil)

			err := fn(ctx)
			assert.NoError(t, err)
		})

		// Execute
		err = service.ResetPassword(ctx, input)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("TokenExpired", func(t *testing.T) {
		// Setup
		mockTransactor := mocks.NewTransactor(t)
		mockUserRepo := mocks.NewUserRepository(t)
		mockAuthRepo := mocks.NewAuthRepository(t)
		mockPublisher := mocks.NewAuthMessagePublisher(t)

		service := auth.NewService(cfg, mockTransactor, mockUserRepo, mockAuthRepo, mockPublisher)

		userID := uuid.New()
		// Create valid token and manually set it as expired
		token, err := auth.NewResetPasswordToken(userID, cfg.ResetPasswordTTL)
		require.NoError(t, err)
		token.ExpiresAt = time.Now().Add(-time.Hour) // Set to past

		input := auth.ResetPasswordInput{
			Token:             token.Value,
			NewPassword:       "newpassword123",
			RepeatNewPassword: "newpassword123",
		}

		// Mock expectations
		mockTransactor.EXPECT().Transact(ctx, mock.AnythingOfType("func(context.Context) error")).Run(func(ctx context.Context, fn func(context.Context) error) {
			mockAuthRepo.EXPECT().GetResetPasswordToken(ctx, input.Token).Return(token, nil)

			err := fn(ctx)
			assert.Error(t, err)
			assert.Equal(t, auth.ErrResetPasswordTokenInvalid, err)
		}).Return(auth.ErrResetPasswordTokenInvalid)

		// Execute
		err = service.ResetPassword(ctx, input)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, auth.ErrResetPasswordTokenInvalid, err)
	})
}

func TestService_ChangePassword(t *testing.T) {
	ctx := context.Background()
	cfg := config.Auth{
		JwtSecret:                 "test-secret",
		JwtTTL:                    time.Hour,
		SessionTTL:                24 * time.Hour,
		ResetPasswordTTL:          time.Hour,
		ResetPasswordFormEndpoint: "http://localhost:3000/reset-password",
	}

	t.Run("Success", func(t *testing.T) {
		// Setup
		mockTransactor := mocks.NewTransactor(t)
		mockUserRepo := mocks.NewUserRepository(t)
		mockAuthRepo := mocks.NewAuthRepository(t)
		mockPublisher := mocks.NewAuthMessagePublisher(t)

		service := auth.NewService(cfg, mockTransactor, mockUserRepo, mockAuthRepo, mockPublisher)

		userID := uuid.New().String()
		oldPassword := "oldpassword123"
		newPassword := "newpassword123"

		hashedOldPassword, err := auth.HashPassword(oldPassword)
		require.NoError(t, err)

		testUser := &user.User{
			ID:       uuid.MustParse(userID),
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: string(hashedOldPassword),
		}

		input := auth.ChangePasswordInput{
			Password:          oldPassword,
			NewPassword:       newPassword,
			RepeatNewPassword: newPassword,
		}

		// Mock expectations
		mockUserRepo.EXPECT().GetByID(ctx, userID).Return(testUser, nil)
		mockUserRepo.EXPECT().Update(ctx, mock.AnythingOfType("*user.User")).Return(nil)

		// Execute
		err = service.ChangePassword(ctx, userID, input)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("WrongCurrentPassword", func(t *testing.T) {
		// Setup
		mockTransactor := mocks.NewTransactor(t)
		mockUserRepo := mocks.NewUserRepository(t)
		mockAuthRepo := mocks.NewAuthRepository(t)
		mockPublisher := mocks.NewAuthMessagePublisher(t)

		service := auth.NewService(cfg, mockTransactor, mockUserRepo, mockAuthRepo, mockPublisher)

		userID := uuid.New().String()
		oldPassword := "oldpassword123"
		wrongPassword := "wrongpassword"

		hashedOldPassword, err := auth.HashPassword(oldPassword)
		require.NoError(t, err)

		testUser := &user.User{
			ID:       uuid.MustParse(userID),
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: string(hashedOldPassword),
		}

		input := auth.ChangePasswordInput{
			Password:          wrongPassword,
			NewPassword:       "newpassword123",
			RepeatNewPassword: "newpassword123",
		}

		// Mock expectations
		mockUserRepo.EXPECT().GetByID(ctx, userID).Return(testUser, nil)

		// Execute
		err = service.ChangePassword(ctx, userID, input)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, auth.ErrWrongCredentials, err)
	})
}

func TestService_GetResetPasswordToken(t *testing.T) {
	ctx := context.Background()
	cfg := config.Auth{
		JwtSecret:                 "test-secret",
		JwtTTL:                    time.Hour,
		SessionTTL:                24 * time.Hour,
		ResetPasswordTTL:          time.Hour,
		ResetPasswordFormEndpoint: "http://localhost:3000/reset-password",
	}

	t.Run("Success", func(t *testing.T) {
		// Setup
		mockTransactor := mocks.NewTransactor(t)
		mockUserRepo := mocks.NewUserRepository(t)
		mockAuthRepo := mocks.NewAuthRepository(t)
		mockPublisher := mocks.NewAuthMessagePublisher(t)

		service := auth.NewService(cfg, mockTransactor, mockUserRepo, mockAuthRepo, mockPublisher)

		tokenValue := "test-token-value"
		userID := uuid.New()
		expectedToken, err := auth.NewResetPasswordToken(userID, cfg.ResetPasswordTTL)
		require.NoError(t, err)
		expectedToken.Value = tokenValue

		// Mock expectations
		mockAuthRepo.EXPECT().GetResetPasswordToken(ctx, tokenValue).Return(expectedToken, nil)

		// Execute
		token, err := service.GetResetPasswordToken(ctx, tokenValue)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedToken, token)
	})

	t.Run("TokenNotFound", func(t *testing.T) {
		// Setup
		mockTransactor := mocks.NewTransactor(t)
		mockUserRepo := mocks.NewUserRepository(t)
		mockAuthRepo := mocks.NewAuthRepository(t)
		mockPublisher := mocks.NewAuthMessagePublisher(t)

		service := auth.NewService(cfg, mockTransactor, mockUserRepo, mockAuthRepo, mockPublisher)

		tokenValue := "nonexistent-token"

		// Mock expectations
		mockAuthRepo.EXPECT().GetResetPasswordToken(ctx, tokenValue).Return(nil, auth.ErrResetPasswordTokenNotFound)

		// Execute
		token, err := service.GetResetPasswordToken(ctx, tokenValue)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, auth.ErrResetPasswordTokenNotFound, err)
		assert.Nil(t, token)
	})
}
