package service

import (
	"douyin-server/dao"
	"errors"
	"gorm.io/gorm"
	"time"
)

func Comment(video_id int64, user_id int64, comment_text string) (dao.Comment, error) {
	comment := dao.Comment{}

	user := dao.User{}
	err := dao.DB.Where("id = ?", user_id).Find(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return comment, errors.New("执行评论操作的用户不存在")
	}
	//去除评论用户的敏感信息
	user.Pwd = ""
	user.Salt = ""

	video := dao.Video{}
	err = dao.DB.Where("id = ?", video_id).Find(&video).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return comment, errors.New("评论的视频不存在")
	}
	video.CommentCount++
	if err := dao.DB.Save(&video).Error; err != nil {
		return comment, err
	}

	comment.UserId = user.Id
	comment.VideoId = video_id
	comment.Content = comment_text
	comment.CreateDate = time.Now().Format("01-02")
	comment.User = user

	if err := dao.DB.Create(&comment).Error; err != nil {
		return comment, err
	}
	return comment, nil
}

func DeleteComment(comment_id int64) error {
	comment := dao.Comment{}
	err := dao.DB.Where("id = ?", comment_id).Find(&comment).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("删除的评论不存在")
	}
	video := dao.Video{}
	err = dao.DB.Where("id = ?", comment.VideoId).Find(&video).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("当前评论所在的视频不存在")
	}
	if video.CommentCount <= 0 {
		return errors.New("当前评论所在的视频评论数异常")
	}
	video.CommentCount--
	if err = dao.DB.Save(&video).Error; err != nil {
		return err
	}
	if err = dao.DB.Delete(&comment).Error; err != nil {
		return err
	}
	return nil
}

func CommentList(video_id int64) ([]dao.Comment, error) {

	//判断video是否存在
	video := dao.Video{}
	err := dao.DB.Where("id = ?", video_id).Find(&video).Error
	if err != nil {
		return nil, err
	}

	//用video_id获取评论列表
	var commentList []dao.Comment //comment表
	err = dao.DB.Where("video_id = ?", video_id).Order("create_date DESC").Find(&commentList).Error

	var returnList []dao.Comment
	for _, c := range commentList {
		var x = dao.Comment{}
		x = c
		user := dao.User{}
		err := dao.DB.Where("id = ?", c.UserId).Find(&user).Error
		if err != nil {
			return nil, err
		}
		user.Pwd = ""
		user.Salt = ""
		x.User = user
		returnList = append(returnList, x)
	}
	if err != nil {
		return nil, err
	}

	return returnList, nil
}
