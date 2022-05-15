package service

import (
	"bytes"
	"douyin-server/dao"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"math/rand"
	"strings"
	"time"
)

func Register(username, password string) (int64, error) {
	if len(username) > 32 {
		return 0, errors.New("用户名过长，不可超过32位")
	}
	if len(password) > 32 {
		return 0, errors.New("密码过长，不可超过32位")
	}

	user := dao.User{}
	dao.DB.Where("name = ?", username).Find(&user)
	if user.Id != 0 {
		return 0, errors.New("用户已存在")
	}

	user.Name = username

	// 加密存储用户密码
	user.Salt = randSalt()
	buf := bytes.Buffer{}
	buf.WriteString(username)
	buf.WriteString(password)
	buf.WriteString(user.Salt)
	pwd, err := bcrypt.GenerateFromPassword(buf.Bytes(), bcrypt.MinCost)
	if err != nil {
		return 0, err
	}
	user.Pwd = string(pwd)

	dao.DB.Create(&user)
	return user.Id, nil
}

func Login(username, password string) (int64, error) {
	user := dao.User{}
	err := dao.DB.Where("name = ?", username).Find(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, errors.New("用户不存在")
	} else if err = bcrypt.CompareHashAndPassword([]byte(user.Pwd), []byte(username+password+user.Salt)); err != nil {
		return 0, errors.New("用户名或密码错误")
	}
	return user.Id, nil
}

func UserInfo(token string) (dao.User, error) {
	user := dao.User{}
	err := dao.DB.Where("name = ?", token).Find(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return user, errors.New("用户不存在")
	}
	return user, nil
}

// 随机盐长度固定为4

func randSalt() string {
	rand.Seed(time.Now().UnixNano())
	buf := strings.Builder{}
	for i := 0; i < 4; i++ {
		// 如果写byte会无法兼容数据库编码
		buf.WriteRune(rune(rand.Intn(256)))
	}
	return buf.String()
}
