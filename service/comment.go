package service

import (
	"douyin-server/dao"
	"errors"
	"gorm.io/gorm"
	"time"
)

// Comment 新增评论
func Comment(videoId int64, userId int64, commentText string) (dao.Comment, error) {
	comment := dao.Comment{}

	//判断进行评论的用户是否存在
	user := dao.User{}
	err := dao.DB.Where("id = ?", userId).Find(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return comment, errors.New("执行评论操作的用户不存在")
	}

	//去除评论用户的敏感信息
	user.Pwd = ""
	user.Salt = ""

	//判断video是否存在
	video := dao.Video{}
	err = dao.DB.Where("id = ?", videoId).Find(&video).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return comment, errors.New("评论的视频不存在")
	}
	//增加视频评论数并更新
	video.CommentCount++
	if err := dao.DB.Save(&video).Error; err != nil {
		return comment, err
	}

	//更新新建的comment对象的信息并保存在数据库中
	comment.UserId = user.Id
	comment.VideoId = videoId
	comment.Content = commentText
	comment.CreateDate = time.Now().Format("01-02")
	comment.User = user
	if err := dao.DB.Create(&comment).Error; err != nil {
		return comment, err
	}
	return comment, nil
}

// DeleteComment 删除评论
func DeleteComment(commentId int64) error {
	comment := dao.Comment{}

	//判断删除的评论是否存在
	err := dao.DB.Where("id = ?", commentId).Find(&comment).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("删除的评论不存在")
	}

	//判断删除评论所在的视频是否存在
	video := dao.Video{}
	err = dao.DB.Where("id = ?", comment.VideoId).Find(&video).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("当前评论所在的视频不存在")
	}
	//如果当前视频评论数不为正数，则删除评论操作是异常的，会导致视频评论数变为负数
	if video.CommentCount <= 0 {
		return errors.New("当前评论所在的视频评论数异常")
	}
	//减少视频评论数并更新
	video.CommentCount--
	if err = dao.DB.Save(&video).Error; err != nil {
		return err
	}
	//删除在评论表中的对应评论行数据
	if err = dao.DB.Delete(&comment).Error; err != nil {
		return err
	}
	return nil
}

// CommentList 获取评论列表
func CommentList(videoId int64) ([]dao.Comment, error) {

	//判断video是否存在
	video := dao.Video{}
	err := dao.DB.Where("id = ?", videoId).Find(&video).Error
	if err != nil {
		return nil, err
	}

	//用video_id获取评论列表，按发布日期进行降序排序
	var commentList []dao.Comment //comment表
	err = dao.DB.Where("video_id = ?", videoId).Order("create_date DESC").Find(&commentList).Error

	//returnList存储需要返回的信息
	var returnList []dao.Comment
	//为每个comment添加user信息
	for _, c := range commentList {
		var x = dao.Comment{}
		x = c
		//获取该评论对应的用户信息
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

	//返回returnList
	return returnList, nil
}
