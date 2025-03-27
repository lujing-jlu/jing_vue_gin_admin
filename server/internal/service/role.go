package service

import (
	"errors"
	"jing_vue_gin_admin/server/config"
	"jing_vue_gin_admin/server/internal/model"
)

type RoleService struct{}

func NewRoleService() *RoleService {
	return &RoleService{}
}

// GetRoleList 获取所有角色列表
func (s *RoleService) GetRoleList() ([]model.Role, error) {
	var roles []model.Role
	if err := config.DB.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// GetRoleByID 根据ID获取角色
func (s *RoleService) GetRoleByID(id uint) (*model.Role, error) {
	var role model.Role
	if err := config.DB.First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// CreateRole 创建角色
func (s *RoleService) CreateRole(role *model.Role) error {
	return config.DB.Create(role).Error
}

// UpdateRole 更新角色信息
func (s *RoleService) UpdateRole(id uint, roleData *model.UpdateRoleRequest) error {
	var role model.Role
	if err := config.DB.First(&role, id).Error; err != nil {
		return errors.New("角色不存在")
	}

	role.Name = roleData.Name
	role.Description = roleData.Description
	role.Permissions = roleData.Permissions

	return config.DB.Save(&role).Error
}

// DeleteRole 删除角色
func (s *RoleService) DeleteRole(id uint) error {
	return config.DB.Delete(&model.Role{}, id).Error
}

// ToggleRoleStatus 切换角色状态
func (s *RoleService) ToggleRoleStatus(id uint) error {
	var role model.Role
	if err := config.DB.First(&role, id).Error; err != nil {
		return errors.New("角色不存在")
	}

	if role.Status == 1 {
		role.Status = 0
	} else {
		role.Status = 1
	}

	return config.DB.Save(&role).Error
}

// GetAllPermissions 获取所有预定义权限
func (s *RoleService) GetAllPermissions() []string {
	return []string{
		model.PermissionUserView,
		model.PermissionUserCreate,
		model.PermissionUserEdit,
		model.PermissionUserDelete,
		model.PermissionRoleView,
		model.PermissionRoleCreate,
		model.PermissionRoleEdit,
		model.PermissionRoleDelete,
	}
} 