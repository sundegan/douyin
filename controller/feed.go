package controller

import (
	"douyin-server/dao"
	"douyin-server/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type FeedResponse struct {
	Response
	VideoList []dao.Video `json:"video_list,omitempty"`
	NextTime  int64       `json:"next_time,omitempty"`
}

func Feed(c *gin.Context) {
	_latestTime := c.Query("latest_time")
	latestTime, err := strconv.ParseInt(_latestTime, 10, 64)
	// 说明时间戳格式有错
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1})
		return
	}

	videoList, nextTime := service.Feed(latestTime)
	// 说明有找到视频
	if nextTime != 0 {
		c.JSON(http.StatusOK, FeedResponse{
			Response:  Response{StatusCode: 0},
			VideoList: videoList,
			NextTime:  nextTime,
		})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	}
}
