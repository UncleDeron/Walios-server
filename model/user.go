package model

import "github.com/aceld/zinx/ziface"

type User struct {
	Conn     ziface.IConnection
	UserId   string
	UserInfo UserInfo
}

type UserInfo struct {
	UserAccount  string         `json:"userAccount"`
	UserId       string         `json:"userId"`
	UserPwd      string         `json:"userPwd"`
	NickName     string         `json:"nickName"`
	Enable       bool           `json:"enable"`
	Status       UserStatusType `json:"status"`
	Avatar       string         `json:"avatar"`
	Introduction string         `json:"introduction"`
	Email        string         `json:"email"`
}

// UserStatusType 用户状态类型
type UserStatusType int

const (
	OffLine UserStatusType = iota
	OnLine
)
