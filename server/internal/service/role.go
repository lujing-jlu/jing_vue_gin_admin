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
func (s *RoleService) GetAllPermissions() map[string][]PermissionInfo {
	return map[string][]PermissionInfo{
		"系统管理": {
			{Value: model.PermissionUserView, Label: "查看用户"},
			{Value: model.PermissionUserCreate, Label: "创建用户"},
			{Value: model.PermissionUserEdit, Label: "编辑用户"},
			{Value: model.PermissionUserDelete, Label: "删除用户"},
			{Value: model.PermissionRoleView, Label: "查看角色"},
			{Value: model.PermissionRoleCreate, Label: "创建角色"},
			{Value: model.PermissionRoleEdit, Label: "编辑角色"},
			{Value: model.PermissionRoleDelete, Label: "删除角色"},
			{Value: model.PermissionSystemConfig, Label: "系统配置"},
			{Value: model.PermissionSystemLog, Label: "系统日志"},
		},
		"内容管理": {
			{Value: model.PermissionArticleView, Label: "查看文章"},
			{Value: model.PermissionArticleCreate, Label: "创建文章"},
			{Value: model.PermissionArticleEdit, Label: "编辑文章"},
			{Value: model.PermissionArticleDelete, Label: "删除文章"},
			{Value: model.PermissionArticlePublish, Label: "发布文章"},
			{Value: model.PermissionCategoryView, Label: "查看分类"},
			{Value: model.PermissionCategoryCreate, Label: "创建分类"},
			{Value: model.PermissionCategoryEdit, Label: "编辑分类"},
			{Value: model.PermissionCategoryDelete, Label: "删除分类"},
			{Value: model.PermissionTagView, Label: "查看标签"},
			{Value: model.PermissionTagCreate, Label: "创建标签"},
			{Value: model.PermissionTagEdit, Label: "编辑标签"},
			{Value: model.PermissionTagDelete, Label: "删除标签"},
			{Value: model.PermissionCommentView, Label: "查看评论"},
			{Value: model.PermissionCommentReply, Label: "回复评论"},
			{Value: model.PermissionCommentDelete, Label: "删除评论"},
			{Value: model.PermissionCommentApprove, Label: "审核评论"},
		},
		"媒体管理": {
			{Value: model.PermissionMediaView, Label: "查看媒体"},
			{Value: model.PermissionMediaUpload, Label: "上传媒体"},
			{Value: model.PermissionMediaEdit, Label: "编辑媒体"},
			{Value: model.PermissionMediaDelete, Label: "删除媒体"},
		},
		"数据统计": {
			{Value: model.PermissionStatView, Label: "查看统计"},
			{Value: model.PermissionStatExport, Label: "导出统计"},
			{Value: model.PermissionStatManage, Label: "管理统计"},
		},
	}
}

type PermissionInfo struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// GetRoles 获取角色列表，支持分页和搜索
func (s *RoleService) GetRoles(page, pageSize int, keyword string, status *int) (map[string]interface{}, error) {
	var roles []model.Role
	var total int64
	
	query := config.DB.Model(&model.Role{})
	
	// 应用搜索条件
	if keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	
	// 计算总记录数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}
	
	// 获取分页数据
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&roles).Error; err != nil {
		return nil, err
	}
	
	return map[string]interface{}{
		"total": total,
		"list":  roles,
	}, nil
} 