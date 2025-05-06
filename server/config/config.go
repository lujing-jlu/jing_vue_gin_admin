package config

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"jing_vue_gin_admin/server/internal/model"
)

var DB *gorm.DB

// InitDB 初始化数据库
func InitDB() error {
	var err error
	DB, err = gorm.Open(sqlite.Open("admin.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	// 自动迁移表结构
	err = DB.AutoMigrate(&model.User{}, &model.Role{}, &model.SystemLog{}, &model.File{})
	if err != nil {
		return err
	}

	return nil
}
