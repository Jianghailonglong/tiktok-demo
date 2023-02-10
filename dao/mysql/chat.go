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
	Flag       int    `json:"flag" gorm:"type:int(1)"`
}

func (ChatRecord) TableName() string {
	return "chat_records"
}

// 插入新的聊天记录
func InsertChatRecord(userId, toUserId, actionType int, content string) (err error) {
	chatRecord := ChatRecord{
		SourceId:   userId,
		Content:    content,
		CreateTime: int(time.Now().Unix()),
		TargetId:   toUserId,
		Flag:       0,
	}
	res := db.Create(&chatRecord)
	if res.Error != nil {
		return errors.New("InsertChatRecord插入失败")
	}
	return
}

// 最好是:  第一次---未读和10条历史
//         第二次开始---未读
func GetChatRecordList(userId, toUserId int) (messageList []common.Message, err error) {
	recordListA2B, recordListB2A := []ChatRecord{}, []ChatRecord{}
	res := db.Where("source_id = ?", userId).Where("target_id = ?", toUserId).Find(&recordListA2B)
	if res.Error != nil {
		return messageList, errors.New("GetChatRecordList查询失败")
	}
	res = db.Where("source_id = ?", toUserId).Where("target_id = ?", userId).Find(&recordListB2A)
	if res.Error != nil {
		return messageList, errors.New("GetChatRecordList查询失败")
	}
	messageList = AppendRecordList(recordListA2B, recordListB2A)
	return messageList, nil
}

//得到未读列表
func GetChatUnreadList(userId, toUserId int) (messageList []common.Message, err error) {
	recordListB2A := []ChatRecord{}
	UnreadDB := db.Where("source_id = ?", toUserId).Where("target_id = ?", userId).Where("flag=?", 0)
	res := UnreadDB.Find(&recordListB2A)
	if res.Error != nil {
		return messageList, errors.New("GetChatRecordList查询失败")
	}
	_ = UnreadDB.Select("flag").Updates(map[string]interface{}{"flag": 1})
	messageList = ChatRecords2Messages(recordListB2A)
	return messageList, nil
}

//按createTime升序排序两方的聊天记录
func AppendRecordList(listA2B, listB2A []ChatRecord) (messageList []common.Message) {
	listA2B = append(listA2B, listB2A...)
	sort.Slice(listA2B, func(i, j int) bool {
		return listA2B[i].CreateTime < listA2B[j].CreateTime
	})
	messageList = ChatRecords2Messages(listA2B)
	return messageList
}

//将[]mysql.ChatRecord转化成[]common.Message
func ChatRecords2Messages(chatRecords []ChatRecord) (messages []common.Message) {
	for idx, l := range chatRecords {
		msg := common.Message{
			Id:         int64(idx + 1),
			ToUserId:   int64(l.TargetId),
			FromUserId: int64(l.SourceId),
			Content:    l.Content,
			CreateTime: int64(l.CreateTime),
		}
		messages = append(messages, msg)
	}
	return messages
}
