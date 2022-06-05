package service

import (
	"douyin-server/dao"
	"path/filepath"
	"strings"
)

const staticRouter = "http://172.25.151.124:3000/"

func Publish(filename, title string, user_id int64, is_generate_ok bool) error {
	fileSuffix := filepath.Ext(filename)

	v := dao.Video{
		AuthorId: user_id,
		Title:    title,
		PlayUrl:  staticRouter + "videos/" + filename,
		CoverUrl: staticRouter + "covers/" + strings.TrimSuffix(filename, fileSuffix) + ".jpg",
	}

	if !is_generate_ok {
		v.CoverUrl = "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg"
	}

	if err := dao.DB.Create(&v).Error; err != nil {
		return err
	}

	return nil
}

func PublishList() (videoList []dao.Video) {
	dao.DB.Model(&dao.Video{}).Joins("left join users on videos.id = users.id").Order("created_at desc").Find(&videoList)
	return
}
