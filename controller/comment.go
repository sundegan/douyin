package controller

import (
	"douyin-server/dao"
	"douyin-server/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type CommentListResponse struct {
	Response
	CommentList []dao.Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment dao.Comment `json:"comment,omitempty"`
}

func CommentAction(c *gin.Context) {
	id, ok := getId(c)
	if !ok {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取用户id失败，请重试"})
		return
	}

	video_id, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取视频id失败"})
		log.Println("出现无法解析成64位整数的视频id")
		return
	}

	actionType := c.Query("action_type")

	if actionType == "1" {
		comment_text := c.Query("comment_text")
		comment, err := service.Comment(video_id, id, comment_text)
		if err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{StatusCode: 1, StatusMsg: "评论异常"},
			})
			return
		} else {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{StatusCode: 0},
				Comment:  comment,
			})
			return
		}
	} else if actionType == "2" {
		comment_id, err := strconv.ParseInt(c.Query("comment_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取评论id失败"})
			log.Println("出现无法解析成64位整数的视频id")
			return
		}
		err = service.DeleteComment(comment_id)
		if err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{StatusCode: 1, StatusMsg: "删除评论异常"},
			})
			return
		} else {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{StatusCode: 0},
			})
		}
	}

}

func CommentList(c *gin.Context) {

	video_id, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "获取视频id失败"})
		log.Println("出现无法解析成64位整数的视频id")
		return
	}

	commentList, err := service.CommentList(video_id)
	if err != nil {
		c.JSON(http.StatusOK, CommentListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
			CommentList: nil,
		})
	}

	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: commentList,
	})
}
