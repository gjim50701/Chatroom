package process

import (
	"encoding/json"
	"ex17/chatRoom/common/message"
	"fmt"
)

func outPutGroupMes(mes *message.Message) {

	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("json.Unmarshal error :", err)
		return
	}

	//格式化顯示訊息
	info := fmt.Sprintf("@%s對大家說: %s", smsMes.UserID, smsMes.Context)
	fmt.Println(info)
	fmt.Println()

}

func outPutSomeoneMes(mes *message.Message) {

	var smsMes message.SmsSomeoneMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("json.Unmarshal error :", err)
		return
	}

	//格式化顯示訊息
	info := fmt.Sprintf("@%s對@%s說: %s", smsMes.UserID, smsMes.SomeoneID, smsMes.Context)
	fmt.Println(info)
	fmt.Println()

}
