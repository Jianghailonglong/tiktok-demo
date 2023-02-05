package mysql

import (
	"errors"
	"sort"
	"tiktok-demo/common"
	"time"
)

type ChatRecord struct {
	Id         int	
	SourceId   int    `json:"source_id"`
	TargetId   int    `json:"target_id"` //消息接收方id
	Content    string `json:"content"`
	CreateTime int    `json:"create_time"`
}

func (ChatRecord) TableName() string {
	return "chat_records"
}

// 插入新的聊天记录
func InsertChatRecord(userId, toUserId, actionType int, content string) (err error) {
	chatRecord := ChatRecord{
		SourceId:         userId,
		Content:    content,
		CreateTime: int(time.Now().Unix()),
		TargetId:   toUserId,
	}
	res := db.Create(&chatRecord)
	if res.Error != nil {
		return errors.New("InsertChatRecord插入失败")
	}
	return
}

func GetChatRecordList(userId, toUserId int)(messageList []common.Message,err error){
	recordListA2B,recordListB2A:=[]ChatRecord{},[]ChatRecord{}
	res := db.Where("source_id = ?", userId).Where("target_id = ?", toUserId).Find(&recordListA2B)
	if res.Error != nil {
		return messageList, errors.New("GetChatRecordList查询失败")
	}
	res = db.Where("source_id = ?", toUserId).Where("target_id = ?", userId).Find(&recordListB2A)
	if res.Error != nil {
		return messageList, errors.New("GetChatRecordList查询失败")
	}
	messageList=AppendRecordList(recordListA2B,recordListB2A)
	return messageList,nil
}

// func AppendRecordList(listA2B,listB2A []ChatRecord)(messageList []common.Message){
// 	for _,l:=range listA2B{
// 		msg:=common.Message{
// 			Id: int64(l.SourceId),
// 			Content: l.Content,
// 			CreateTime: fmt.Sprint(l.CreateTime),
// 		}
// 		messageList=append(messageList, msg)
// 	}
// 	for _,l:=range listB2A{
// 		msg:=common.Message{
// 			Id: int64(l.SourceId),
// 			Content: l.Content,
// 			CreateTime: fmt.Sprint(l.CreateTime),
// 		}
// 		messageList=append(messageList, msg)
// 	}
// 	sort.Slice(messageList,func(i, j int) bool {
// 		ti,_:=strconv.Atoi(messageList[i].CreateTime)
// 		tj,_:=strconv.Atoi(messageList[j].CreateTime)
// 		return ti<tj
// 	})
// 	return messageList
// }
func AppendRecordList(listA2B,listB2A []ChatRecord)(messageList []common.Message){
	listA2B=append(listA2B, listB2A...)
	sort.Slice(listA2B,func(i, j int) bool {
		return listA2B[i].CreateTime<listA2B[j].CreateTime
	})
	for _,l:=range listA2B{
		msg:=common.Message{
			Id: int64(l.SourceId),
			Content: l.Content,
			CreateTime: time.Unix(int64(l.CreateTime),0).Format("2006-01-02 15:04:05"),
		}
		messageList=append(messageList, msg)
	}
	
	return messageList
}