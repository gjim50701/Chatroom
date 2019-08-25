package process

import (
	"encoding/json"
	"ex17/chatRoom/client/model"
	"ex17/chatRoom/client/utils"
	"ex17/chatRoom/common/message"
	"fmt"
)

//客戶端自己的map(在線名單)
var onlineUsers = make(map[string]*message.User, 10)

//因為在客戶端的很多地方會使用到curUser 將其作為一全局 並在用戶成功登入後 完成初始化
var curUser model.CurUser

//在客戶端顯示當前在線的用戶列表
func outPutOnlineUsers() {

	fmt.Printf("\n")
	fmt.Println("-----當前在線用戶列表-----")
	for id := range onlineUsers {

		fmt.Printf("@%s\n", id)
	}
	fmt.Printf("\n\n")
}

//處理返回NotifyUserStatusMes
func updataUserStatus(notifyUserStatusMes *message.NotifyUserStatusMes) {

	//適當優化一下(有可能在線的狀態改變)
	user, ok := onlineUsers[notifyUserStatusMes.UserID]
	if !ok {
		user = &message.User{
			UserID:     notifyUserStatusMes.UserID,
			UserStatus: notifyUserStatusMes.Status,
		}

	}

	user.UserStatus = notifyUserStatusMes.Status

	if user.UserStatus == message.UserOnline {
		onlineUsers[notifyUserStatusMes.UserID] = user
	} else if user.UserStatus == message.UserOffline {
		delete(onlineUsers, notifyUserStatusMes.UserID)
		fmt.Printf("@%s已登出系統\n", notifyUserStatusMes.UserID)
	}

	outPutOnlineUsers()
}

func exit() {

	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType

	var statusMes message.NotifyUserStatusMes
	statusMes.UserID = curUser.UserID
	statusMes.Status = message.UserOffline

	data, err := json.Marshal(statusMes)
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

}
