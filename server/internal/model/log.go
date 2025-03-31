package model

// SystemLog 系统日志模型
type SystemLog struct {
	Base
	UserID    uint   `json:"user_id"`               // 操作用户ID
	Username  string `gorm:"size:32" json:"username"` // 操作用户名
	Module    string `gorm:"size:32" json:"module"`   // 操作模块
	Action    string `gorm:"size:32" json:"action"`   // 操作类型
	Resource  string `gorm:"size:64" json:"resource"` // 操作资源
	Detail    string `gorm:"size:255" json:"detail"`  // 操作详情
	IP        string `gorm:"size:32" json:"ip"`       // 操作IP
	UserAgent string `gorm:"size:255" json:"user_agent"` // 用户代理
	Status    int    `gorm:"default:1" json:"status"` // 操作状态：1成功，0失败
}

// LogListRequest 日志列表请求参数
type LogListRequest struct {
	Module    string `form:"module" json:"module"`
	Action    string `form:"action" json:"action"`
	Username  string `form:"username" json:"username"`
	StartTime string `form:"start_time" json:"start_time"`
	EndTime   string `form:"end_time" json:"end_time"`
	Page      int    `form:"page" json:"page" binding:"required,min=1"`
	PageSize  int    `form:"page_size" json:"page_size" binding:"required,min=1,max=100"`
}

// LogListResponse 日志列表响应
type LogListResponse struct {
	Total int64       `json:"total"`
	List  []SystemLog `json:"list"`
}

// 日志操作模块常量
const (
	LogModuleUser   = "user"   // 用户模块
	LogModuleRole   = "role"   // 角色模块
	LogModuleSystem = "system" // 系统模块
	LogModuleAuth   = "auth"   // 认证模块
)

// 日志操作类型常量
const (
	LogActionLogin  = "login"  // 登录
	LogActionLogout = "logout" // 登出
	LogActionCreate = "create" // 创建
	LogActionUpdate = "update" // 更新
	LogActionDelete = "delete" // 删除
	LogActionView   = "view"   // 查看
	LogActionExport = "export" // 导出
	LogActionImport = "import" // 导入
	LogActionEnable = "enable" // 启用
	LogActionDisable = "disable" // 禁用
) 