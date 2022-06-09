package service

import (
	"douyin-server/dao"
)

const staticRouter = "http://122.9.113.111:8080/"

// Publish 向数据库插入视频的记录
func Publish(videoName, coverName, title string, userId int64, isGenerateOk bool) error {

	v := dao.Video{
		AuthorId: userId,
		Title:    title,
		PlayUrl:  staticRouter + "videos/" + videoName,
		CoverUrl: staticRouter + "covers/" + coverName,
	}

	// 若生成封面失败，视频的封面地址会被替换为默认封面的地址
	if !isGenerateOk {
		v.CoverUrl = "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg"
	}

	if err := dao.DB.Create(&v).Error; err != nil {
		return err
	}
	return nil
}

// PublishList 查询作者id和传入的id相同的视频记录，并倒序列出
func PublishList(id int64) (videoList []dao.Video) {
	dao.DB.Model(&dao.Video{}).Where("author_id = ?", id).Order("created_at DESC").Find(&videoList)
	return
}
