package service

import (
	"douyin-server/dao"
	"errors"
	"gorm.io/gorm"
	"strconv"
)

// ActionUser 关注列表和粉丝列表使用的数据结构
type ActionUser struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

// Action 对用户进行关注取关操作
func Action(userId int64, toUserId int64, actionType string) error {
	// 判断当前用户是否存在
	user := dao.User{}
	err := dao.DB.Where("id = ?", userId).Find(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("执行关注和取关操作的用户不存在")
	}

	// 判断关注用户是否存在
	toUser := dao.User{}
	err = dao.DB.Where("id = ?", toUserId).Find(&toUser).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("关注的用户不存在")
	}

	// 关注操作
	if actionType == "1" {
		isExist := dao.HExists(dao.RDB_FOLLOW, strconv.FormatInt(userId, 10), strconv.FormatInt(toUserId, 10))
		if isExist {
			return errors.New("已关注该用户")
		} else {
			// 关注列表增加数据
			dao.HSet(
				dao.RDB_FOLLOW,
				strconv.FormatInt(userId, 10),
				strconv.FormatInt(toUserId, 10),
				toUser.Name,
			)
			// 粉丝列表增加数据
			dao.HSet(
				dao.RDB_FANS,
				strconv.FormatInt(toUserId, 10),
				strconv.FormatInt(userId, 10),
				user.Name,
			)
		}
	} else if actionType == "2" { // 取关操作
		isExist := dao.HExists(dao.RDB_FOLLOW, strconv.FormatInt(userId, 10), strconv.FormatInt(toUserId, 10))
		if !isExist {
			return errors.New("已经取关该用户")
		} else {
			// 关注列表删除数据
			dao.HDel(
				dao.RDB_FOLLOW,
				strconv.FormatInt(userId, 10),
				strconv.FormatInt(toUserId, 10),
			)
			// 粉丝列表删除数据
			dao.HDel(
				dao.RDB_FANS,
				strconv.FormatInt(toUserId, 10),
				strconv.FormatInt(userId, 10),
			)
		}
	}
	return nil
}

// FollowList 返回用户的关注列表
func FollowList(userId int64) ([]ActionUser, error) {
	var userList []ActionUser
	// 从Redis的DB1获取该用户关注的所有用户id
	userIdList := dao.HKeys(dao.RDB_FOLLOW, strconv.FormatInt(userId, 10))
	// 遍历用户id,填充关注列表中用户的数据
	for _, id := range userIdList {
		var actionUser ActionUser
		actionUser.Id, _ = strconv.ParseInt(id, 10, 64)
		actionUser.Name = dao.HGet(dao.RDB_FOLLOW, strconv.FormatInt(userId, 10), id) // 用户的名称
		actionUser.FollowCount = dao.HLen(dao.RDB_FOLLOW, id)                         // 用户的关注总数
		actionUser.FollowerCount = dao.HLen(dao.RDB_FANS, id)                         // 用户的粉丝总数
		actionUser.IsFollow = true
		userList = append(userList, actionUser)
	}
	return userList, nil
}

// FollowerList 返回用户的粉丝列表
func FollowerList(userId int64) ([]ActionUser, error) {
	var userList []ActionUser
	// 从Redis的DB2获取该用户的所有粉丝id
	userIdList := dao.HKeys(dao.RDB_FANS, strconv.FormatInt(userId, 10))
	// 遍历粉丝id,填充粉丝列表中用户的数据
	for _, id := range userIdList {
		var actionUser ActionUser
		actionUser.Id, _ = strconv.ParseInt(id, 10, 64)
		actionUser.Name = dao.HGet(dao.RDB_FANS, strconv.FormatInt(userId, 10), id) // 用户的名称
		actionUser.FollowCount = dao.HLen(dao.RDB_FOLLOW, id)                       // 用户的关注总数
		actionUser.FollowerCount = dao.HLen(dao.RDB_FANS, id)                       // 用户的粉丝总数
		isExist := dao.HExists(dao.RDB_FOLLOW, strconv.FormatInt(userId, 10), id)   // 判断是否互粉
		if isExist {                                                                // 互粉,则is_follow字段设为true
			actionUser.IsFollow = true
		} else {
			actionUser.IsFollow = false
		}
		userList = append(userList, actionUser)
	}
	return userList, nil
}
