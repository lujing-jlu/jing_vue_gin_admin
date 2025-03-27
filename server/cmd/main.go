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

	// 创建测试用户
	var count int64
	config.DB.Model(&model.User{}).Count(&count)
	if count == 0 {
		// 使用bcrypt正确哈希密码
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("Failed to hash password:", err)
		}

		user := &model.User{
			Username: "admin",
			Password: string(hashedPassword),
			Nickname: "管理员",
			Role:     "admin",
			Status:   1,
		}
		if err := config.DB.Create(user).Error; err != nil {
			log.Fatal("Failed to create test user:", err)
		}
		log.Println("Created test user: admin/123456")

		// 创建管理员角色
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
			},
		}
		if err := config.DB.Create(adminRole).Error; err != nil {
			log.Fatal("Failed to create admin role:", err)
		}
		log.Println("Created admin role")
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

	// 公开路由
	public := r.Group("/api")
	{
		public.POST("/login", userHandler.Login)
	}

	// 需要认证的路由
	authorized := r.Group("/api")
	authorized.Use(middleware.AuthMiddleware())
	{
		// 用户管理API
		authorized.GET("/users", middleware.RequirePermission(model.PermissionUserView), userHandler.GetUserList)
		authorized.POST("/users", middleware.RequirePermission(model.PermissionUserCreate), userHandler.CreateUser)
		authorized.PUT("/users/:id", middleware.RequirePermission(model.PermissionUserEdit), userHandler.UpdateUser)
		authorized.PUT("/users/:id/status", middleware.RequirePermission(model.PermissionUserEdit), userHandler.ToggleUserStatus)

		// 角色管理API
		authorized.GET("/roles", middleware.RequirePermission(model.PermissionRoleView), roleHandler.GetRoleList)
		authorized.POST("/roles", middleware.RequirePermission(model.PermissionRoleCreate), roleHandler.CreateRole)
		authorized.PUT("/roles/:id", middleware.RequirePermission(model.PermissionRoleEdit), roleHandler.UpdateRole)
		authorized.DELETE("/roles/:id", middleware.RequirePermission(model.PermissionRoleDelete), roleHandler.DeleteRole)
		authorized.PUT("/roles/:id/status", middleware.RequirePermission(model.PermissionRoleEdit), roleHandler.ToggleRoleStatus)
		authorized.GET("/permissions", middleware.RequirePermission(model.PermissionRoleView), roleHandler.GetAllPermissions)
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
