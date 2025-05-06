package handler

import (
	"jing_vue_gin_admin/server/internal/model"
	"jing_vue_gin_admin/server/internal/service"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

// FileHandler 文件处理器
type FileHandler struct {
	fileService *service.FileService
}

// NewFileHandler 创建文件处理器
func NewFileHandler() *FileHandler {
	return &FileHandler{
		fileService: service.NewFileService(),
	}
}

// UploadFile 上传文件
func (h *FileHandler) UploadFile(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "message": "用户未认证"})
		return
	}

	// 绑定请求
	var req model.FileUploadRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "请求参数错误", "error": err.Error()})
		return
	}

	// 校验文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "文件上传失败", "error": err.Error()})
		return
	}
	req.File = file

	// 上传文件
	result, err := h.fileService.UploadFile(req, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "文件上传失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "文件上传成功", "data": result})
}

// GetFileList 获取文件列表
func (h *FileHandler) GetFileList(c *gin.Context) {
	var req model.FileListRequest

	// 获取分页参数
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// 获取其他筛选参数
	req.Page = page
	req.PageSize = pageSize
	req.FileName = c.Query("file_name")
	req.Category = c.Query("category")
	req.FileType = c.Query("file_type")
	req.StartDate = c.Query("start_date")
	req.EndDate = c.Query("end_date")

	// 获取上传者ID
	uploadedByStr := c.Query("uploaded_by")
	if uploadedByStr != "" {
		uploadedBy, err := strconv.ParseUint(uploadedByStr, 10, 32)
		if err == nil {
			req.UploadedBy = uint(uploadedBy)
		}
	}

	// 查询文件列表
	result, err := h.fileService.GetFiles(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "获取文件列表失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "获取文件列表成功", "data": result})
}

// GetFile 获取文件详情
func (h *FileHandler) GetFile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "无效的文件ID"})
		return
	}

	file, err := h.fileService.GetFileByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": "获取文件失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "获取文件成功", "data": file})
}

// UpdateFile 更新文件信息
func (h *FileHandler) UpdateFile(c *gin.Context) {
	var req model.FileUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "请求参数错误", "error": err.Error()})
		return
	}

	file, err := h.fileService.UpdateFile(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "更新文件失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "更新文件成功", "data": file})
}

// DeleteFile 删除文件
func (h *FileHandler) DeleteFile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "无效的文件ID"})
		return
	}

	if err := h.fileService.DeleteFile(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "删除文件失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "删除文件成功"})
}

// BatchDeleteFiles 批量删除文件
func (h *FileHandler) BatchDeleteFiles(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "请求参数错误", "error": err.Error()})
		return
	}

	for _, id := range req.IDs {
		if err := h.fileService.DeleteFile(id); err != nil {
			// 记录错误但继续删除
			continue
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "批量删除文件成功"})
}

// DownloadFile 下载文件
func (h *FileHandler) DownloadFile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "无效的文件ID"})
		return
	}

	file, err := h.fileService.GetFileByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": "获取文件失败", "error": err.Error()})
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(file.FilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": "文件不存在"})
		return
	}

	// 增加下载次数
	go h.fileService.IncrementDownloadCount(uint(id))

	// 下载文件
	c.Header("Content-Disposition", "attachment; filename="+file.FileName)
	c.Header("Content-Type", "application/octet-stream")
	c.File(file.FilePath)
}

// GetFileStats 获取文件统计信息
func (h *FileHandler) GetFileStats(c *gin.Context) {
	stats, err := h.fileService.GetFileStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "获取文件统计失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "获取文件统计成功", "data": stats})
}

// GetCategories 获取文件分类列表
func (h *FileHandler) GetCategories(c *gin.Context) {
	categories, err := h.fileService.GetCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "获取文件分类失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "获取文件分类成功", "data": categories})
}

// GetFileTypes 获取文件类型列表
func (h *FileHandler) GetFileTypes(c *gin.Context) {
	types, err := h.fileService.GetFileTypes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "获取文件类型失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "获取文件类型成功", "data": types})
} 