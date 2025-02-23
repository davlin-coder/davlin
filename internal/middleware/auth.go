package middleware

import (
	"net/http"
	"strings"

	"github.com/davlin-coder/davlin/internal/resource/tools"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT认证中间件
func Auth(jwtManager *tools.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取token
		authorization := c.GetHeader("Authorization")
		if authorization == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证令牌"})
			c.Abort()
			return
		}

		// 解析Bearer token
		parts := strings.Split(authorization, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的认证格式"})
			c.Abort()
			return
		}

		// 验证token
		claims, err := jwtManager.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌"})
			c.Abort()
			return
		}

		// 将用户信息存储在上下文中
		c.Set("username", claims.Username)
		c.Set("user_id", claims.Subject)

		c.Next()
	}
}
