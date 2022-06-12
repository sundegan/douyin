package service

import (
	"context"
	"douyin-server/dao"
	"errors"
	"strconv"
)

// ActionUser 关注列表和粉丝列表使用的数据结构
type ActionUser struct {
	IsFollow      bool   `json:"is_follow"`
	Id            int64  `json:"id"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	Name          string `json:"name"`
}

// Action 对用户进行关注取关操作
func Action(userId int64, toUserId int64, actionType string) error {
	// 判断当前用户是否存在
	userName, err := UserInfoByField(userId, "Name")
	if err != nil {
		return errors.New("执行关注或取关操作的用户不存在")
	}

	// 判断关注用户是否存在
	toUserName, err := UserInfoByField(toUserId, "Name")
	if err != nil {
		return errors.New("关注的用户不存在")
	}

	// 关注操作
	if actionType == "1" {
		followed := dao.RdbFollow.HExists(context.Background(), strconv.FormatInt(userId, 10), strconv.FormatInt(toUserId, 10)).Val()
		if followed {
			return errors.New("已关注该用户")
		} else {
			// 关注列表增加数据
			dao.RdbFollow.HSet(
				context.Background(),
				strconv.FormatInt(userId, 10),
				strconv.FormatInt(toUserId, 10),
				toUserName,
			)
			// 粉丝列表增加数据
			dao.RdbFans.HSet(
				context.Background(),
				strconv.FormatInt(toUserId, 10),
				strconv.FormatInt(userId, 10),
				userName,
			)
		}
	} else if actionType == "2" { // 取关操作
		followed := dao.RdbFollow.HExists(context.Background(), strconv.FormatInt(userId, 10), strconv.FormatInt(toUserId, 10)).Val()
		if !followed {
			return errors.New("已经取关该用户")
		} else {
			// 关注列表删除数据
			dao.RdbFollow.HDel(
				context.Background(),
				strconv.FormatInt(userId, 10),
				strconv.FormatInt(toUserId, 10),
			)
			// 粉丝列表删除数据
			dao.RdbFans.HDel(
				context.Background(),
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
	userIdList := dao.RdbFollow.HKeys(context.Background(), strconv.FormatInt(userId, 10)).Val()
	// 遍历用户id,填充关注列表中用户的数据
	for _, id := range userIdList {
		var actionUser ActionUser
		actionUser.Id, _ = strconv.ParseInt(id, 10, 64)
		actionUser.Name = dao.RdbFollow.HGet(context.Background(), strconv.FormatInt(userId, 10), id).Val() // 用户的名称
		actionUser.FollowCount = dao.RdbFollow.HLen(context.Background(), id).Val()                         // 用户的关注总数
		actionUser.FollowerCount = dao.RdbFans.HLen(context.Background(), id).Val()                         // 用户的粉丝总数
		actionUser.IsFollow = true                                                                          // 在关注列表中的用户都为关注状态
		userList = append(userList, actionUser)
	}
	return userList, nil
}

// FollowerList 返回用户的粉丝列表
func FollowerList(userId int64) ([]ActionUser, error) {
	var userList []ActionUser
	// 从Redis的DB2获取该用户的所有粉丝id
	userIdList := dao.RdbFans.HKeys(context.Background(), strconv.FormatInt(userId, 10)).Val()
	// 遍历粉丝id,填充粉丝列表中用户的数据
	for _, id := range userIdList {
		var actionUser ActionUser
		actionUser.Id, _ = strconv.ParseInt(id, 10, 64)
		actionUser.Name = dao.RdbFans.HGet(context.Background(), strconv.FormatInt(userId, 10), id).Val() // 用户的名称
		actionUser.FollowCount = dao.RdbFollow.HLen(context.Background(), id).Val()                       // 用户的关注总数
		actionUser.FollowerCount = dao.RdbFans.HLen(context.Background(), id).Val()                       // 用户的粉丝总数

		// 判断是否互粉
		followed := dao.RdbFollow.HExists(context.Background(), strconv.FormatInt(userId, 10), id).Val()
		if followed { // 互粉,则is_follow字段设为true
			actionUser.IsFollow = true
		} else {
			actionUser.IsFollow = false
		}
		userList = append(userList, actionUser)
	}
	return userList, nil
}
