package main

import (
	"fmt"

	"github.com/davlin-coder/davlin/internal/config"
	"github.com/davlin-coder/davlin/internal/model"
	"github.com/davlin-coder/davlin/internal/resource/container"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	c, err := container.NewContainer()
	if err != nil {
		panic(err)
	}

	// 执行数据库迁移
	err = c.Invoke(func(db *gorm.DB) error {
		return db.AutoMigrate(
			&model.User{},
		)
	})
	if err != nil {
		panic(err)
	}

	err = c.Invoke(func(conf *config.Config, router *gin.Engine) error {
		port := fmt.Sprintf(":%d", conf.APP.Port)
		fmt.Println("Service is starting at " + port + "...")
		return router.Run(port)
	})
	if err != nil {
		panic(err)
	}
}
