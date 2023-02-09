package conf

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Config = new(AppConfig)

type AppConfig struct {
	Mode         string `mapstructure:"mode"`
	Port         int    `mapstructure:"port"`
	Name         string `mapstructure:"name"`
	Version      string `mapstructure:"version"`
	StartTime    string `mapstructure:"start_time"`
	MachineID    int    `mapstructure:"machine_id"`
	*LogConfig   `mapstructure:"log"`
	*MySQLConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
	*AuthConfig  `mapstructure:"auth"`
	*MinioConfig `mapstructure:"minio"`
}

type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	Port         int    `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host         string `mapstructure:"host"`
	Password     string `mapstructure:"password"`
	Port         int    `mapstructure:"port"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
}

type LogConfig struct {
	Level             string `mapstructure:"level"`
	Filename          string `mapstructure:"filename"`
	MaxSize           int    `mapstructure:"max_size"`
	MaxAge            int    `mapstructure:"max_age"`
	MaxBackups        int    `mapstructure:"max_backups"`
	GormSlowThreshold int    `mapstructure:"gorm_slowthreshold"`
}

type AuthConfig struct {
	JwtExpire int    `mapstructure:"jwt_expire"`
	JwtSecret string `mapstructure:"jwt_secret"`
}

type MinioConfig struct {
	EndPoint        string            `mapstructure:"endpoint"`
	AccessKeyID     string            `mapstructure:"access_key_id"`
	SecretAccessKey string            `mapstructure:"secret_access_key"`
	UseSsL          bool              `mapstructure:"use_ssL"`
	Video           MinioBucketConfig `mapstructure:"video"`
	Image           MinioBucketConfig `mapstructure:"image"`
}

type MinioBucketConfig struct {
	BucketName  string `mapstructure:"bucket_name"`
	ContentType string `mapstructure:"content_type"`
	URL         string `mapstructure:"url"`
}

func InitConfig() error {
	// 根据文件位置修改
	workDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd failed, err:%v", err)
	}
	viper.SetConfigFile(workDir + "/conf/config.yaml")
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		err := viper.Unmarshal(&Config)
		if err != nil {
			return
		}
	})
	err = viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("ReadInConfig failed, err:%v", err)
	}
	if err = viper.Unmarshal(&Config); err != nil {
		return fmt.Errorf("unmarshal to Conf failed, err:%v", err)
	}
	return err
}
