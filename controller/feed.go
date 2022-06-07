package controller

import (
	"douyin-server/dao"
	"douyin-server/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

type FeedResponse struct {
	Response
	VideoList []dao.Video `json:"video_list,omitempty"`
	NextTime  int64       `json:"next_time,omitempty"`
}

func Feed(c *gin.Context) {
	token := c.Query("token")

	var id int64
	// 用户携带了token
	if token != "" {
		_id, err := dao.RDB.Get(dao.Ctx, token).Result()
		if err == nil {
			// token续期
			dao.RDB.Expire(dao.Ctx, token, 12*time.Hour)
		} else {
			// token不存在，记录该ip此次访问
			ipAddress := c.ClientIP()
			times, _ := dao.RDB.Get(dao.Ctx, ipAddress).Int64()
			if times > 10 {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "疑似恶意访问，请勿携带非法token！"})
				c.Abort()
				return
			}
			dao.RDB.Set(dao.Ctx, ipAddress, times+1, time.Minute)
		}

		id, err = strconv.ParseInt(_id, 10, 64)
		if err != nil {
			log.Println("出现无法解析成64位整数的token")
		}
	}

	_latestTime := c.Query("latest_time")
	latestTime, err := strconv.ParseInt(_latestTime, 10, 64)
	// 说明时间戳格式有错
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1})
		return
	}

	videoList, nextTime := service.Feed(id, latestTime)
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
