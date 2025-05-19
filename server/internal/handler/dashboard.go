package handler

import (
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"jing_vue_gin_admin/server/config"
	"jing_vue_gin_admin/server/internal/model"
)

// 仪表盘数据结构
type DashboardStats struct {
	UserCount         int                       `json:"userCount"`
	RoleCount         int                       `json:"roleCount"`
	LogCount          int                       `json:"logCount"`
	FileCount         int                       `json:"fileCount"`
	TodayLoginCount   int                       `json:"todayLoginCount"`
	StorageUsage      int64                     `json:"storageUsage"`
	UserActivities    []UserActivityData        `json:"userActivities"`
	FileTypeDistData  []TypeDistributionData    `json:"fileTypeDistribution"`
	LogTypeDistData   []TypeDistributionData    `json:"logTypeDistribution"`
	RecentRegData     []DateCountData           `json:"recentRegistrations"`
}

// 日期-数量数据结构
type DateCountData struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// 类型分布数据结构
type TypeDistributionData struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

// 用户活跃度数据结构
type UserActivityData struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// 系统信息数据结构
type SystemInfo struct {
	Version      string `json:"version"`
	Environment  string `json:"environment"`
	ServerTime   string `json:"serverTime"`
	GoVersion    string `json:"goVersion"`
	GormVersion  string `json:"gormVersion"`
	GinVersion   string `json:"ginVersion"`
	Uptime       string `json:"uptime"`
	MemoryUsage  string `json:"memoryUsage"`
	CPUUsage     string `json:"cpuUsage"`
}

// DashboardHandler 处理仪表盘相关的请求
type DashboardHandler struct{}

// NewDashboardHandler 创建仪表盘处理器
func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{}
}

// GetDashboardStats 获取仪表盘统计数据
func (h *DashboardHandler) GetDashboardStats(c *gin.Context) {
	var userCount, roleCount, logCount, fileCount, todayLoginCount int64
	var storageUsage int64

	// 获取用户总数
	config.DB.Model(&model.User{}).Count(&userCount)

	// 获取角色总数
	config.DB.Model(&model.Role{}).Count(&roleCount)

	// 获取日志总数
	config.DB.Model(&model.SystemLog{}).Count(&logCount)

	// 获取文件总数
	config.DB.Model(&model.File{}).Count(&fileCount)

	// 获取今日登录总数
	today := time.Now().Format("2006-01-02")
	config.DB.Model(&model.SystemLog{}).
		Where("module = ? AND action = ? AND created_at LIKE ? AND status = ?", "auth", "login", today+"%", 1).
		Count(&todayLoginCount)

	// 获取存储使用量
	config.DB.Model(&model.File{}).Select("COALESCE(SUM(file_size), 0)").Row().Scan(&storageUsage)

	// 获取用户活跃度数据
	userActivities := h.getUserActivities()

	// 获取文件类型分布
	fileTypeDist := h.getFileTypeDistribution()

	// 获取日志类型分布
	logTypeDist := h.getLogTypeDistribution()

	// 获取最近注册趋势
	recentReg := h.getRecentRegistrations()

	stats := DashboardStats{
		UserCount:        int(userCount),
		RoleCount:        int(roleCount),
		LogCount:         int(logCount),
		FileCount:        int(fileCount),
		TodayLoginCount:  int(todayLoginCount),
		StorageUsage:     storageUsage,
		UserActivities:   userActivities,
		FileTypeDistData: fileTypeDist,
		LogTypeDistData:  logTypeDist,
		RecentRegData:    recentReg,
	}

	c.JSON(http.StatusOK, stats)
}

// GetSystemInfo 获取系统信息
func (h *DashboardHandler) GetSystemInfo(c *gin.Context) {
	// 获取启动时间，这里使用当前时间代替
	systemStartTime := time.Now().Add(-72 * time.Hour) // 假设系统运行了72小时
	uptime := time.Since(systemStartTime).String()

	// 简化版本的系统信息，避免使用gopsutil
	info := SystemInfo{
		Version:     "1.0.0",
		Environment: "production",
		ServerTime:  time.Now().Format("2006-01-02 15:04:05"),
		GoVersion:   runtime.Version(),
		GormVersion: "v2.0",
		GinVersion:  "v1.9.1",
		Uptime:      uptime,
		MemoryUsage: "30%", // 硬编码模拟值
		CPUUsage:    "10%", // 硬编码模拟值
	}

	c.JSON(http.StatusOK, info)
}

// getUserActivities 获取用户活跃度数据（最近7天）
func (h *DashboardHandler) getUserActivities() []UserActivityData {
	var result []UserActivityData
	
	// 最近7天的日期
	for i := 6; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		var count int64
		
		// 查询当天的登录日志数量
		config.DB.Model(&model.SystemLog{}).
			Where("module = ? AND action = ? AND created_at LIKE ? AND status = ?", "auth", "login", date+"%", 1).
			Count(&count)
		
		result = append(result, UserActivityData{
			Date:  date,
			Count: int(count),
		})
	}
	
	return result
}

// getFileTypeDistribution 获取文件类型分布
func (h *DashboardHandler) getFileTypeDistribution() []TypeDistributionData {
	var results []struct {
		FileType string
		Count    int
	}
	
	config.DB.Model(&model.File{}).
		Select("file_type as FileType, COUNT(*) as Count").
		Group("file_type").
		Order("Count DESC").
		Limit(5).
		Scan(&results)
	
	var data []TypeDistributionData
	for _, r := range results {
		data = append(data, TypeDistributionData{
			Type:  r.FileType,
			Count: r.Count,
		})
	}
	
	return data
}

// getLogTypeDistribution 获取日志类型分布
func (h *DashboardHandler) getLogTypeDistribution() []TypeDistributionData {
	var results []struct {
		Module string
		Count  int
	}
	
	config.DB.Model(&model.SystemLog{}).
		Select("module as Module, COUNT(*) as Count").
		Group("module").
		Order("Count DESC").
		Limit(5).
		Scan(&results)
	
	var data []TypeDistributionData
	for _, r := range results {
		data = append(data, TypeDistributionData{
			Type:  r.Module,
			Count: r.Count,
		})
	}
	
	return data
}

// getRecentRegistrations 获取最近注册用户趋势（最近30天）
func (h *DashboardHandler) getRecentRegistrations() []DateCountData {
	var result []DateCountData
	
	// 最近30天的日期
	for i := 29; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		var count int64
		
		// 查询当天注册的用户数量
		config.DB.Model(&model.User{}).
			Where("created_at LIKE ?", date+"%").
			Count(&count)
		
		result = append(result, DateCountData{
			Date:  date,
			Count: int(count),
		})
	}
	
	return result
}

// GetUserActivities 获取用户活跃度数据API
func (h *DashboardHandler) GetUserActivities(c *gin.Context) {
	c.JSON(http.StatusOK, h.getUserActivities())
}

// GetFileTypeDistribution 获取文件类型分布API
func (h *DashboardHandler) GetFileTypeDistribution(c *gin.Context) {
	c.JSON(http.StatusOK, h.getFileTypeDistribution())
}

// GetLogTypeDistribution 获取日志类型分布API
func (h *DashboardHandler) GetLogTypeDistribution(c *gin.Context) {
	c.JSON(http.StatusOK, h.getLogTypeDistribution())
}

// GetRecentRegistrations 获取最近注册用户趋势API
func (h *DashboardHandler) GetRecentRegistrations(c *gin.Context) {
	c.JSON(http.StatusOK, h.getRecentRegistrations())
} 