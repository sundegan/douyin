package controller

import (
	"douyin-server/service"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

// FavoriteAction 点赞/取消点赞视频
func FavoriteAction(c *gin.Context) {
	//从上下文中获取执行当前操作的用户的id
	id, ok := getId(c)
	if !ok {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取用户id失败，请重试"})
		return
	}

	//获取执行操作的视频id
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取视频id失败"})
		log.Println("出现无法解析成64位整数的视频id")
		return
	}
	//获取当前操作类型
	actionType := c.Query("action_type")

	err = service.Favorite(videoId, id, actionType)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "用户或视频不存在！"})
	} else if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	}
}

// FavoriteList 点赞视频列表
func FavoriteList(c *gin.Context) {
	_id := c.Query("user_id")
	id, err := strconv.ParseInt(_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取用户id失败，请重试"})
		return
	}

	//获取点赞视频列表
	videoList, err := service.FavoriteList(id)
	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
			VideoList: nil,
		})
	}

	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videoList,
	})
}
