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

// 系统管理权限
const (
	PermissionUserView   = "user:view"   // 查看用户
	PermissionUserCreate = "user:create" // 创建用户
	PermissionUserEdit   = "user:update" // 编辑用户
	PermissionUserDelete = "user:delete" // 删除用户
	
	PermissionRoleView   = "role:view"   // 查看角色
	PermissionRoleCreate = "role:create" // 创建角色
	PermissionRoleEdit   = "role:update" // 编辑角色
	PermissionRoleDelete = "role:delete" // 删除角色

	PermissionSystemConfig = "system:config" // 系统配置
	PermissionSystemLog    = "log:view"      // 查看日志
	PermissionLogDelete    = "log:delete"    // 删除日志
)

// 内容管理权限
const (
	PermissionArticleView   = "article:view"   // 查看文章
	PermissionArticleCreate = "article:create" // 创建文章
	PermissionArticleEdit   = "article:edit"   // 编辑文章
	PermissionArticleDelete = "article:delete" // 删除文章
	PermissionArticlePublish = "article:publish" // 发布文章

	PermissionCategoryView   = "category:view"   // 查看分类
	PermissionCategoryCreate = "category:create" // 创建分类
	PermissionCategoryEdit   = "category:edit"   // 编辑分类
	PermissionCategoryDelete = "category:delete" // 删除分类

	PermissionTagView   = "tag:view"   // 查看标签
	PermissionTagCreate = "tag:create" // 创建标签
	PermissionTagEdit   = "tag:edit"   // 编辑标签
	PermissionTagDelete = "tag:delete" // 删除标签

	PermissionCommentView    = "comment:view"    // 查看评论
	PermissionCommentReply   = "comment:reply"   // 回复评论
	PermissionCommentDelete  = "comment:delete"  // 删除评论
	PermissionCommentApprove = "comment:approve" // 审核评论
)

// 媒体管理权限
const (
	PermissionMediaView   = "media:view"   // 查看媒体
	PermissionMediaUpload = "media:upload" // 上传媒体
	PermissionMediaEdit   = "media:edit"   // 编辑媒体
	PermissionMediaDelete = "media:delete" // 删除媒体
)

// 数据统计权限
const (
	PermissionStatView    = "stat:view"    // 查看统计
	PermissionStatExport  = "stat:export"  // 导出统计
	PermissionStatManage  = "stat:manage"  // 管理统计
)

// 文件管理权限
const (
	PermissionFileView   = "file:view"   // 查看文件
	PermissionFileUpload = "file:upload" // 上传文件
	PermissionFileUpdate = "file:update" // 更新文件
	PermissionFileDelete = "file:delete" // 删除文件
) 