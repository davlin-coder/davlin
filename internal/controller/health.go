package controller

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// HealthController defines the health check controller interface
type HealthController interface {
	Check(c *gin.Context)
}

// healthController implements the HealthController interface
type healthController struct {}

// NewHealthController creates a new health check controller instance
func NewHealthController() HealthController {
	return &healthController{}
}

// Check performs health status check
func (ctrl *healthController) Check(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"message": "Service is running normally",
	})
}