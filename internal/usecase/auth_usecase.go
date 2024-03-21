package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prawirdani/golang-restapi/internal/entity"
	"github.com/prawirdani/golang-restapi/internal/model"
	"github.com/prawirdani/golang-restapi/internal/repository"
	"github.com/prawirdani/golang-restapi/pkg/httputil"
	"github.com/prawirdani/golang-restapi/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	db          *pgxpool.Pool
	userRepo    *repository.UserRepository
	jwtProvider *utils.JWTProvider
}

func NewAuthUseCase(db *pgxpool.Pool, ur *repository.UserRepository, jp *utils.JWTProvider) *AuthUseCase {
	return &AuthUseCase{
		db:          db,
		userRepo:    ur,
		jwtProvider: jp,
	}
}

func (u *AuthUseCase) CreateNewUser(ctx context.Context, request *model.UserRegisterRequest) error {
	if err := utils.ValidateStruct(request); err != nil {
		return err
	}

	// Hashing request password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	newUser := entity.User{
		ID:       uuid.New(),
		Name:     request.Name,
		Email:    request.Email,
		Password: string(hashedPassword),
	}

	tx, err := u.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := u.userRepo.Create(ctx, newUser, tx); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (u *AuthUseCase) Login(ctx context.Context, request *model.UserLoginRequest) (*string, error) {
	// Validate request body
	if err := utils.ValidateStruct(request); err != nil {
		return nil, err
	}

	// Query user from database by request email
	user, _ := u.userRepo.SelectByEmail(ctx, request.Email, u.db)

	// Helper function to check user password
	isPasswordMatch := func(storedPassword, requestPassword string) bool {
		err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(requestPassword))
		return err == nil
	}

	// Check if user exist and check is the password matched
	if user == nil || !isPasswordMatch(user.Password, request.Password) {
		return nil, httputil.ErrUnauthorized("Check your credentials")
	}

	// Generate jwt token
	tokenString, err := u.jwtProvider.CreateToken(user)
	if err != nil {
		return nil, err
	}

	return tokenString, nil
}
