package controller

import (
	"douyin-server/dao"
	"douyin-server/service"
	"errors"
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type VideoListResponse struct {
	Response
	VideoList []dao.Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	id, ok := getId(c)
	if !ok {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取用户id失败，请重试"})
		return
	}

	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	filename := filepath.Base(data.Filename)
	fileSuffix := filepath.Ext(data.Filename)
	finalName := fmt.Sprintf("%d_%s", id, filename)
	saveFile := filepath.Join("./static-server/videos/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	saveCover := filepath.Join("./static-server/covers/", strings.TrimSuffix(finalName, fileSuffix)+".jpg")

	err = coverGenerator(saveFile, saveCover)
	isGenerateOK := true
	if err != nil {
		fmt.Print(err)
		isGenerateOK = false
	}
	if err := service.Publish(finalName, c.PostForm("title"), id, isGenerateOK); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})

}

func PublishList(c *gin.Context) {
	id, ok := getId(c)
	if !ok {
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

func coverGenerator(videoDst, coverDst string) error {
	//ffmpeg command example: ffmpeg -ss 00:00:30 -i 666051400.mp4 -vframes 1 0.jpg
	cmd := exec.Command("ffmpeg", "-ss", "00:00:00", "-i", videoDst, "-vframes", "1", coverDst)
	err := cmd.Run()
	if err != nil {
		return errors.New("提取帧失败，使用默认封面")
	}
	return nil
}
