package snowflake

import (
	"fmt"
	"github.com/sony/sonyflake"
	"tiktok-demo/logger"
	"time"
)

var sonyFlake *sonyflake.Sonyflake

// InitSonyFlake 初始化配置
func InitSonyFlake(machineId uint16) (err error) {
	t, err := time.Parse("2006-01-02", "2023-02-01") // 初始化一个开始的时间
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}
	settings := sonyflake.Settings{ // 生成全局配置
		StartTime: t,
		MachineID: func() (uint16, error) {
			return machineId, nil // 指定机器ID
		},
	}
	sonyFlake = sonyflake.NewSonyflake(settings) // 用配置生成sonyFlake节点
	return
}

// GenID 返回生成的id值
func GenID() (id uint64, err error) { // 拿到sonyFlake节点生成id值
	if sonyFlake == nil {
		err = fmt.Errorf("snoy flake not inited")
		return
	}
	id, err = sonyFlake.NextID()
	return
}
