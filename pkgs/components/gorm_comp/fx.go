package gorm_comp

import (
	"go.uber.org/fx"
	"gorm.io/gorm"
)

func ProvideGormOpt() *GormOpt {
	return LoadDatabaseConfigs()
}

func ProvideGormDB(gormDB *GormDB) *gorm.DB {
	return gormDB.GetDB()
}

var GormComponentFx = fx.Options(
	fx.Provide(ProvideGormOpt),
	fx.Provide(ProvideGormDB),
	fx.Provide(NewGormDB),
)
