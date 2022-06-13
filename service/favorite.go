package service

import (
	"douyin-server/dao"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"strconv"
)

// Favorite 点赞/取消点赞视频
func Favorite(videoId int64, userId int64, actionType string) (err error) {
	// 先通过布隆过滤器过滤无效的用户id
	if !userIdFilter.TestString(strconv.FormatInt(userId, 10)) {
		return errors.New("当前操作用户不存在")
	}

	tryTimes := 0

	// 先获取视频的作者id，并判断视频是否存在
	var authorId int64
	err = dao.DB.Model(&dao.Video{}).Where("id = ?", videoId).Select("author_id").Find(&authorId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("点赞的视频不存在")
	}

Transaction:
	if tryTimes == 3 {
		return errors.New("服务器繁忙，请稍后重试")
	}
	tryTimes++

	// 开启事务，后面的查询直接查数据库并加锁，以免出现数据不一致现象
	tx := dao.DB.Begin()

	// 获取当前用户点赞视频总数以及视频作者被点赞总数，先获取在User表中id较小的数据，避免循环等待导致死锁
	var userFavoriteCount, authorTotalFavorited int64
	if userId == authorId {
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(&dao.User{}).Where("id = ?", userId).
			Select("favorite_count").Find(&userFavoriteCount).Error
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Model(&dao.User{}).Where("id = ?", userId).
			Select("total_favorited").Find(&authorTotalFavorited).Error
		if err != nil {
			tx.Rollback()
			return
		}
	} else if userId < authorId {
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(&dao.User{}).Where("id = ?", userId).
			Select("favorite_count").Find(&userFavoriteCount).Error
		if err != nil {
			tx.Rollback()
			return
		}

		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(&dao.User{}).Where("id = ?", authorId).
			Select("total_favorited").Find(&authorTotalFavorited).Error
		if err != nil {
			tx.Rollback()
			log.Println("出现非法作者的视频，应检查")
			return
		}
	} else {
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(&dao.User{}).Where("id = ?", authorId).
			Select("total_favorited").Find(&authorTotalFavorited).Error
		if err != nil {
			tx.Rollback()
			log.Println("出现非法作者的视频，应检查")
			return
		}

		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(&dao.User{}).Where("id = ?", userId).
			Select("favorite_count").Find(&userFavoriteCount).Error
		if err != nil {
			tx.Rollback()
			return
		}
	}

	// 判断video是否存在，不等待的方式获取锁，不成功再重试
	var videoFavoriteCount int64
	err = tx.Clauses(clause.Locking{Strength: "UPDATE", Options: "NOWAIT"}).Model(&dao.Video{}).
		Where("id = ?", videoId).Select("favorite_count").Find(&videoFavoriteCount).Error
	if err != nil {
		tx.Rollback()
		goto Transaction
	}

	if actionType == "1" { //点赞视频
		// 通过从点赞关联表中查询结果，判断当前用户是否已点赞过该视频，该错误可忽略
		rows, _ := dao.DB.Model(&dao.Favorite{}).Where("user_id = ? AND video_id = ?", userId, videoId).Rows()
		if rows.Next() {
			tx.Rollback()
			return errors.New("已点赞过该视频")
		}

		// 更新数值
		userFavoriteCount++
		videoFavoriteCount++
		authorTotalFavorited++

		// 在点赞关联表新增关联
		favorite := dao.Favorite{
			UserId:  userId,
			VideoId: videoId,
		}
		if err = dao.DB.Model(&dao.Favorite{}).Create(&favorite).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else if actionType == "2" { //取消点赞
		// 如果当前视频点赞数不为正数，则取消点赞操作是异常的，会导致视频点赞数变为负数
		if userFavoriteCount <= 0 {
			tx.Rollback()
			return errors.New("取消点赞异常，用户点赞数为非正数")
		} else if videoFavoriteCount <= 0 {
			tx.Rollback()
			return errors.New("取消点赞异常，视频点赞数为非正数")
		}

		//更新数值
		userFavoriteCount--
		videoFavoriteCount--
		authorTotalFavorited--

		// 从点赞关联表中删除关联
		favorite := dao.Favorite{}
		if err = dao.DB.Model(&dao.Favorite{}).Where("user_id = ? AND video_id = ?", userId, videoId).
			Delete(&favorite).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	//更新表中与点赞相关的数据
	err = tx.Model(&dao.Video{}).Where("id = ?", videoId).Update("favorite_count", videoFavoriteCount).Error
	if err != nil {
		tx.Rollback()
		return errors.New("视频点赞数修改失败")
	}
	err = tx.Model(&dao.User{Id: userId}).Update("favorite_count", userFavoriteCount).Error
	if err != nil {
		tx.Rollback()
		return errors.New("用户点赞数修改失败")
	}
	err = tx.Model(&dao.User{Id: authorId}).Update("total_favorited", authorTotalFavorited).Error
	if err != nil {
		tx.Rollback()
		return errors.New("作者被点赞数修改失败")
	}

	tx.Commit()
	return nil
}

// FavoriteList 获取点赞列表
func FavoriteList(userId int64) ([]dao.Video, error) {
	//依据用户id查询当前用户点赞的所有视频的视频id，存入favorite数组
	favorite := []dao.Favorite{}
	err := dao.DB.Model(&dao.Favorite{}).Where("user_id = ?", userId).Find(&favorite).Error
	if err != nil {
		return nil, err
	}
	//依据favorite数组中的数据查询视频列表
	videoList := []dao.Video{}
	for _, f := range favorite {
		video := dao.Video{}
		err = dao.DB.Model(&dao.Video{}).Where("id = ?", f.VideoId).Find(&video).Error
		if err != nil {
			return nil, err
		}
		videoList = append(videoList, video)
	}
	return videoList, nil
}
