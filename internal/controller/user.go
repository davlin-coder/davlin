package controller

import (
	"net/http"

	"github.com/davlin-coder/davlin/internal/model"
	"github.com/davlin-coder/davlin/internal/service"
	"github.com/gin-gonic/gin"
)

// UserController defines the user controller interface
type UserController interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	SendVerificationCode(c *gin.Context)
}

// userController implements UserController interface
type userController struct {
	userService         service.UserService
	verificationService service.VerificationService
}

// NewUserController creates a new user controller instance
func NewUserController(userService service.UserService, verificationService service.VerificationService) UserController {
	return &userController{
		userService:         userService,
		verificationService: verificationService,
	}
}

// Register handles user registration
func (ctrl *userController) Register(c *gin.Context) {
	var registerRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Code     string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 验证验证码
	if err := ctrl.verificationService.VerifyCode(registerRequest.Email, registerRequest.Code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &model.User{
		Username: registerRequest.Username,
		Password: registerRequest.Password,
		Email:    registerRequest.Email,
	}

	if err := ctrl.userService.Register(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}

// Login handles user login
func (ctrl *userController) Login(c *gin.Context) {
	var loginInfo struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password"`
		Code     string `json:"code"`
	}

	if err := c.ShouldBindJSON(&loginInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的登录信息"})
		return
	}

	// 验证至少提供了一种验证方式
	if loginInfo.Password == "" && loginInfo.Code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供密码或验证码"})
		return
	}

	token, err := ctrl.userService.Login(loginInfo.Email, loginInfo.Password, loginInfo.Code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "登录失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// SendVerificationCode sends verification code to user's email
func (ctrl *userController) SendVerificationCode(c *gin.Context) {
	var request struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的邮箱地址"})
		return
	}

	if err := ctrl.verificationService.SendVerificationCode(request.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "验证码已发送"})
}
