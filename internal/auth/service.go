package auth

import (
	"context"
	"errors"
	"learn_golang/internal/model"
	"learn_golang/internal/repository"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{UserRepo: userRepo}
}

// Register user
func (s *AuthService) Register(ctx context.Context, user model.User) error {
	log.Printf("Auth.Register: registering email=%s name=%s", user.Email, user.Name)
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Auth.Register: bcrypt error: %v", err)
		return err
	}
	user.Password = string(hashed)

	_, err = s.UserRepo.CreateUser(ctx, user)
	if err != nil {
		log.Printf("Auth.Register: create user error: %v", err)
	}
	return err
}

// Login ตรวจสอบ username/password
func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	log.Printf("Auth.Login: attempting login for email=%s", email)
	user, err := s.UserRepo.GetUserByEmail(ctx, email)
	if err != nil {
		log.Printf("Auth.Login: GetUserByEmail error: %v", err)
		return "", errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Printf("Auth.Login: password mismatch for email=%s: %v", email, err)
		return "", errors.New("invalid email or password")
	}

	// สร้าง JWT token
	token, err := GenerateToken(user.ID.Hex())
	if err != nil {
		log.Printf("Auth.Login: GenerateToken error: %v", err)
		return "", err
	}

	log.Printf("Auth.Login: login successful for email=%s userID=%s", email, user.ID.Hex())
	return token, nil
}
