package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"jing_vue_gin_admin/server/internal/model"
	"jing_vue_gin_admin/server/internal/service"
	"strings"
)

// bodyLogWriter 自定义响应写入器，用于捕获响应内容
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// OperationLog 操作日志中间件
func OperationLog(module string, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求体
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 包装响应写入器
		blw := &bodyLogWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = blw

		// 继续处理请求
		c.Next()

		// 获取请求信息
		userID, _ := c.Get("userID")
		username, _ := c.Get("username")
		
		// 获取客户端IP和用户代理
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()
		
		// 获取资源标识
		resource := c.FullPath()
		
		// 构建详情
		detail := ""
		if len(requestBody) > 0 {
			if len(requestBody) > 200 {
				detail = string(requestBody[:200]) + "..."
			} else {
				detail = string(requestBody)
			}
		}
		
		// 状态判断
		status := 1
		if c.Writer.Status() >= 400 {
			status = 0
		}
		
		// 创建日志服务
		logService := service.NewLogService()
		
		// 添加操作日志
		var userIDUint uint
		if userID != nil {
			userIDUint = userID.(uint)
		}
		
		var usernameStr string
		if username != nil {
			usernameStr = username.(string)
		}
		
		// 记录日志
		logService.AddOperationLog(
			userIDUint,
			usernameStr,
			module,
			action,
			resource,
			detail,
			clientIP,
			userAgent,
			status,
		)
	}
}

// LoginLogMiddleware 登录日志中间件
func LoginLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果不是登录接口，则跳过
		if !strings.Contains(c.FullPath(), "/login") {
			c.Next()
			return
		}
		
		// 获取请求体
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}
		
		// 解析用户名
		username := ""
		if len(requestBody) > 0 {
			var loginReq model.LoginRequest
			if err := json.Unmarshal(requestBody, &loginReq); err == nil {
				username = loginReq.Username
			}
		}
		
		// 包装响应写入器
		blw := &bodyLogWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = blw
		
		// 继续处理请求
		c.Next()
		
		// 获取客户端IP和用户代理
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()
		
		// 获取资源标识
		resource := c.FullPath()
		
		// 状态判断
		status := 1
		if c.Writer.Status() >= 400 {
			status = 0
		}
		
		// 创建日志服务
		logService := service.NewLogService()
		
		// 添加登录日志
		logService.AddOperationLog(
			0, // 未登录用户ID为0
			username,
			model.LogModuleAuth,
			model.LogActionLogin,
			resource,
			"",
			clientIP,
			userAgent,
			status,
		)
	}
} 