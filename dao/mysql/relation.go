package mysql

import (
	"tiktok-demo/logger"
)

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
	SUBSCRIBED       = 1
	UNSUBSCRIBED     = 0
	RECORD_NOT_FOUND = "record not found"
)

// 根据关注着和被关注着找到关联关系
func GetRelation(userId int, toUserId int) (*Relation, error) {
	relation := Relation{}
	if err := db.
		Where("user_id = ?", userId).
		Where("to_user_id = ?", toUserId).
		Take(&relation).Error; nil != err {
		logger.Log.Error(err.Error())
		return nil, err
	}

	return &relation, nil
}

// 在原来的关系上修改关注关系
func UpdateRelation(relation *Relation, action int) error {
	// 更新失败，返回错误。
	if err := db.Model(Relation{}).
		Where("id = ?", relation.Id).
		Update("subscribed", action).Error; nil != err {
		// 更新失败，打印错误日志。
		logger.Log.Error(err.Error())
		return err
	}
	// 更新成功。
	return nil
}

// 原表中没有关注记录，新增一条记录
func AddRelation(userId int, toUserId int) error {
	relation := Relation{
		UserId:     userId,
		ToUserId:   toUserId,
		Subscribed: SUBSCRIBED,
	}

	// 插入失败，返回err.
	if err := db.Select("UserId", "ToUserId", "Subscribed").Create(&relation).Error; nil != err {
		logger.Log.Error(err.Error())
		return err
	}
	// 插入成功
	return nil
}

// 获取粉丝数量
func GetFollowerCnt(userId int64) (int64, error) {
	var cnt int64

	if err := db.Model(&Relation{}).
		Where("to_user_id = ?", userId).
		Where("subscribed = ?", SUBSCRIBED).
		Count(&cnt).Error; nil != err {

		return 0, err
	}

	return cnt, nil
}

// 获取关注数量
func GetFollowCnt(userId int64) (int64, error) {
	var cnt int64

	if err := db.Model(&Relation{}).
		Where("user_id = ?", userId).
		Where("subscribed = ?", SUBSCRIBED).
		Count(&cnt).Error; nil != err {

		return 0, err
	}

	return cnt, nil

}

// 获取当前用户的所有粉丝id
func GetFollowedIdList(userId int64) ([]int64, error) {
	var idList []int64
	if err := db.Model(&Relation{}).
		Where("to_user_id = ?", userId).
		Where("subscribed = ?", SUBSCRIBED).
		Pluck("user_id", &idList).Error; nil != err {

		return nil, err
	}

	return idList, nil
}

// 获取当前用户的所有关注id
func GetFollowIdList(userId int64) ([]int64, error) {
	var idList []int64
	if err := db.Model(&Relation{}).
		Where("user_id = ?", userId).
		Where("subscribed = ?", SUBSCRIBED).
		Pluck("to_user_id", &idList).Error; nil != err {

		return nil, err
	}

	return idList, nil
}

// 判断用户A是否是B的粉丝
func IsFollow(userAId int64, userBId int64) (bool, error) {

	relation := Relation{}

	if err := db.Where("user_id = ?", userAId).
		Where("to_user_id = ?", userBId).
		Where("subscribed = ?", SUBSCRIBED).
		Take(&relation).Error; nil != err {
		if RECORD_NOT_FOUND == err.Error() {
			// 查不到数据, 未关注
			return false, nil
		}
		// 其它错误
		return false, err
	}
	return true, nil
}
