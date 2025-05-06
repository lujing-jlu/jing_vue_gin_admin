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
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 创建默认管理员角色和用户
	createDefaultAdminRoleAndUser()

	// 初始化Gin框架
	r := gin.Default()

	// 允许跨域
	r.Use(middleware.Cors())

	// 创建各种处理器实例
	userHandler := handler.NewUserHandler()
	roleHandler := handler.NewRoleHandler()
	logHandler := handler.NewLogHandler()
	fileHandler := handler.NewFileHandler() // 添加文件处理器

	// API路由组
	api := r.Group("/api")

	// 公共路由组
	public := api.Group("")
	public.Use(middleware.LoginLogMiddleware())
	{
		public.POST("/login", userHandler.Login)
	}

	// 需要身份验证的路由组
	auth := api.Group("")
	auth.Use(middleware.AuthMiddleware())
	{
		// 用户相关路由
		userRoutes := auth.Group("/users")
		{
			userRoutes.GET("", middleware.RequirePermission("user:view"), middleware.OperationLog("user", "view"), userHandler.GetUsers)
			userRoutes.GET("/:id", middleware.RequirePermission("user:view"), middleware.OperationLog("user", "view"), userHandler.GetUser)
			userRoutes.POST("", middleware.RequirePermission("user:create"), middleware.OperationLog("user", "create"), userHandler.CreateUser)
			userRoutes.PUT("/:id", middleware.RequirePermission("user:update"), middleware.OperationLog("user", "update"), userHandler.UpdateUser)
			userRoutes.DELETE("/:id", middleware.RequirePermission("user:delete"), middleware.OperationLog("user", "delete"), userHandler.DeleteUser)
			userRoutes.PUT("/:id/status", middleware.RequirePermission("user:update"), middleware.OperationLog("user", "update"), userHandler.UpdateUserStatus)
			userRoutes.GET("/current", userHandler.GetCurrentUser)
		}

		// 角色相关路由
		roleRoutes := auth.Group("/roles")
		{
			roleRoutes.GET("", middleware.RequirePermission("role:view"), middleware.OperationLog("role", "view"), roleHandler.GetRoles)
			roleRoutes.GET("/:id", middleware.RequirePermission("role:view"), middleware.OperationLog("role", "view"), roleHandler.GetRole)
			roleRoutes.POST("", middleware.RequirePermission("role:create"), middleware.OperationLog("role", "create"), roleHandler.CreateRole)
			roleRoutes.PUT("/:id", middleware.RequirePermission("role:update"), middleware.OperationLog("role", "update"), roleHandler.UpdateRole)
			roleRoutes.DELETE("/:id", middleware.RequirePermission("role:delete"), middleware.OperationLog("role", "delete"), roleHandler.DeleteRole)
			roleRoutes.GET("/all/permissions", middleware.RequirePermission("role:view"), roleHandler.GetAllPermissions)
		}

		// 日志相关路由
		logRoutes := auth.Group("/logs")
		{
			logRoutes.GET("", middleware.RequirePermission("log:view"), middleware.OperationLog("system", "view"), logHandler.GetLogList)
			logRoutes.DELETE("", middleware.RequirePermission("log:delete"), middleware.OperationLog("system", "delete"), logHandler.DeleteLogs)
			logRoutes.DELETE("/clear", middleware.RequirePermission("log:delete"), middleware.OperationLog("system", "delete"), logHandler.ClearLogs)
			logRoutes.GET("/modules", middleware.RequirePermission("log:view"), logHandler.GetLogModules)
			logRoutes.GET("/actions", middleware.RequirePermission("log:view"), logHandler.GetLogActions)
		}

		// 文件相关路由
		fileRoutes := auth.Group("/files")
		{
			fileRoutes.POST("", middleware.RequirePermission("file:upload"), middleware.OperationLog("system", "create"), fileHandler.UploadFile)
			fileRoutes.GET("", middleware.RequirePermission("file:view"), middleware.OperationLog("system", "view"), fileHandler.GetFileList)
			fileRoutes.GET("/:id", middleware.RequirePermission("file:view"), middleware.OperationLog("system", "view"), fileHandler.GetFile)
			fileRoutes.PUT("", middleware.RequirePermission("file:update"), middleware.OperationLog("system", "update"), fileHandler.UpdateFile)
			fileRoutes.DELETE("/:id", middleware.RequirePermission("file:delete"), middleware.OperationLog("system", "delete"), fileHandler.DeleteFile)
			fileRoutes.POST("/batch/delete", middleware.RequirePermission("file:delete"), middleware.OperationLog("system", "delete"), fileHandler.BatchDeleteFiles)
			fileRoutes.GET("/download/:id", middleware.OperationLog("system", "view"), fileHandler.DownloadFile)
			fileRoutes.GET("/stats", middleware.RequirePermission("file:view"), fileHandler.GetFileStats)
			fileRoutes.GET("/categories", middleware.RequirePermission("file:view"), fileHandler.GetCategories)
			fileRoutes.GET("/types", middleware.RequirePermission("file:view"), fileHandler.GetFileTypes)
		}
	}

	// 启动服务器
	r.Run(":8080")
}

func createDefaultAdminRoleAndUser() {
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
				model.PermissionLogDelete,
				model.PermissionFileView,
				model.PermissionFileUpload,
				model.PermissionFileUpdate,
				model.PermissionFileDelete,
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
}
