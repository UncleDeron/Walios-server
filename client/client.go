package main

import (
	"encoding/json"
	"fmt"
	"github.com/UncleDeron/frogChat-server/model"
	"io"
	"net"
	"time"

	"github.com/aceld/zinx/znet"
)

/*
	模拟客户端
*/
func main() {

	conn, err := net.Dial("tcp", "127.0.0.1:23746")
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}

	for {
		//发封包message消息
		dp := znet.NewDataPack()
		data := model.ClientMsg{
			Action: model.Login,
			Data: model.LoginMsgData{
				Username: "chenhda",
				Password: "123456",
			},
		}
		msg, err := json.Marshal(data)
		if err != nil {
			fmt.Println("json marshal err, exit!", err)
			continue
		}
		msg, _ = dp.Pack(znet.NewMsgPackage(0, msg))
		_, err = conn.Write(msg)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		//先读出流中的head部分
		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData) //ReadFull 会把msg填充满为止
		if err != nil {
			fmt.Println("read head error")
			break
		}
		//将headData字节流 拆包到msg中
		msgHead, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("server unpack err:", err)
			return
		}

		if msgHead.GetDataLen() > 0 {
			//msg 是有data数据的，需要再次读取data数据
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetDataLen())

			//根据dataLen从io中读取字节流
			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("server unpack data err:", err)
				return
			}

			fmt.Println("==> Test Router:[Ping] Recv Msg: ID=", msg.ID, ", len=", msg.DataLen, ", data=", string(msg.Data))
		}

		time.Sleep(1 * time.Second)
	}
}
