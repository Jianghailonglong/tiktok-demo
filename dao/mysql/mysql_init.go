package mysql

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"tiktok-demo/conf"
	"tiktok-demo/logger"
)

var db *gorm.DB

// InitMysql 初始化MySQL连接
func InitMysql() (err error) {
	// "user:password@tcp(host:port)/dbname"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Config.User, conf.Config.MySQLConfig.Password, conf.Config.MySQLConfig.Host,
		conf.Config.MySQLConfig.Port, conf.Config.MySQLConfig.DBName)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.NewGormLogger(conf.Config.GormSlowThreshold),
	})
	if err != nil {
		return
	}
	DB, err := db.DB()
	if err != nil {
		logger.Log.Error("db.DB failed", zap.Any("error", err.Error()))
		return err
	}
	DB.SetMaxIdleConns(conf.Config.MySQLConfig.MaxIdleConns)
	DB.SetMaxOpenConns(conf.Config.MySQLConfig.MaxOpenConns)
	logger.Log.Info("init mysql success")
	return
}

// Close 关闭MySQL连接
func Close() error {
	DB, err := db.DB()
	if err != nil {
		logger.Log.Error("db.DB failed", zap.Any("error", err.Error()))
		return err
	}
	err = DB.Close()
	if err != nil {
		logger.Log.Error("DB.Close failed", zap.Any("error", err.Error()))
		return err
	}
	return nil
}
