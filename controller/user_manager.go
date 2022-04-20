package controller

import (
	"sync"
)

type UserManager struct {
	// 当前在线的User集合
	Users map[string]*UserController
	// 保护 Users 的锁
	uLock sync.RWMutex
}

// UserManagerInstance 全局 UserManager
var UserManagerInstance *UserManager

// 初始化 UserManager
func init() {
	UserManagerInstance = &UserManager{
		Users: make(map[string]*UserController),
	}
}

// AddUser 添加一个 user
func (um *UserManager) AddUser(userId string, user *UserController) {
	um.uLock.Lock()
	um.Users[userId] = user
	um.uLock.Unlock()
}
