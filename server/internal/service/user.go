package service

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"jing_vue_gin_admin/server/config"
	"jing_vue_gin_admin/server/internal/model"
	"jing_vue_gin_admin/server/pkg/jwt"
	"log"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) Login(req *model.LoginRequest) (*model.LoginResponse, error) {
	log.Printf("Login attempt for user: %s", req.Username)
	
	var user model.User
	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		log.Printf("User not found: %v", err)
		return nil, errors.New("用户不存在")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		log.Printf("Password mismatch for user: %s", req.Username)
		return nil, errors.New("密码错误")
	}

	if user.Status != 1 {
		log.Printf("User account disabled: %s", req.Username)
		return nil, errors.New("账号已被禁用")
	}

	token, expiresAt, err := jwt.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		return nil, err
	}

	log.Printf("Login successful for user: %s", req.Username)
	return &model.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      user,
	}, nil
}

func (s *UserService) CreateUser(user *model.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return config.DB.Create(user).Error
} 