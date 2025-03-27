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
		authorized.GET("/users", userHandler.GetUserList)          // 获取用户列表
		authorized.POST("/users", userHandler.CreateUser)          // 创建用户
		authorized.PUT("/users/:id", userHandler.UpdateUser)       // 更新用户
		authorized.PUT("/users/:id/status", userHandler.ToggleUserStatus) // 切换用户状态

		// 角色管理API
		authorized.GET("/roles", roleHandler.GetRoleList)           // 获取角色列表
		authorized.POST("/roles", roleHandler.CreateRole)           // 创建角色
		authorized.PUT("/roles/:id", roleHandler.UpdateRole)        // 更新角色
		authorized.DELETE("/roles/:id", roleHandler.DeleteRole)     // 删除角色
		authorized.PUT("/roles/:id/status", roleHandler.ToggleRoleStatus) // 切换角色状态
		authorized.GET("/permissions", roleHandler.GetAllPermissions)      // 获取所有权限
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
