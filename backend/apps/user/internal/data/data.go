package data

import (
	"context"
	"errors"

	"StructForge/backend/common/data/database"
	"StructForge/backend/common/log"

	"gorm.io/gorm"
)

// Data 数据访问层
type Data struct {
	db database.Database
}

// NewData 创建数据访问层实例
func NewData(db database.Database) (*Data, func(), error) {
	ctx := context.Background()

	if db == nil {
		return nil, nil, errors.New("database instance is nil")
	}

	// 执行数据库迁移
	if err := db.AutoMigrate(ctx, &User{}, &UserProfile{}, &EmailVerification{}); err != nil {
		log.Error(ctx, "数据库迁移失败",
			log.ErrorField(err),
		)
		return nil, nil, err
	}

	log.Info(ctx, "用户服务数据访问层初始化成功")

	cleanup := func() {
		log.Info(ctx, "用户服务数据访问层关闭")
	}

	return &Data{
		db: db,
	}, cleanup, nil
}

// DB 获取 GORM 数据库实例
func (d *Data) DB() *gorm.DB {
	if d.db == nil {
		return nil
	}
	return d.db.GetDB()
}
