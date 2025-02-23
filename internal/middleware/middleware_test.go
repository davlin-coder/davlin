package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCors(t *testing.T) {
	// 设置测试模式
	gin.SetMode(gin.TestMode)

	// 创建一个新的Gin引擎
	r := gin.New()
	r.Use(Cors())

	// 添加测试路由
	r.GET("/test", func(c *gin.Context) {
		c.String(200, "test")
	})

	// 测试OPTIONS请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "/test", nil)
	r.ServeHTTP(w, req)

	// 验证OPTIONS请求的响应
	assert.Equal(t, 204, w.Code)
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "Content-Type, Authorization", w.Header().Get("Access-Control-Allow-Headers"))

	// 测试正常GET请求
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	// 验证GET请求的响应
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestLogger(t *testing.T) {
	// 设置测试模式
	gin.SetMode(gin.TestMode)

	// 创建一个新的Gin引擎
	r := gin.New()
	r.Use(Logger())

	// 添加测试路由
	r.GET("/test", func(c *gin.Context) {
		time.Sleep(10 * time.Millisecond) // 模拟处理时间
		c.String(200, "test")
	})

	// 发送测试请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "test", w.Body.String())
}