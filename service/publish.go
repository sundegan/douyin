package service

import (
	"douyin-server/dao"
	"path/filepath"
	"strings"
)

const staticRouter = "http://122.9.113.111:8080/"

func Publish(filename, title string, userId int64, isGenerateOk bool) error {
	fileSuffix := filepath.Ext(filename)

	v := dao.Video{
		AuthorId: userId,
		Title:    title,
		PlayUrl:  staticRouter + "videos/" + filename,
		CoverUrl: staticRouter + "covers/" + strings.TrimSuffix(filename, fileSuffix) + ".jpg",
	}

	if !isGenerateOk {
		v.CoverUrl = "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg"
	}

	if err := dao.DB.Create(&v).Error; err != nil {
		return err
	}
	return nil
}

func PublishList(id int64) (videoList []dao.Video) {
	dao.DB.Model(&dao.Video{}).Where("author_id = ?", id).Order("created_at DESC").Find(&videoList)
	return
}
