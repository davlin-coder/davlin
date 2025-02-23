package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	// 设置测试模式
	gin.SetMode(gin.TestMode)

	// 创建一个新的Gin引擎
	r := gin.New()

	// 创建健康检查控制器
	healthCtrl := NewHealthController()

	// 注册路由
	r.GET("/health", healthCtrl.Check)

	// 创建一个测试请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)

	// 执行请求
	r.ServeHTTP(w, req)

	// 断言响应状态码
	assert.Equal(t, http.StatusOK, w.Code)

	// 断言响应内容
	expectedBody := `{"status":"healthy","message":"Service is running normally"}`
	assert.JSONEq(t, expectedBody, w.Body.String())
}