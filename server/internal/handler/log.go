package handler

import (
	"github.com/gin-gonic/gin"
	"jing_vue_gin_admin/server/internal/model"
	"jing_vue_gin_admin/server/internal/service"
	"net/http"
)

type LogHandler struct {
	logService *service.LogService
}

// NewLogHandler 创建日志处理器
func NewLogHandler() *LogHandler {
	return &LogHandler{
		logService: service.NewLogService(),
	}
}

// GetLogList 获取日志列表
func (h *LogHandler) GetLogList(c *gin.Context) {
	var req model.LogListRequest
	
	// 设置默认值
	req.Page = 1
	req.PageSize = 10
	
	// 绑定请求参数
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}
	
	// 获取日志列表
	resp, err := h.logService.GetLogList(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取日志列表失败"})
		return
	}
	
	c.JSON(http.StatusOK, resp)
}

// DeleteLogs 批量删除日志
func (h *LogHandler) DeleteLogs(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}
	
	if err := h.logService.DeleteLogs(req.IDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除日志失败"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// ClearLogs 清空日志
func (h *LogHandler) ClearLogs(c *gin.Context) {
	if err := h.logService.ClearLogs(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "清空日志失败"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "清空成功"})
}

// GetLogModules 获取日志模块列表
func (h *LogHandler) GetLogModules(c *gin.Context) {
	modules := []string{
		model.LogModuleUser,
		model.LogModuleRole,
		model.LogModuleSystem,
		model.LogModuleAuth,
	}
	
	c.JSON(http.StatusOK, modules)
}

// GetLogActions 获取日志操作类型列表
func (h *LogHandler) GetLogActions(c *gin.Context) {
	actions := []string{
		model.LogActionLogin,
		model.LogActionLogout,
		model.LogActionCreate,
		model.LogActionUpdate,
		model.LogActionDelete,
		model.LogActionView,
		model.LogActionExport,
		model.LogActionImport,
		model.LogActionEnable,
		model.LogActionDisable,
	}
	
	c.JSON(http.StatusOK, actions)
} 