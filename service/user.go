package service

import (
	"bytes"
	"douyin-server/dao"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// LoginLimit 中间件服务，限制注册登录操作过于频繁。
func LoginLimit(ipAddress string) bool {
	// 错误可忽略
	times, _ := dao.RDB.Get(dao.Ctx, ipAddress).Int64()
	if times > 10 {
		return false
	} else {
		dao.RDB.Set(dao.Ctx, ipAddress, times+1, time.Minute)
	}
	return true
}

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

func UserInfo(id int64) (dao.User, error) {
	user := dao.User{}
	err := dao.DB.Where("id = ?", id).Find(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return user, errors.New("用户不存在")
	}
	return user, nil
}

// CreateToken 生成随机token，并存储到redis中，返回token
func CreateToken(id int64) (token int64) {
	// redis存储64位整数更节省空间
	token = int64(rand.Uint64())

	// 检测token有无冲突
	_, err := dao.RDB.Get(dao.Ctx, strconv.FormatInt(token, 10)).Result()
	for err == nil {
		token = int64(rand.Uint64())
		_, err = dao.RDB.Get(dao.Ctx, strconv.FormatInt(token, 10)).Result()
	}

	dao.RDB.Set(dao.Ctx, strconv.FormatInt(token, 10), id, 12*time.Hour)

	return
}

// 随机盐长度固定为4
func randSalt() string {
	buf := strings.Builder{}
	for i := 0; i < 4; i++ {
		// 如果写byte会无法兼容mysql编码
		buf.WriteRune(rune(rand.Intn(256)))
	}
	return buf.String()
}
