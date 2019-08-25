package main

import (
	"ex17/chatRoom/common/message"
	"ex17/chatRoom/server/process"
	"ex17/chatRoom/server/utils"
	"fmt"
	"io"
	"net"
)

type Processor struct {
	Conn net.Conn
}

//根據客戶端發送消息的種類 決定調用哪個函數處理
func (this *Processor) serverProcessMes(mes *message.Message) (err error) {

	//群聊測試 看是否有能接收客戶端傳來的訊息
	fmt.Println("mes =", mes)

	switch mes.Type {
	case message.LoginMesType:
		//處理登入邏輯
		//創建一個UserProcess實例
		up := &process.UserProcess{
			Conn: this.Conn,
		}
		err = up.ServerProcessLogin(mes)
	case message.RegisterMesType:
		//註冊處理
		up := &process.UserProcess{
			Conn: this.Conn,
		}
		err = up.ServerProcessRegister(mes)
	case message.SmsMesType:
		smsProcess := &process.SmsProcess{}
		smsProcess.SendGroupMes(mes)
	case message.SmsSomeoneMesType:
		smsProcess := &process.SmsProcess{}
		smsProcess.SendSomeoneMes(mes)
	case message.NotifyUserStatusMesType:
		up := &process.UserProcess{
			Conn: this.Conn,
		}
		err = up.Exit(mes)
	default:
		fmt.Println("消息類型不存在 無法處理...")
	}
	return
}

func (this *Processor) process02() (err error) {

	for {
		//將讀取數據包 封裝成一函數readPkg() 返回 Message ,Err
		//同樣創建一Transfer實例
		tf := &utils.Transfer{
			Conn: this.Conn,
		}
		mes, err := tf.ReadPkg()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客戶端退出了 服務器也退出...")
				return err
			}
			fmt.Println("readPkg(conn) error :", err)
			return err
		}

		err = this.serverProcessMes(&mes)
		if err != nil {
			return err
		}

	}

}
