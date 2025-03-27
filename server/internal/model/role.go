package model

type Role struct {
	Base
	Name        string   `gorm:"size:32;uniqueIndex" json:"name"`         // 角色名称
	Description string   `gorm:"size:128" json:"description"`             // 角色描述
	Permissions []string `gorm:"serializer:json" json:"permissions"`      // 角色权限列表
	Status      int      `gorm:"default:1" json:"status"`                // 1: 启用, 0: 禁用
}

type CreateRoleRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Permissions []string `json:"permissions" binding:"required"`
}

type UpdateRoleRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Permissions []string `json:"permissions" binding:"required"`
}

// 预定义权限列表
var (
	PermissionUserView   = "user:view"   // 查看用户
	PermissionUserCreate = "user:create" // 创建用户
	PermissionUserEdit   = "user:edit"   // 编辑用户
	PermissionUserDelete = "user:delete" // 删除用户
	
	PermissionRoleView   = "role:view"   // 查看角色
	PermissionRoleCreate = "role:create" // 创建角色
	PermissionRoleEdit   = "role:edit"   // 编辑角色
	PermissionRoleDelete = "role:delete" // 删除角色
) 