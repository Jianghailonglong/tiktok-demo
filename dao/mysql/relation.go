package mysql

import "tiktok-demo/logger"

type Relation struct {
	Id         int `gorm:"column:id"`
	UserId     int `gorm:"column:user_id"`
	ToUserId   int `gorm:"column:to_user_id"`
	Subscribed int `gorm:"column:subscribed"`
}

func (Relation) TableName() string {
	return "relations"
}

const (
	SUBSCRIBED   = 1
	UNSUBSCRIBED = 0
)

// 根据关注着和被关注着找到关联关系
func GetRelation(userId int, toUserId int) (*Relation, error) {
	relation := Relation{}
	if err := db.
		Where("user_id = ?", userId).
		Where("to_user_id = ?", toUserId).
		Take(&relation).Error; nil != err {
		logger.Log.Info(err.Error())
		return nil, err
	}

	return &relation, nil
}

// 在原来的关系上修改关注关系
func UpdateRelation(relation *Relation, action int) (bool, error) {
	// 更新失败，返回错误。
	if err := db.Model(Relation{}).
		Where("id = ?", relation.Id).
		Update("subscribed", action).Error; nil != err {
		// 更新失败，打印错误日志。
		logger.Log.Info(err.Error())
		return false, err
	}
	// 更新成功。
	return true, nil
}

// 原表中没有关注记录，新增一条记录
func AddRelation(userId int, toUserId int) (bool, error) {
	relation := Relation{
		UserId:     userId,
		ToUserId:   toUserId,
		Subscribed: SUBSCRIBED,
	}

	// 插入失败，返回err.
	if err := db.Select("UserId", "ToUserId", "Subscribed").Create(&relation).Error; nil != err {
		logger.Log.Info(err.Error())
		return false, err
	}
	// 插入成功
	return true, nil
}
