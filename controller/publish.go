package controller

import (
	"douyin-server/dao"
	"douyin-server/service"
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

type VideoListResponse struct {
	Response
	VideoList []dao.Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token_ := c.PostForm("token")
	token, err := strconv.ParseInt(token_, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	}
	id, ok := service.CheckToken(token)
	if !ok {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
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
	user := usersLoginInfo[token_]
	finalName := fmt.Sprintf("%d_%s", user.Id, filename)
	saveFile := filepath.Join("./public/videos/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	saveCover := filepath.Join("./public/covers", finalName[:len(finalName)-3]+"jpg")
	cmd := exec.Command("ffmpeg", "-i", saveFile, "-ss", "0", "-f", "image2", saveCover)
	err = cmd.Run()
	// if err != nil {
	// 	c.JSON(http.StatusOK, Response{
	// 		StatusCode: 1,
	// 		StatusMsg:  err.Error(),
	// 	})
	// 	return
	// }
	if err := service.Publish(finalName, c.PostForm("title"), id); err != nil {
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

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: DemoVideos,
	})
}
