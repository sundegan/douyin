package controller

import (
	"douyin-server/dao"
	"douyin-server/service"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type VideoListResponse struct {
	Response
	VideoList []dao.Video `json:"video_list"`
}

// Publish 获取用户投稿的视频并保存到本地
func Publish(c *gin.Context) {
	id, ok := getId(c)
	if !ok {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取用户id失败，请重试"})
		return
	}

	// 获取用户上传的视频
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	// 视频文件的后缀，也即视频的格式
	fileSuffix := filepath.Ext(data.Filename)
	// 通过用户id和当前时间戳拼接成最终存放的视频文件名
	filename := fmt.Sprintf("%x_%x%s", id, time.Now().UnixNano(), fileSuffix)
	// 拼接存放视频的本地路径
	saveFile := filepath.Join("./static-server/videos/", filename)

	// 保存视频文件到本地
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	// 封面的文件名和最终存放的视频文件名一致，因为封面是图片，所以把后缀改为jpg
	covername := strings.TrimSuffix(filename, fileSuffix) + ".jpg"
	// 拼接存放封面的本地路径
	saveCover := filepath.Join("./static-server/covers/", covername)

	isGenerateOK := coverGenerator(saveFile, saveCover)

	if err = service.Publish(filename, covername, c.PostForm("title"), id, isGenerateOK); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  filename + " uploaded successfully",
	})

}

// PublishList 返回用户的投稿视频列表
func PublishList(c *gin.Context) {
	_id := c.Query("user_id")
	id, err := strconv.ParseInt(_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取用户id失败，请重试"})
		return
	}
	videoList := service.PublishList(id)
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videoList,
	})
}

// coverGenerator 会通过命令行调用ffmpeg提取视频的第一帧作为封面
func coverGenerator(videoDst, coverDst string) bool {
	cmd := exec.Command("ffmpeg", "-ss", "00:00:00", "-i", videoDst, "-vframes", "1", coverDst)
	err := cmd.Run()
	if err != nil {
		log.Println("生成帧出错：", err)
	}
	return err == nil
}
