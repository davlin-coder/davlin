package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/davlin-coder/davlin/internal/model"
	"github.com/davlin-coder/davlin/internal/resource/email"
	"github.com/davlin-coder/davlin/internal/resource/redis"
	"github.com/davlin-coder/davlin/internal/resource/template"
	"github.com/davlin-coder/davlin/internal/resource/tools"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	Register(user *model.User) error
	Login(username, password, code string) (string, error)
}

type userService struct {
	db                  *gorm.DB
	jwt                 *tools.JWTManager
	verificationService VerificationService
}

// NewUserService creates a new user service instance
func NewUserService(db *gorm.DB, jwt *tools.JWTManager, verificationService VerificationService) UserService {
	return &userService{
		db:                  db,
		jwt:                 jwt,
		verificationService: verificationService,
	}
}

// Register handles user registration service
func (s *userService) Register(user *model.User) error {
	// Encrypt password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// Set creation time
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Save to database
	result := s.db.Create(user)
	return result.Error
}

// Login handles user login service
func (s *userService) Login(email, password, code string) (string, error) {
	var user model.User
	// Find user by email
	result := s.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return "", errors.New("用户不存在")
	}

	// 如果提供了验证码，先验证验证码
	if code != "" {
		if err := s.verificationService.VerifyCode(email, code); err != nil {
			return "", errors.New("密码错误")
		}
	} else {
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			return "", errors.New("密码错误")
		}
	}

	// Generate JWT token
	token, err := s.jwt.GenerateToken(user.ID, user.Username)
	if err != nil {
		return "", fmt.Errorf("生成token失败: %v", err)
	}
	return token, nil
}

// VerificationService 定义验证码服务接口
type VerificationService interface {
	SendVerificationCode(email string) error
	VerifyCode(email, code string) error
}

// verificationService 实现验证码服务接口
type verificationService struct {
	redisClient redis.RedisClient
	emailer     email.EmailSender
	tm          template.TemplateManager
}

// NewVerificationService 创建验证码服务实例
func NewVerificationService(redisClient redis.RedisClient, emailer email.EmailSender, tm template.TemplateManager) VerificationService {
	return &verificationService{
		redisClient: redisClient,
		emailer:     emailer,
		tm:          tm,
	}
}

// SendVerificationCode 生成并发送验证码
func (s *verificationService) SendVerificationCode(email string) error {
	// 生成6位随机验证码
	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	// 设置验证码有效期为15分钟
	err := s.redisClient.Set(context.Background(), fmt.Sprintf("verification_code:%s", email), code, 15*time.Minute)
	if err != nil {
		return fmt.Errorf("保存验证码失败: %v", err)
	}

	// 使用模板渲染邮件内容
	emailContent, err := s.tm.ExecuteTemplate("verification_email", map[string]interface{}{
		"Code":          code,
		"ExpireMinutes": 15,
	})
	if err != nil {
		return fmt.Errorf("渲染邮件模板失败: %v", err)
	}

	// 发送验证码邮件
	err = s.emailer.SendHTMLEmail([]string{email}, "注册验证码", emailContent)
	if err != nil {
		return fmt.Errorf("发送验证码失败: %v", err)
	}

	return nil
}

// VerifyCode 验证验证码
func (s *verificationService) VerifyCode(email, code string) error {
	savedCode, err := s.redisClient.Get(context.Background(), fmt.Sprintf("verification_code:%s", email))
	if err != nil {
		return fmt.Errorf("验证码无效或已过期")
	}

	if savedCode != code {
		return fmt.Errorf("验证码错误")
	}

	// 删除已使用的验证码
	_ = s.redisClient.Del(context.Background(), fmt.Sprintf("verification_code:%s", email))

	return nil
}
