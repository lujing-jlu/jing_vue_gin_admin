package service

import (
	"jing_vue_gin_admin/server/config"
	"jing_vue_gin_admin/server/internal/model"
	"time"
)

type LogService struct{}

// NewLogService 创建日志服务
func NewLogService() *LogService {
	return &LogService{}
}

// CreateLog 创建系统日志
func (s *LogService) CreateLog(log *model.SystemLog) error {
	return config.DB.Create(log).Error
}

// GetLogList 获取日志列表
func (s *LogService) GetLogList(req *model.LogListRequest) (*model.LogListResponse, error) {
	var total int64
	var logs []model.SystemLog

	db := config.DB.Model(&model.SystemLog{})

	// 添加查询条件
	if req.Module != "" {
		db = db.Where("module = ?", req.Module)
	}
	if req.Action != "" {
		db = db.Where("action = ?", req.Action)
	}
	if req.Username != "" {
		db = db.Where("username LIKE ?", "%"+req.Username+"%")
	}
	if req.StartTime != "" {
		startTime, err := time.Parse("2006-01-02 15:04:05", req.StartTime)
		if err == nil {
			db = db.Where("created_at >= ?", startTime)
		}
	}
	if req.EndTime != "" {
		endTime, err := time.Parse("2006-01-02 15:04:05", req.EndTime)
		if err == nil {
			db = db.Where("created_at <= ?", endTime)
		}
	}

	// 查询总数
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	// 查询分页数据
	offset := (req.Page - 1) * req.PageSize
	if err := db.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, err
	}

	return &model.LogListResponse{
		Total: total,
		List:  logs,
	}, nil
}

// DeleteLogs 批量删除日志
func (s *LogService) DeleteLogs(ids []uint) error {
	return config.DB.Where("id IN ?", ids).Delete(&model.SystemLog{}).Error
}

// ClearLogs 清空日志
func (s *LogService) ClearLogs() error {
	return config.DB.Exec("DELETE FROM system_logs").Error
}

// AddOperationLog 添加操作日志（快捷方法）
func (s *LogService) AddOperationLog(userID uint, username string, module string, action string, resource string, detail string, ip string, userAgent string, status int) error {
	log := &model.SystemLog{
		UserID:    userID,
		Username:  username,
		Module:    module,
		Action:    action,
		Resource:  resource,
		Detail:    detail,
		IP:        ip,
		UserAgent: userAgent,
		Status:    status,
	}
	return s.CreateLog(log)
} 