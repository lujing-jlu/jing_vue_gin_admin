package main

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"jing_vue_gin_admin/server/config"
	"jing_vue_gin_admin/server/internal/handler"
	"jing_vue_gin_admin/server/internal/middleware"
	"jing_vue_gin_admin/server/internal/model"
	"log"
)

func main() {
	// 初始化数据库
	if err := config.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// 创建测试用户和角色
	var count int64
	config.DB.Model(&model.User{}).Count(&count)
	if count == 0 {
		// 先创建管理员角色
		adminRole := &model.Role{
			Name:        "admin",
			Description: "系统管理员",
			Status:      1,
			Permissions: []string{
				model.PermissionUserView,
				model.PermissionUserCreate,
				model.PermissionUserEdit,
				model.PermissionUserDelete,
				model.PermissionRoleView,
				model.PermissionRoleCreate,
				model.PermissionRoleEdit,
				model.PermissionRoleDelete,
				model.PermissionSystemConfig,
				model.PermissionSystemLog,
			},
		}
		if err := config.DB.Create(adminRole).Error; err != nil {
			log.Fatal("Failed to create admin role:", err)
		}
		log.Println("Created admin role")

		// 再创建管理员用户
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("Failed to hash password:", err)
		}

		user := &model.User{
			Username: "admin",
			Password: string(hashedPassword),
			Nickname: "管理员",
			Role:     adminRole.Name,
			Status:   1,
		}
		if err := config.DB.Create(user).Error; err != nil {
			log.Fatal("Failed to create test user:", err)
		}
		log.Println("Created test user: admin/123456")
	}

	r := gin.Default()

	// 允许跨域
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 创建处理器
	userHandler := handler.NewUserHandler()
	roleHandler := handler.NewRoleHandler()
	logHandler := handler.NewLogHandler()

	// 公开路由
	public := r.Group("/api")
	public.Use(middleware.LoginLogMiddleware())
	{
		public.POST("/login", userHandler.Login)
	}

	// 需要认证的路由
	authorized := r.Group("/api")
	authorized.Use(middleware.AuthMiddleware())
	{
		// 用户管理API
		authorized.GET("/users", middleware.RequirePermission(model.PermissionUserView), middleware.OperationLog(model.LogModuleUser, model.LogActionView), userHandler.GetUserList)
		authorized.POST("/users", middleware.RequirePermission(model.PermissionUserCreate), middleware.OperationLog(model.LogModuleUser, model.LogActionCreate), userHandler.CreateUser)
		authorized.PUT("/users/:id", middleware.RequirePermission(model.PermissionUserEdit), middleware.OperationLog(model.LogModuleUser, model.LogActionUpdate), userHandler.UpdateUser)
		authorized.PUT("/users/:id/status", middleware.RequirePermission(model.PermissionUserEdit), middleware.OperationLog(model.LogModuleUser, model.LogActionUpdate), userHandler.ToggleUserStatus)

		// 个人信息API（无需特殊权限）
		authorized.PUT("/user/profile", middleware.OperationLog(model.LogModuleUser, model.LogActionUpdate), userHandler.UpdateProfile)
		authorized.PUT("/user/password", middleware.OperationLog(model.LogModuleUser, model.LogActionUpdate), userHandler.ChangePassword)

		// 角色管理API
		authorized.GET("/roles", middleware.RequirePermission(model.PermissionRoleView), middleware.OperationLog(model.LogModuleRole, model.LogActionView), roleHandler.GetRoleList)
		authorized.POST("/roles", middleware.RequirePermission(model.PermissionRoleCreate), middleware.OperationLog(model.LogModuleRole, model.LogActionCreate), roleHandler.CreateRole)
		authorized.PUT("/roles/:id", middleware.RequirePermission(model.PermissionRoleEdit), middleware.OperationLog(model.LogModuleRole, model.LogActionUpdate), roleHandler.UpdateRole)
		authorized.DELETE("/roles/:id", middleware.RequirePermission(model.PermissionRoleDelete), middleware.OperationLog(model.LogModuleRole, model.LogActionDelete), roleHandler.DeleteRole)
		authorized.PUT("/roles/:id/status", middleware.RequirePermission(model.PermissionRoleEdit), middleware.OperationLog(model.LogModuleRole, model.LogActionUpdate), roleHandler.ToggleRoleStatus)
		authorized.GET("/permissions", middleware.RequirePermission(model.PermissionRoleView), roleHandler.GetAllPermissions)

		// 系统日志API
		authorized.GET("/logs", middleware.RequirePermission(model.PermissionSystemLog), middleware.OperationLog(model.LogModuleSystem, model.LogActionView), logHandler.GetLogList)
		authorized.DELETE("/logs", middleware.RequirePermission(model.PermissionSystemLog), middleware.OperationLog(model.LogModuleSystem, model.LogActionDelete), logHandler.DeleteLogs)
		authorized.DELETE("/logs/clear", middleware.RequirePermission(model.PermissionSystemLog), middleware.OperationLog(model.LogModuleSystem, model.LogActionDelete), logHandler.ClearLogs)
		authorized.GET("/logs/modules", middleware.RequirePermission(model.PermissionSystemLog), logHandler.GetLogModules)
		authorized.GET("/logs/actions", middleware.RequirePermission(model.PermissionSystemLog), logHandler.GetLogActions)
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
