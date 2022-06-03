package service

import (
	"douyin-server/dao"
	"time"
)

// Feed 选择发布时间在latestTime之前的视频
func Feed(latestTime int64) (videoList []dao.Video, nextTime int64) {
	timeStamp := time.UnixMilli(latestTime)
	dao.DB.Model(&dao.Video{}).Preload("Author").Where("created_at <= ?", timeStamp).Limit(30).Find(&videoList)

	eraseSensitiveField(&videoList)

	// 没有视频，返回空
	if len(videoList) == 0 {
		return
	}
	// 返回这次视频最近的投稿时间-1，下次即可获取比这次视频旧的视频
	nextTime = videoList[0].CreatedAt.UnixMilli() - 1
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
