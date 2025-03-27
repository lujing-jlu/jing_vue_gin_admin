package model

import (
	"time"
)

type User struct {
	Base
	Username string `gorm:"uniqueIndex;size:32" json:"username"`
	Password string `gorm:"size:128" json:"-"`
	Nickname string `gorm:"size:32" json:"nickname"`
	Avatar   string `gorm:"size:256" json:"avatar"`
	Role     string `gorm:"size:16" json:"role"`
	Status   int    `gorm:"default:1" json:"status"` // 1: 正常, 0: 禁用
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	User      User      `json:"user"`
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

type UpdateUserRequest struct {
	Nickname string `json:"nickname" binding:"required"`
	Role     string `json:"role" binding:"required"`
} 