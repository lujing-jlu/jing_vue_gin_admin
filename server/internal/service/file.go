package service

import (
	"errors"
	"fmt"
	"io"
	"jing_vue_gin_admin/server/config"
	"jing_vue_gin_admin/server/internal/model"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/gorm"
)

// FileService 文件服务
type FileService struct{}

// NewFileService 创建文件服务实例
func NewFileService() *FileService {
	return &FileService{}
}

// UploadFile 上传文件
func (s *FileService) UploadFile(req model.FileUploadRequest, userID uint) (*model.File, error) {
	// 创建上传目录
	uploadDir := "uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			return nil, fmt.Errorf("创建上传目录失败: %w", err)
		}
	}

	// 获取文件信息
	fileHeader := req.File
	filename := fileHeader.Filename
	fileSize := fileHeader.Size
	fileType := getFileType(filename)

	// 创建唯一文件名
	uniqueFilename := generateUniqueFilename(filename)
	filePath := filepath.Join(uploadDir, uniqueFilename)

	// 保存文件
	srcFile, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer srcFile.Close()

	destFile, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("创建文件失败: %w", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, srcFile); err != nil {
		return nil, fmt.Errorf("保存文件失败: %w", err)
	}

	// 创建文件记录
	category := req.Category
	if category == "" {
		category = "other"
	}

	file := &model.File{
		FileName:    filename,
		FileSize:    fileSize,
		FileType:    fileType,
		FilePath:    filePath,
		Category:    category,
		Description: req.Description,
		UploadedBy:  userID,
	}

	if err := config.DB.Create(file).Error; err != nil {
		// 删除物理文件
		os.Remove(filePath)
		return nil, fmt.Errorf("保存文件记录失败: %w", err)
	}

	return file, nil
}

// GetFiles 获取文件列表
func (s *FileService) GetFiles(req model.FileListRequest) (*model.FileListResponse, error) {
	var files []model.File
	var total int64
	
	query := config.DB.Model(&model.File{})

	// 应用过滤条件
	if req.FileName != "" {
		query = query.Where("file_name LIKE ?", "%"+req.FileName+"%")
	}
	if req.Category != "" {
		query = query.Where("category = ?", req.Category)
	}
	if req.FileType != "" {
		query = query.Where("file_type = ?", req.FileType)
	}
	if req.UploadedBy != 0 {
		query = query.Where("uploaded_by = ?", req.UploadedBy)
	}
	if req.StartDate != "" && req.EndDate != "" {
		query = query.Where("created_at BETWEEN ? AND ?", req.StartDate, req.EndDate+" 23:59:59")
	} else if req.StartDate != "" {
		query = query.Where("created_at >= ?", req.StartDate)
	} else if req.EndDate != "" {
		query = query.Where("created_at <= ?", req.EndDate+" 23:59:59")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("获取文件总数失败: %w", err)
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&files).Error; err != nil {
		return nil, fmt.Errorf("获取文件列表失败: %w", err)
	}

	return &model.FileListResponse{
		Total: int(total),
		List:  files,
	}, nil
}

// GetFileByID 通过ID获取文件信息
func (s *FileService) GetFileByID(id uint) (*model.File, error) {
	var file model.File
	if err := config.DB.First(&file, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("文件不存在")
		}
		return nil, fmt.Errorf("获取文件信息失败: %w", err)
	}
	return &file, nil
}

// UpdateFile 更新文件信息
func (s *FileService) UpdateFile(req model.FileUpdateRequest) (*model.File, error) {
	file, err := s.GetFileByID(req.ID)
	if err != nil {
		return nil, err
	}

	// 更新文件信息
	updates := map[string]interface{}{}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.FileName != "" && req.FileName != file.FileName {
		updates["file_name"] = req.FileName
	}

	if len(updates) > 0 {
		if err := config.DB.Model(file).Updates(updates).Error; err != nil {
			return nil, fmt.Errorf("更新文件信息失败: %w", err)
		}
	}

	return file, nil
}

// DeleteFile 删除文件
func (s *FileService) DeleteFile(id uint) error {
	file, err := s.GetFileByID(id)
	if err != nil {
		return err
	}

	// 开启事务
	return config.DB.Transaction(func(tx *gorm.DB) error {
		// 删除数据库记录
		if err := tx.Delete(file).Error; err != nil {
			return fmt.Errorf("删除文件记录失败: %w", err)
		}

		// 删除物理文件
		if err := os.Remove(file.FilePath); err != nil {
			// 记录错误但不中断事务
			fmt.Printf("删除物理文件失败: %v\n", err)
		}

		return nil
	})
}

// IncrementDownloadCount 增加下载次数
func (s *FileService) IncrementDownloadCount(id uint) error {
	if err := config.DB.Model(&model.File{}).Where("id = ?", id).
		UpdateColumn("downloads", gorm.Expr("downloads + ?", 1)).Error; err != nil {
		return fmt.Errorf("更新下载次数失败: %w", err)
	}
	return nil
}

// GetFileStats 获取文件统计信息
func (s *FileService) GetFileStats() (*model.FileStats, error) {
	var stats model.FileStats
	var totalCount int64
	var totalSize int64

	// 获取文件总数和总大小
	if err := config.DB.Model(&model.File{}).Count(&totalCount).Error; err != nil {
		return nil, fmt.Errorf("获取文件总数失败: %w", err)
	}
	if err := config.DB.Model(&model.File{}).Select("COALESCE(SUM(file_size), 0)").Scan(&totalSize).Error; err != nil {
		return nil, fmt.Errorf("获取文件总大小失败: %w", err)
	}

	stats.TotalFiles = int(totalCount)
	stats.TotalSize = totalSize

	// 获取分类统计
	var categories []model.FileCategoryCount
	if err := config.DB.Model(&model.File{}).
		Select("category, COUNT(*) as count").
		Group("category").
		Scan(&categories).Error; err != nil {
		return nil, fmt.Errorf("获取分类统计失败: %w", err)
	}
	stats.Categories = categories

	// 获取文件类型统计
	var types []model.FileTypeCount
	if err := config.DB.Model(&model.File{}).
		Select("file_type, COUNT(*) as count").
		Group("file_type").
		Scan(&types).Error; err != nil {
		return nil, fmt.Errorf("获取文件类型统计失败: %w", err)
	}
	stats.Types = types

	// 获取最近上传
	var recentUploads []model.File
	if err := config.DB.Model(&model.File{}).
		Order("created_at DESC").
		Limit(5).
		Find(&recentUploads).Error; err != nil {
		return nil, fmt.Errorf("获取最近上传失败: %w", err)
	}
	stats.RecentUploads = recentUploads

	// 获取热门下载
	var popularDownloads []model.File
	if err := config.DB.Model(&model.File{}).
		Order("downloads DESC").
		Limit(5).
		Find(&popularDownloads).Error; err != nil {
		return nil, fmt.Errorf("获取热门下载失败: %w", err)
	}
	stats.PopularDownloads = popularDownloads

	return &stats, nil
}

// GetCategories 获取所有文件分类
func (s *FileService) GetCategories() ([]string, error) {
	var categories []string
	if err := config.DB.Model(&model.File{}).
		Distinct().
		Pluck("category", &categories).Error; err != nil {
		return nil, fmt.Errorf("获取分类列表失败: %w", err)
	}
	return categories, nil
}

// GetFileTypes 获取所有文件类型
func (s *FileService) GetFileTypes() ([]string, error) {
	var types []string
	if err := config.DB.Model(&model.File{}).
		Distinct().
		Pluck("file_type", &types).Error; err != nil {
		return nil, fmt.Errorf("获取文件类型列表失败: %w", err)
	}
	return types, nil
}

// 辅助函数

// getFileType 根据文件名获取文件类型
func getFileType(filename string) string {
	ext := filepath.Ext(filename)
	if ext == "" {
		return "unknown"
	}
	return strings.ToLower(ext[1:]) // 去掉开头的点
}

// generateUniqueFilename 生成唯一的文件名
func generateUniqueFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	baseName := filepath.Base(originalFilename)
	baseName = strings.TrimSuffix(baseName, ext)
	
	// 使用时间戳和随机数生成唯一名称
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%s_%d%s", baseName, timestamp, ext)
} 