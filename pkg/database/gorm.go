package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"github.com/mhusainh/FastTix/config"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDatabase(mysqlconfig config.MySQLConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlconfig.User,
		mysqlconfig.Password,
		mysqlconfig.Host,
		mysqlconfig.Port,
		mysqlconfig.Database,
	)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
}
