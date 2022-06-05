package service

import (
	"douyin-server/dao"
	"errors"
	"gorm.io/gorm"
)

func Favorite(video_id int64, user_id int64, action_type string) error {
	user := dao.User{}
	err := dao.DB.Where("id = ?", user_id).Find(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("执行点赞操作的用户不存在")
	}
	video := dao.Video{}
	err = dao.DB.Where("id = ?", video_id).Find(&video).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("点赞的视频不存在")
	}
	author := dao.User{}
	err = dao.DB.Where("id = ?", video.AuthorId).Find(&author).Error
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

		user.FavoriteCount++
		video.FavoriteCount++
		author.TotalFavorited++

		favorite := dao.Favorite{
			UserId:  user_id,
			VideoId: video_id,
		}

		if err := dao.DB.Create(&favorite).Error; err != nil {
			return err
		}
	} else if action_type == "2" {

		if user.FavoriteCount <= 0 {
			return errors.New("取消点赞异常，用户点赞数为非正数")
		} else if video.FavoriteCount <= 0 {
			return errors.New("取消点赞异常，视频点赞数为非正数")
		}

		user.FavoriteCount--
		video.FavoriteCount--

		favorite := dao.Favorite{}

		if err := dao.DB.Where("user_id = ? AND video_id = ?", user_id, video_id).
			Delete(&favorite).Error; err != nil {
			return err
		}
	}

	err = dao.DB.Save(&user).Error
	if err != nil {
		return errors.New("用户点赞数修改失败")
	}

	err = dao.DB.Save(&video).Error
	if err != nil {
		return errors.New("视频点赞数修改失败")
	}

	err = dao.DB.Save(&author).Error
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
