package controller

import (
	"douyin-server/service"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func FavoriteAction(c *gin.Context) {
	_id, ok := c.Get("id")
	if !ok {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取用户id失败，请重试"})
		return
	}

	id, ok := _id.(int64)
	if !ok {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取用户id失败，请重试"})
		log.Println("出现无法解析成64位整数的token")
		return
	}

	video_id, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取视频id失败"})
		log.Println("出现无法解析成64位整数的视频id")
		return
	}
	action_type := c.Query("action_type")
	fmt.Println("user_id: " + strconv.Itoa(int(id)) + " video_id: " +
		strconv.Itoa(int(video_id)) +
		" action_type: " + action_type)

	err = service.Favorite(video_id, id, action_type)
	if err == errors.New("已点赞过该视频") {
		c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "已点赞过该视频"})
	} else if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	}
}

func FavoriteList(c *gin.Context) {
	_id, ok := c.Get("id")
	if !ok {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取用户id失败，请重试"})
		return
	}

	id, ok := _id.(int64)
	if !ok {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取用户id失败，请重试"})
		log.Println("出现无法解析成64位整数的token")
		return
	}

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
