package service

import (
	"douyin-server/dao"
	"errors"
	"gorm.io/gorm"
)

func Favorite(video_id int64, user_id int64, action_type string) error {
	video := dao.Video{}
	err := dao.DB.Model(&dao.Video{}).Where("id = ?", video_id).
		Select("author_id", "favorite_count").Find(&video).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("点赞的视频不存在")
	}

	userFavoriteCount, authorTotalFavorited := 0, 0
	err = dao.DB.Model(&dao.User{}).Where("id = ?", user_id).
		Select("favorite_count").Find(&userFavoriteCount).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("执行点赞操作的用户不存在")
	}
	err = dao.DB.Model(&dao.User{}).Where("id = ?", video.AuthorId).
		Select("total_favorited").Find(&authorTotalFavorited).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("视频的作者不存在")
	}

	if action_type == "1" {
		rows, err := dao.DB.Table("favorites").Where("user_id = ? AND video_id = ?", user_id, video_id).Rows()
		if err != nil {
			return err
		}
		if rows.Next() {
			return errors.New("已点赞过该视频")
		}

		userFavoriteCount++
		video.FavoriteCount++
		authorTotalFavorited++

		favorite := dao.Favorite{
			UserId:  user_id,
			VideoId: video_id,
		}

		if err := dao.DB.Create(&favorite).Error; err != nil {
			return err
		}
	} else if action_type == "2" {

		if userFavoriteCount <= 0 {
			return errors.New("取消点赞异常，用户点赞数为非正数")
		} else if video.FavoriteCount <= 0 {
			return errors.New("取消点赞异常，视频点赞数为非正数")
		}

		userFavoriteCount--
		video.FavoriteCount--
		authorTotalFavorited--

		favorite := dao.Favorite{}

		if err := dao.DB.Where("user_id = ? AND video_id = ?", user_id, video_id).
			Delete(&favorite).Error; err != nil {
			return err
		}
	}

	err = dao.DB.Model(&dao.Video{}).Where("id = ?", video_id).Update("favorite_count", video.FavoriteCount).Error
	if err != nil {
		return errors.New("视频点赞数修改失败")
	}

	err = dao.DB.Model(&dao.User{}).Where("id = ?", user_id).Update("favorite_count", userFavoriteCount).Error
	if err != nil {
		return errors.New("用户点赞数修改失败")
	}

	err = dao.DB.Model(&dao.User{}).Where("id = ?", video.AuthorId).Update("total_favorited", authorTotalFavorited).Error
	if err != nil {
		return errors.New("作者被点赞数修改失败")
	}

	return nil
}

func FavoriteList(user_id int64) ([]dao.Video, error) {
	favorite := []dao.Favorite{}
	err := dao.DB.Where("user_id = ?", user_id).Find(&favorite).Error
	if err != nil {
		return nil, err
	}
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
