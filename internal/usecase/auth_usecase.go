package usecase

import (
	"auth-service/internal/domain"
	"auth-service/internal/repository"
	"auth-service/internal/utils"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type authUsecase struct {
	userRepo     domain.UserRepository
	redisRepo    *repository.RedisRepository
	emailService utils.EmailService
}

func NewAuthUsecase(
	userRepo domain.UserRepository,
	redisRepo *repository.RedisRepository,
	emailService utils.EmailService,
) domain.UserUsecase {
	return &authUsecase{
		userRepo:     userRepo,
		redisRepo:    redisRepo,
		emailService: emailService,
	}
}

func (u *authUsecase) Register(user *domain.User) error {
	// Check if user exists
	existingUser, _ := u.userRepo.GetByEmail(user.Email)
	if existingUser != nil {
		return errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	// Create user
	user.ID = uuid.New().String()
	user.Password = hashedPassword
	user.IsActive = false
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if err := u.userRepo.Create(user); err != nil {
		return err
	}

	// Generate and send OTP
	otp := utils.GenerateOTP()
	err = u.redisRepo.StoreOTP(context.Background(), user.Email, otp)
	if err != nil {
		return err
	}

	return u.emailService.SendOTP(user.Email, otp)
}

func (u *authUsecase) Login(email, password string) (string, error) {
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return "", domain.ErrInactiveUser
	}

	if !user.IsActive {
		return "", domain.ErrInvalidCredentials
	}

	// Generate JWT token
	return utils.GenerateJWT(user)
}

func (u *authUsecase) VerifyOTP(email, otp string) error {
	storedOTP, err := u.redisRepo.GetOTP(context.Background(), email)
	if err != nil {
		return errors.New("invalid or expired OTP")
	}

	if storedOTP != otp {
		return errors.New("invalid OTP")
	}

	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return err
	}

	return u.userRepo.UpdateActive(user.ID, true)
}

func (u *authUsecase) ResendOTP(email string) error {
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return err
	}

	if user.IsActive {
		return errors.New("account already verified")
	}

	otp := utils.GenerateOTP()
	err = u.redisRepo.StoreOTP(context.Background(), email, otp)
	if err != nil {
		return err
	}

	return utils.SendOTPEmail(email, otp)
}

func (u *authUsecase) GetUserByID(id string) (*domain.User, error) {
	if id == "" {
		return nil, errors.New("invalid user ID")
	}

	user, err := u.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *authUsecase) DeleteUser(id string, permanent bool) error {
	// First check if user exists and is active
	user, err := u.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	if user.IsDeleted {
		return domain.ErrUserNotFound
	}

	if permanent {
		return u.userRepo.HardDelete(id)
	}

	return u.userRepo.SoftDelete(id)
}
