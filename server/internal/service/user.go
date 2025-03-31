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

// GetUserList 获取所有用户列表
func (s *UserService) GetUserList() ([]model.User, error) {
	var users []model.User
	if err := config.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	if err := config.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(id uint, userData *model.UpdateUserRequest) error {
	var user model.User
	if err := config.DB.First(&user, id).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 更新用户信息
	user.Nickname = userData.Nickname
	user.Role = userData.Role

	return config.DB.Save(&user).Error
}

// ToggleUserStatus 切换用户状态
func (s *UserService) ToggleUserStatus(id uint) error {
	var user model.User
	if err := config.DB.First(&user, id).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 切换状态：1 -> 0, 0 -> 1
	if user.Status == 1 {
		user.Status = 0
	} else {
		user.Status = 1
	}

	return config.DB.Save(&user).Error
}

// UpdateProfile 更新用户个人信息
func (s *UserService) UpdateProfile(id uint, nickname string, avatar string) error {
	// 查找用户
	var user model.User
	if err := config.DB.First(&user, id).Error; err != nil {
		return err
	}

	// 更新用户信息
	updates := map[string]interface{}{
		"nickname": nickname,
	}
	
	// 如果提供了头像URL，则更新头像
	if avatar != "" {
		updates["avatar"] = avatar
	}

	return config.DB.Model(&user).Updates(updates).Error
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(id uint, oldPassword string, newPassword string) error {
	// 查找用户
	var user model.User
	if err := config.DB.First(&user, id).Error; err != nil {
		return err
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("原密码不正确")
	}

	// 生成新密码的哈希值
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 更新密码
	return config.DB.Model(&user).Update("password", string(hashedPassword)).Error
} 