package process

import (
	"encoding/json"
	"ex17/chatRoom/common/message"
	"ex17/chatRoom/server/utils"
	"fmt"
	"net"
)

type SmsProcess struct {
}

func (this *SmsProcess) SendGroupMes(mes *message.Message) {

	//取出mes的內容
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("json.Unmarshal error:", err)
		return
	}

	data, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal error:", err)
		return
	}

	//遍歷服務器端的map 將訊息轉發出去
	for id, up := range userMgr.onlineUsers {
		if id == smsMes.UserID {
			continue
		}
		this.SendEachOnlineUser(data, up.Conn)
	}

}

func (this *SmsProcess) SendSomeoneMes(mes *message.Message) {

	//取出mes的內容
	var smsMes message.SmsSomeoneMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("json.Unmarshal error:", err)
		return
	}

	data, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal error:", err)
		return
	}

	up := userMgr.onlineUsers[smsMes.SomeoneID]

	this.SendEachOnlineUser(data, up.Conn)

}

func (this *SmsProcess) SendEachOnlineUser(data []byte, conn net.Conn) {

	tf := &utils.Transfer{
		Conn: conn,
	}

	err := tf.WritePkg(data)
	if err != nil {
		fmt.Println("轉發消息失敗 error:", err)
		return
	}
}
