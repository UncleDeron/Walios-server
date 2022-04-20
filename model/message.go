package model

type MsgType int

type ResponseStatus int

const (
	ClientT MsgType = iota
	UserInfoT
	TextMsgT
)

type ClientActionType int

const (
	Login ClientActionType = iota
	Logout
)

const (
	loginsuccess ResponseStatus = iota
)

type LoginMsgData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ClientMsg struct {
	Action ClientActionType `json:"action"`
	Data   LoginMsgData     `json:"data"`
}

type loginResMsg struct {
}
