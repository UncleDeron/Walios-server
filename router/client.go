package router

import (
	"fmt"
	"github.com/UncleDeron/Walios/pb"
	"github.com/UncleDeron/frogChat-server/controller"
	"github.com/aceld/zinx/znet"
	"google.golang.org/protobuf/proto"
)
import "github.com/aceld/zinx/ziface"

type ClientRouter struct {
	znet.BaseRouter
}

func (r *ClientRouter) Handle(request ziface.IRequest) {
	fmt.Println("receive msg from client:", string(request.GetData()))
	data := request.GetData()
	clientMsg := &pb.ClientMsg{}
	err := proto.Unmarshal(data, clientMsg)
	if err != nil {
		fmt.Println("proto unmarshal error:", err)
		return
	}
	switch clientMsg.Action {
	case pb.ClientActionType_LOGIN:
		controller.Login(request.GetConnection(), clientMsg.GetLoginMsgData())
	case pb.ClientActionType_REGISTER:
		controller.RegisterUser(request.GetConnection(), clientMsg.GetRegisterMsgData())
	}
}
