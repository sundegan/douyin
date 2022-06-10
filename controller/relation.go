package controller

import (
	"douyin-server/service"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserListResponse struct {
	Response
	UserList []service.ActionUser `json:"user_list"`
}

// RelationAction 登录用户对其他用户进行关注或取消关注操作
func RelationAction(c *gin.Context) {
	// 获取用户id
	userId, ok := getId(c)
	if !ok {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取用户id失败，请重试"})
		return
	}

	// 获取关注用户id
	toUserId, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取对方用户id失败"})
		return
	}

	// 获取操作类型,1表示关注,2表示取消关注
	actionType := c.Query("action_type")

	// 调用service层Action函数完成关注和取消操作
	err = service.Action(userId, toUserId, actionType)
	if err == errors.New("已关注该用户") {
		c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "已关注该用户"})
	} else if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	}
}

// FollowList 返回用户的关注列表
func FollowList(c *gin.Context) {
	// 获取用户id
	_id := c.Query("user_id")
	id, err := strconv.ParseInt(_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取用户id失败，请重试"})
		return
	}
	
	// 调用service层FollowList函数获得该用户的关注列表
	userList, err := service.FollowList(id)
	if err != nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
			UserList: nil,
		})
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: userList,
	})
}

// FollowerList 返回用户的粉丝列表
func FollowerList(c *gin.Context) {
	// 获取用户id
	_id := c.Query("user_id")
	id, err := strconv.ParseInt(_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取用户id失败，请重试"})
		return
	}

	// 调用service层FollowerList函数获得该用户的粉丝列表
	userList, err := service.FollowerList(id)
	if err != nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
			UserList: nil,
		})
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: userList,
	})
}
