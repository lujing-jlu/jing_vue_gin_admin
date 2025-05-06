package model

import "mime/multipart"

// File 文件模型
type File struct {
	Base
	FileName    string `json:"file_name" gorm:"size:255;not null"`     // 文件名
	FileSize    int64  `json:"file_size" gorm:"not null"`              // 文件大小
	FileType    string `json:"file_type" gorm:"size:50;not null"`      // 文件类型
	FilePath    string `json:"file_path" gorm:"size:255;not null"`     // 文件路径
	Category    string `json:"category" gorm:"size:50;default:'other'"`// 文件分类
	Description string `json:"description" gorm:"size:500"`            // 文件描述
	UploadedBy  uint   `json:"uploaded_by" gorm:"not null"`            // 上传者ID
	Downloads   int    `json:"downloads" gorm:"default:0"`             // 下载次数
}

// FileUploadRequest 文件上传请求
type FileUploadRequest struct {
	File        *multipart.FileHeader `form:"file" binding:"required"`
	Category    string                `form:"category"`
	Description string                `form:"description"`
}

// FileListRequest 文件列表请求
type FileListRequest struct {
	FileName    string `form:"file_name"`
	Category    string `form:"category"`
	FileType    string `form:"file_type"`
	StartDate   string `form:"start_date"`
	EndDate     string `form:"end_date"`
	UploadedBy  uint   `form:"uploaded_by"`
	Page        int    `form:"page" binding:"required,min=1"`
	PageSize    int    `form:"page_size" binding:"required,min=1,max=100"`
}

// FileListResponse 文件列表响应
type FileListResponse struct {
	Total int    `json:"total"`
	List  []File `json:"list"`
}

// FileUpdateRequest 文件更新请求
type FileUpdateRequest struct {
	ID          uint   `json:"id" binding:"required"`
	Category    string `json:"category"`
	Description string `json:"description"`
	FileName    string `json:"file_name"`
}

// FileCategoryCount 文件分类统计
type FileCategoryCount struct {
	Category string `json:"category"`
	Count    int    `json:"count"`
}

// FileTypeCount 文件类型统计
type FileTypeCount struct {
	FileType string `json:"file_type"`
	Count    int    `json:"count"`
}

// FileStats 文件统计
type FileStats struct {
	TotalFiles      int                `json:"total_files"`
	TotalSize       int64              `json:"total_size"`
	Categories      []FileCategoryCount `json:"categories"`
	Types           []FileTypeCount     `json:"types"`
	RecentUploads   []File             `json:"recent_uploads"`
	PopularDownloads []File            `json:"popular_downloads"`
} 