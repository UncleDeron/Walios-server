package controller

import (
	"github.com/UncleDeron/Walios/pb"
	"github.com/UncleDeron/Walios/protocol"
	"github.com/UncleDeron/frogChat-server/dao"
	"github.com/UncleDeron/frogChat-server/model"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"google.golang.org/protobuf/proto"
)

type UserController struct {
	Conn     ziface.IConnection
	UserInfo *model.UserInfo
}

func RegisterUser(conn ziface.IConnection, registerData *pb.RegisterMsgData) {
	var resData *pb.ResData
	// 尝试从 redis 中取出对应的用户
	if user, _ := dao.UserDaoInstance.GetUserByAccount(registerData.UserAccount); user != nil {
		// 用户已存在
		resData = &pb.ResData{
			Code: pb.ResponseCode_FAILED,
			Msg:  "该账号已存在",
			Data: nil,
		}
	} else {
		regData := &model.UserInfo{
			UserAccount:  registerData.UserAccount,
			UserPwd:      registerData.UserPwd,
			NickName:     registerData.NickName,
			Enable:       true,
			Status:       model.OffLine,
			Introduction: "",
			Email:        registerData.Email,
		}
		err := dao.UserDaoInstance.AddUser(regData)
		if err != nil {
			resData = &pb.ResData{
				Code: pb.ResponseCode_FAILED,
				Msg:  err.Error(),
				Data: nil,
			}
		} else {
			resData = &pb.ResData{
				Code: pb.ResponseCode_SUCCESS,
				Msg:  "注册成功",
				Data: nil,
			}
		}
	}
	msg, err := proto.Marshal(resData)
	if err != nil {
		zlog.Error("Msg Marshal failed: ", err)
		return
	}
	err = conn.SendMsg(protocol.RegisterResponse, msg)
	if err != nil {
		zlog.Error("SendMsg failed: ", err)
	}

}

func Login(conn ziface.IConnection, loginInfo *pb.LoginMsgData) {
	var resData *pb.ResData
	// 尝试从 redis 中取出对应的用户
	user, err := dao.UserDaoInstance.GetUserByAccount(loginInfo.UserAccount)
	if err != nil {
		resData = &pb.ResData{
			Code: pb.ResponseCode_FAILED,
			Msg:  err.Error(),
			Data: nil,
		}
		zlog.Info("login failed: ", err)
	} else if !user.Enable {
		resData = &pb.ResData{
			Code: pb.ResponseCode_FAILED,
			Msg:  "该用户已被封禁",
			Data: nil,
		}
	} else if loginInfo.UserPwd == user.UserPwd { // 校验密码通过
		userController := &UserController{
			Conn:     conn,
			UserInfo: user,
		}
		// 将 userController 添加到用户管理器的 user 列表中
		UserManagerInstance.AddUser(user.UserId, userController)
		zlog.Info("用户ID:", user.UserId, ": 登录成功")
		resData = &pb.ResData{
			Code: pb.ResponseCode_SUCCESS,
			Msg:  "登录成功",
			Data: &pb.ResData_UserInfo{
				UserInfo: &pb.UserInfo{
					UserAccount:  user.UserAccount,
					UserId:       user.UserId,
					UserPwd:      "",
					NickName:     user.NickName,
					Enable:       true,
					Status:       pb.UserStatus(model.OnLine),
					Avatar:       user.Avatar,
					Introduction: user.Introduction,
				},
			},
		}
	} else {
		resData = &pb.ResData{
			Code: pb.ResponseCode_FAILED,
			Msg:  "密码与用户名不匹配",
			Data: nil,
		}
	}
	msg, err := proto.Marshal(resData)
	if err != nil {
		zlog.Error("Msg Marshal failed: ", err)
		return
	}
	err = conn.SendMsg(protocol.LoginResponse, msg)
	if err != nil {
		zlog.Error("SendMsg failed: ", err)
	}
}
