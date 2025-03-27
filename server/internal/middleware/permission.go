package middleware

import (
	"github.com/gin-gonic/gin"
	"jing_vue_gin_admin/server/config"
	"jing_vue_gin_admin/server/internal/model"
	"net/http"
)

// RequirePermission 检查用户是否具有指定权限
func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文中获取用户ID
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
			c.Abort()
			return
		}

		// 获取用户信息
		var user model.User
		if err := config.DB.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
			c.Abort()
			return
		}

		// 获取用户角色
		var role model.Role
		if err := config.DB.Where("name = ?", user.Role).First(&role).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "角色不存在"})
			c.Abort()
			return
		}

		// 检查角色状态
		if role.Status != 1 {
			c.JSON(http.StatusForbidden, gin.H{"error": "角色已禁用"})
			c.Abort()
			return
		}

		// 检查是否具有所需权限
		hasPermission := false
		for _, perm := range role.Permissions {
			if perm == permission {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "没有操作权限"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAdmin 检查用户是否为管理员
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文中获取用户ID
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
			c.Abort()
			return
		}

		// 获取用户信息
		var user model.User
		if err := config.DB.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
			c.Abort()
			return
		}

		// 检查是否为管理员
		if user.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "需要管理员权限"})
			c.Abort()
			return
		}

		c.Next()
	}
} 