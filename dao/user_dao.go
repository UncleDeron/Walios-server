package dao

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/UncleDeron/frogChat-server/model"
	"github.com/garyburd/redigo/redis"
	"github.com/google/uuid"
)

var (
	UserDaoInstance *UserDao
)

// UserDao 结构体
type UserDao struct {
	Pool *redis.Pool
}

// NewUserDao 使用工厂模式创建 UserDao 实例
func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	userDao = &UserDao{
		Pool: pool,
	}
	return
}

// GetUserByAccount 根据 userAccount 返回一个用户
func (dao *UserDao) GetUserByAccount(account string) (user *model.UserInfo, err error) {
	user = &model.UserInfo{}
	// 通过 userAccount 去redis 查讯用户
	conn := dao.Pool.Get()
	defer conn.Close()
	res, err := redis.String(conn.Do("hget", "users", account))
	if err != nil {
		if err == redis.ErrNil { // 没有找到该 account 对应的用户
			err = errors.New("用户不存在")
		}
		return nil, err
	}

	// 将 res 解析为 user 实例
	err = json.Unmarshal([]byte(res), user)
	if err != nil {
		fmt.Println("解析 user 出错:", err)
		return nil, err
	}
	return user, nil
}

func (dao *UserDao) AddUser(registerData *model.UserInfo) (err error) {
	conn := dao.Pool.Get()
	defer conn.Close()
	user := &model.UserInfo{
		UserAccount:  registerData.UserAccount,
		UserId:       uuid.NewString(),
		UserPwd:      registerData.UserPwd,
		NickName:     registerData.NickName,
		Email:        registerData.Email,
		Introduction: "",
		Avatar:       "default.jpg",
		Status:       model.OffLine,
		Enable:       true,
	}
	data, err := json.Marshal(user)
	if err != nil {
		err = errors.New("序列化新用户失败: " + err.Error())
	}
	_, err = conn.Do("hset", "users", user.UserAccount, data)
	if err != nil {
		fmt.Println(err)
	}
	return
}
