package service

import (
	"douyin-server/dao"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

// Comment 新增评论
func Comment(videoId int64, userId int64, commentText string) (comment dao.Comment, err error) {
	// 判断进行评论的用户是否存在
	user, err := UserInfo(userId)
	if err != nil {
		return comment, errors.New("执行评论操作的用户不存在")
	}

	// 开启事务，以便获取行锁
	tx := dao.DB.Begin()
	// 判断video是否存在
	var videoCommentCount int64
	err = tx.Model(&dao.Video{}).Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", videoId).Select("comment_count").Find(&videoCommentCount).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return comment, errors.New("评论的视频不存在")
	}
	// 增加视频评论数并更新
	videoCommentCount++
	if err = tx.Model(&dao.Video{}).Where("id = ?", videoId).
		Update("comment_count", videoCommentCount).Error; err != nil {
		tx.Rollback()
		return comment, err
	}
	tx.Commit()

	// 新建comment对象并保存在数据库中
	comment.UserId = user.Id
	comment.VideoId = videoId
	comment.Content = commentText
	comment.CreateDate = time.Now().Format("01-02")
	comment.User = user
	if err = dao.DB.Model(&dao.Comment{}).Create(&comment).Error; err != nil {
		return comment, err
	}
	tx.Commit()
	return comment, nil
}

// DeleteComment 删除评论
func DeleteComment(commentId int64) (err error) {
	comment := dao.Comment{}

	// 判断删除的评论是否存在
	err = dao.DB.Model(&dao.Comment{}).Where("id = ?", commentId).Find(&comment).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("删除的评论不存在")
	}

	// 判断删除评论所在的视频是否存在
	// 开启事务，以便获取行锁
	tx := dao.DB.Begin()
	// 判断video是否存在
	var videoCommentCount int64
	err = tx.Model(&dao.Video{}).Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", comment.VideoId).Select("comment_count").Find(&videoCommentCount).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return errors.New("当前评论所在的视频不存在")
	}
	// 如果当前视频评论数不为正数，则删除评论操作是异常的，会导致视频评论数变为负数
	if videoCommentCount <= 0 {
		tx.Rollback()
		return errors.New("当前评论所在的视频评论数异常")
	}
	// 减少视频评论数并更新
	videoCommentCount--
	if err = tx.Model(&dao.Video{}).Where("id = ?", comment.VideoId).
		Update("comment_count", videoCommentCount).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	// 删除在评论表中的对应评论行数据
	if err = dao.DB.Model(&dao.Comment{}).Delete(&comment).Error; err != nil {
		return err
	}
	return nil
}

// CommentList 获取评论列表
func CommentList(videoId int64) (returnList []dao.Comment, err error) {
	//判断video是否存在
	video := dao.Video{}
	err = dao.DB.Model(&dao.Video{}).Where("id = ?", videoId).Find(&video).Error
	if err != nil {
		return nil, err
	}

	// 用videoId获取评论列表，按发布日期进行降序排序
	var commentList []dao.Comment //comment表
	err = dao.DB.Model(&dao.Comment{}).Where("video_id = ?", videoId).
		Order("create_date DESC").Find(&commentList).Error
	if err != nil {
		return nil, err
	}

	// returnList存储需要返回的信息
	// 为每个comment添加user信息
	for _, c := range commentList {
		var x = dao.Comment{}
		x = c
		//获取该评论对应的用户信息
		user, err := UserInfo(c.UserId)
		if err != nil {
			return nil, err
		}
		x.User = user
		returnList = append(returnList, x)
	}

	//返回returnList
	return returnList, nil
}
