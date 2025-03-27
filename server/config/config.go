package config

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"jing_vue_gin_admin/server/internal/model"
)

var DB *gorm.DB

func InitDB() error {
	var err error
	DB, err = gorm.Open(sqlite.Open("admin.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	// 自动迁移数据库表
	err = DB.AutoMigrate(&model.User{})
	if err != nil {
		return err
	}

	return nil
}
