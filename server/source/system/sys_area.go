package system

import (
	"context"
	sysModel "server/model/system"
	"server/service/system"

	"gorm.io/gorm"
)

const initOrderArea = system.InitOrderSystem + 1

type initArea struct{}

// auto run
func init() {
	system.RegisterInit(initOrderArea, &initArea{})
}

func (i *initArea) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(&sysModel.SysArea{})
}

func (i *initArea) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&sysModel.SysArea{})
}

func (i *initArea) InitializerName() string {
	return sysModel.SysArea{}.TableName()
}

func (i *initArea) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	// 检查是否已有数据
	var count int64
	db.Model(&sysModel.SysArea{}).Count(&count)
	if count > 0 {
		return ctx, nil
	}

	// 插入一些初始区域数据（示例）
	areas := []sysModel.SysArea{
		{
			N:     "北京",
			I:     11,
			P:     0,
			Y:     "b",
			Level: 1,
		},
		{
			N:     "北京",
			I:     1101,
			P:     11,
			Y:     "b",
			Level: 2,
		},
		{
			N:     "东城",
			I:     110101,
			P:     1101,
			Y:     "d",
			Level: 3,
		},
		{
			N:     "西城",
			I:     110102,
			P:     1101,
			Y:     "x",
			Level: 3,
		},
		{
			N:     "朝阳",
			I:     110105,
			P:     1101,
			Y:     "c",
			Level: 3,
		},
		{
			N:     "丰台",
			I:     110106,
			P:     1101,
			Y:     "f",
			Level: 3,
		},
	}

	if err := db.Create(&areas).Error; err != nil {
		return ctx, err
	}
	return ctx, nil
}

func (i *initArea) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	var count int64
	if err := db.Model(&sysModel.SysArea{}).Count(&count).Error; err != nil {
		return false
	}
	return count > 0
}
