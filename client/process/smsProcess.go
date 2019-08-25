package process

import (
	"encoding/json"
	"ex17/chatRoom/client/utils"
	"ex17/chatRoom/common/message"
	"fmt"
)

type SmsProcess struct {
}

//發送群聊消息
func (this *SmsProcess) SendGroupMes(context string) (err error) {

	var mes message.Message
	mes.Type = message.SmsMesType

	var smsMes message.SmsMes
	smsMes.Context = context
	smsMes.UserID = curUser.UserID
	smsMes.UserStatus = curUser.UserStatus

	data, err := json.Marshal(smsMes)
	if err != nil {
		fmt.Println("json.Marshal error :", err)
		return
	}
	mes.Data = string(data)

	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal error :", err)
		return
	}

	tf := &utils.Transfer{
		Conn: curUser.Conn,
	}

	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("註冊發送訊息錯誤 error :", err)
		return
	}

	return

}

//發送私聊訊息
func (this *SmsProcess) SendSomeoneMes(userID, context string) (err error) {

	var mes message.Message
	mes.Type = message.SmsSomeoneMesType

	var smsMes message.SmsSomeoneMes
	smsMes.SomeoneID = userID
	smsMes.SmsMes.Context = context
	smsMes.SmsMes.UserID = curUser.UserID
	smsMes.SmsMes.UserStatus = curUser.UserStatus

	data, err := json.Marshal(smsMes)
	if err != nil {
		fmt.Println("json.Marshal error :", err)
		return
	}
	mes.Data = string(data)

	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal error :", err)
		return
	}

	tf := &utils.Transfer{
		Conn: curUser.Conn,
	}

	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("註冊發送訊息錯誤 error :", err)
		return
	}

	return

}
