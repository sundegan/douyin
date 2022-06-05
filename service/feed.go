package service

import (
	"douyin-server/dao"
	"time"
)

// Feed 选择发布时间在latestTime之前的视频
func Feed(id int64, latestTime int64) (videoList []dao.Video, nextTime int64) {
	timeStamp := time.UnixMilli(latestTime)
	dao.DB.Model(&dao.Video{}).Preload("Author").Where("created_at <= ?", timeStamp).Order("created_at desc").Limit(30).Find(&videoList)

	eraseSensitiveField(&videoList)

	// 没有视频，返回空
	if len(videoList) == 0 {
		return
	}
	// 返回这次视频最近的投稿时间-1，下次即可获取比这次视频旧的视频
	nextTime = videoList[len(videoList)-1].CreatedAt.UnixMilli() - 1

	// 说明当前获取视频的客户端已登陆账号
	if id != 0 {
		for i := range videoList {
			rows, err := dao.DB.Model(&dao.Favorite{}).Where("user_id = ? AND video_id = ?", id, videoList[i].Id).Rows()
			if err != nil {
				continue
			}
			if rows.Next() {
				videoList[i].IsFavorite = true
			}
		}
	}
	return
}

// 去除视频中作者的敏感信息
func eraseSensitiveField(videoList *[]dao.Video) {
	for i := range *videoList {
		author := &(*videoList)[i].Author
		author.Id = 0
		author.Pwd = ""
		author.Salt = ""
	}
}
