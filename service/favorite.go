package service

import (
	"douyin-server/dao"
	"errors"
	"gorm.io/gorm"
)

// Favorite 点赞/取消点赞视频
func Favorite(videoId int64, userId int64, actionType string) error {
	video := dao.Video{}

	//判断video是否存在
	err := dao.DB.Model(&dao.Video{}).Where("id = ?", videoId).
		Select("author_id", "favorite_count").Find(&video).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("点赞的视频不存在")
	}

	//获取当前用户点赞视频总数以及视频作者被点赞总数
	userFavoriteCount, authorTotalFavorited := 0, 0
	err = dao.DB.Model(&dao.User{}).Where("id = ?", userId).
		Select("favorite_count").Find(&userFavoriteCount).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("执行点赞操作的用户不存在")
	}
	err = dao.DB.Model(&dao.User{}).Where("id = ?", video.AuthorId).
		Select("total_favorited").Find(&authorTotalFavorited).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("视频的作者不存在")
	}

	//获取视频作者信息
	author := dao.User{}
	err = dao.DB.Where("id = ?", video.AuthorId).Find(&author).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("视频的作者不存在")
	}

	if actionType == "1" { //点赞视频
		//通过从点赞关联表中查询结果，判断当前用户是否已点赞过该视频
		rows, err := dao.DB.Table("favorites").Where("user_id = ? AND video_id = ?", userId, videoId).Rows()
		if err != nil {
			return err
		}
		if rows.Next() {
			return errors.New("已点赞过该视频")
		}

		//更新数值
		userFavoriteCount++
		video.FavoriteCount++
		authorTotalFavorited++

		//在点赞关联表新增关联
		favorite := dao.Favorite{
			UserId:  userId,
			VideoId: videoId,
		}
		if err := dao.DB.Create(&favorite).Error; err != nil {
			return err
		}
	} else if actionType == "2" { //取消点赞
		//如果当前视频点赞数不为正数，则取消点赞操作是异常的，会导致视频点赞数变为负数
		if userFavoriteCount <= 0 {
			return errors.New("取消点赞异常，用户点赞数为非正数")
		} else if video.FavoriteCount <= 0 {
			return errors.New("取消点赞异常，视频点赞数为非正数")
		}

		//更新数值
		userFavoriteCount--
		video.FavoriteCount--
		authorTotalFavorited--

		//从点赞关联表中删除关联
		favorite := dao.Favorite{}
		if err := dao.DB.Where("user_id = ? AND video_id = ?", userId, videoId).
			Delete(&favorite).Error; err != nil {
			return err
		}
	}

	//更新表中与点赞相关的数据
	err = dao.DB.Model(&dao.Video{}).Where("id = ?", videoId).Update("favorite_count", video.FavoriteCount).Error
	if err != nil {
		return errors.New("视频点赞数修改失败")
	}
	err = dao.DB.Model(&dao.User{}).Where("id = ?", userId).Update("favorite_count", userFavoriteCount).Error
	if err != nil {
		return errors.New("用户点赞数修改失败")
	}
	err = dao.DB.Model(&dao.User{}).Where("id = ?", video.AuthorId).Update("total_favorited", authorTotalFavorited).Error
	if err != nil {
		return errors.New("作者被点赞数修改失败")
	}
	err = dao.DB.Save(&author).Error
	if err != nil {
		return errors.New("作者被点赞数修改失败")
	}

	return nil
}

// Favorite 获取点赞列表
func FavoriteList(userId int64) ([]dao.Video, error) {
	//依据用户id查询当前用户点赞的所有视频的视频id，存入favorite数组
	favorite := []dao.Favorite{}
	err := dao.DB.Where("user_id = ?", userId).Find(&favorite).Error
	if err != nil {
		return nil, err
	}
	//依据favorite数组中的数据查询视频列表
	videoList := []dao.Video{}
	for _, f := range favorite {
		video := dao.Video{}
		err := dao.DB.Where("id = ?", f.VideoId).Find(&video).Error
		if err != nil {
			return nil, err
		}
		videoList = append(videoList, video)
	}
	return videoList, nil
}
