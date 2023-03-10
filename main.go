package main

import (
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"tiktok-demo/conf"
	"tiktok-demo/controller"
	"tiktok-demo/dao/mysql"
	"tiktok-demo/dao/redis"
	"tiktok-demo/logger"
	"tiktok-demo/middleware/kafka"
	"tiktok-demo/middleware/snowflake"
	"tiktok-demo/router"
	"tiktok-demo/service"
)

// 初始化项目所有依赖
func initDependencies() error {
	err := conf.InitConfig()
	if err != nil {
		return err
	}
	err = logger.InitLogger()
	if err != nil {
		return err
	}
	err = mysql.InitMysql()
	if err != nil {
		return err
	}
	err = redis.InitRedis()
	if err != nil {
		return err
	}
	err = controller.InitTrans("en")
	if err != nil {
		return err
	}
	err = service.InitMinio()
	if err != nil {
		return err
	}
	err = snowflake.InitSonyFlake(uint16(conf.Config.MachineID))
	if err != nil {
		return err
	}
	err = kafka.InitProducers()
	if err != nil {
		return err
	}
	err = kafka.InitConsumerGroups()
	return err
}

func main() {
	defer func() {
		err := mysql.Close()
		if err != nil {
			return
		}
	}()
	err := initDependencies()
	if err != nil {
		fmt.Printf("initDependencies failed, err:%v\n", err)
	}
	r := gin.New()
	// 替换gin框架日志，自定义GinRecovery
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	//websocket监听
	go controller.Manager.Start()
	// 路由设置
	router.InitRouters(r)
	pprof.Register(r)
	// 自定义修改端口
	err = r.Run(":8000")
	if err != nil {
		return
	}
}
