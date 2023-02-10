package mysql

import (
	"errors"
	"sort"
	"time"
)

type WsChatRecord struct {
	Id         int
	SourceId   int    `json:"source_id"`
	TargetId   int    `json:"target_id"` //消息接收方id
	Content    string `json:"content"`
	CreateTime int    `json:"create_time"`
	Flag       int    `json:"flag"` //是否已读 {0,1}
}

func (WsChatRecord) TableName() string {
	return "ws_chat_records"
}
func InsertWsChatRecord(userId, toUserId int, content string, flag int) error {
	wsChatRecord := WsChatRecord{
		SourceId:   userId,
		Content:    content,
		CreateTime: int(time.Now().Unix()),
		TargetId:   toUserId,
		Flag:       flag,
	}
	res := db.Create(&wsChatRecord)
	if res.Error != nil {
		return errors.New("InsertWsChatRecord插入失败")
	}
	return nil
}
func GetWsUnreadRecords(userId, toUserId int) (unreads []WsChatRecord, startTime int, err error) {
	unreadDB := db.Where("source_id=?", toUserId).Where("target_id=?", userId).Where("flag=?", 0).Order("create_time asc")
	//得到unreads升序结果
	res := unreadDB.Find(&unreads)
	if res.Error != nil {
		return unreads, 0, res.Error
	}

	if len(unreads) == 0 {
		startTime = int(time.Now().Unix())
	} else {
		startTime = unreads[0].CreateTime
	}

	//未读改为已读
	res = unreadDB.Select("flag").Updates(map[string]interface{}{"flag": 1})
	if res.Error != nil {
		return unreads, startTime, res.Error
	}
	return unreads, startTime, nil
}

func GetWsHistorysByTime(userId, toUserId, startTime int) (historys []WsChatRecord, err error) {
	//找到两方的历史记录
	myHist, yourHist := []WsChatRecord{}, []WsChatRecord{}
	res := db.Where("create_time<?", startTime).Where("source_id=?", userId).Where("target_id=?", toUserId).Find(&myHist)
	if res.Error != nil {
		return []WsChatRecord{}, res.Error
	}
	res = db.Where("create_time<?", startTime).Where("source_id=?", toUserId).Where("target_id=?", userId).Find(&yourHist)
	if res.Error != nil {
		return []WsChatRecord{}, res.Error
	}
	//合并两方的历史记录并排序
	historys = append(myHist, yourHist...)
	sort.Slice(historys, func(i, j int) bool {
		return historys[i].CreateTime < historys[j].CreateTime
	})
	return historys, nil
}
