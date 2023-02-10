package service

import (
	"sort"
	"tiktok-demo/dao/mysql"
)

func InsertWsMsg(userId, toUserId int, content string, flag int) error{
	// 对方账号是否存在
	_, err := mysql.CheckUserExist(toUserId)
	if err != nil {
		return err
	}
	err = mysql.InsertWsChatRecord(userId, toUserId, content,flag)
	return err
}

func UnreadAndHistory(userId, toUserId int) (msgRecords []mysql.WsChatRecord,err error){
	//得到自己的未读数据
	unreads,startTime,err:=mysql.GetWsUnreadRecords(userId, toUserId)
	if err!=nil{
		return 
	}
	//得到双方的历史数据
	historys,err:=mysql.GetWsHistorysByTime(userId, toUserId,startTime)
	if err!=nil{
		return 
	}
	//将未读消息和历史20条消息升序返回
	if len(historys)>20{
		historys=historys[len(historys)-20:]
	}
	msgRecords=append(unreads,historys...)
	sort.Slice(msgRecords,func(i, j int) bool {
		return msgRecords[i].CreateTime<msgRecords[j].CreateTime
	})
	return msgRecords,nil
}
