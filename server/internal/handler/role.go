package handler

import (
	"github.com/gin-gonic/gin"
	"jing_vue_gin_admin/server/internal/model"
	"jing_vue_gin_admin/server/internal/service"
	"net/http"
	"strconv"
)

type RoleHandler struct {
	roleService *service.RoleService
}

func NewRoleHandler() *RoleHandler {
	return &RoleHandler{
		roleService: service.NewRoleService(),
	}
}

// GetRoleList 获取角色列表
func (h *RoleHandler) GetRoleList(c *gin.Context) {
	roles, err := h.roleService.GetRoleList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取角色列表失败"})
		return
	}
	c.JSON(http.StatusOK, roles)
}

// CreateRole 创建角色
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req model.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	role := &model.Role{
		Name:        req.Name,
		Description: req.Description,
		Permissions: req.Permissions,
		Status:      1, // 默认启用
	}

	if err := h.roleService.CreateRole(role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建角色失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "创建角色成功"})
}

// UpdateRole 更新角色
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "角色ID格式错误"})
		return
	}

	var req model.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	if err := h.roleService.UpdateRole(uint(id), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新角色成功"})
}

// DeleteRole 删除角色
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "角色ID格式错误"})
		return
	}

	if err := h.roleService.DeleteRole(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除角色失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除角色成功"})
}

// ToggleRoleStatus 切换角色状态
func (h *RoleHandler) ToggleRoleStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "角色ID格式错误"})
		return
	}

	if err := h.roleService.ToggleRoleStatus(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "角色状态已更新"})
}

// GetAllPermissions 获取所有预定义权限
func (h *RoleHandler) GetAllPermissions(c *gin.Context) {
	permissions := h.roleService.GetAllPermissions()
	c.JSON(http.StatusOK, permissions)
} 