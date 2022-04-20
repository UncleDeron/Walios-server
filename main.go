package main

import (
	"github.com/UncleDeron/frogChat-server/dao"
	"github.com/UncleDeron/frogChat-server/model"
	router "github.com/UncleDeron/frogChat-server/router"
	"github.com/UncleDeron/frogChat-server/utils"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
)

// DoConnectionBegin 创建连接的时候执行
func DoConnectionBegin(conn ziface.IConnection) {
	zlog.Debug("DoConnecionBegin is Called ... ")

	//设置两个链接属性，在连接创建之后

	err := conn.SendMsg(2, []byte("DoConnection BEGIN..."))
	if err != nil {
		zlog.Error(err)
	}

}

// DoConnectionLost 连接断开的时候执行
func DoConnectionLost(conn ziface.IConnection) {
	//在连接销毁之前，查询conn的Name，Home属性

	zlog.Debug("DoConneciotnLost is Called ... ")
}

func main() {
	// 初始化线程池
	conf := &utils.RedisConfig{}
	err := utils.LoadConfig("./conf/redis.json", conf)
	if err != nil {
		zlog.Error("线程池配置读取失败")
	}
	pool := initPool(conf)
	// 初始化 userDao
	dao.UserDaoInstance = dao.NewUserDao(pool)

	// 初始化zinx
	s := znet.NewServer()
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)
	s.AddRouter(uint32(model.ClientT), &router.ClientRouter{})
	s.Serve()

}
