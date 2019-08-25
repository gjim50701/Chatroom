package process

import (
	"encoding/json"
	"ex17/chatRoom/common/message"
	"ex17/chatRoom/server/model"
	"ex17/chatRoom/server/utils"
	"fmt"
	"net"
)

type UserProcess struct {
	Conn net.Conn
	//增加一屬性 表示該Conn是哪個用戶
	UserID string
}

//專門處理登入請求

func (this *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {

	//從mes中取出mes.Data 並反序列化成LoginMes
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail error :", err)
		return
	}

	var resMes message.Message
	resMes.Type = message.LoginResMesType

	var loginResMes message.LoginResMes

	//登入需要到redis數據庫完成驗證
	user, err := model.MyUserDao.Login(loginMes.UserID, loginMes.UserPwd)

	if err != nil {
		if err == model.ERROR_USER_NOEXISTS {
			loginResMes.CodeID = 500
			loginResMes.Error = err.Error()
		} else if err == model.ERROR_USER_PWD {
			loginResMes.CodeID = 403
			loginResMes.Error = err.Error()
		} else {
			loginResMes.CodeID = 505
			loginResMes.Error = "服務器內部錯誤..."
		}
	} else {
		loginResMes.CodeID = 200

		//用戶登入成功 把該用戶放入到userMgr中
		//將登入的用戶ID賦給this
		this.UserID = loginMes.UserID
		userMgr.AddOnlineUser(this)

		//通知其他在線用戶上線
		this.NotifyOthersOnlineUser(loginMes.UserID, message.UserOnline)

		//將當前在線的用戶id 放入到loginResMes.UsersID
		//遍歷userMgr.onlineUsers
		for id := range userMgr.onlineUsers {
			loginResMes.UsersID = append(loginResMes.UsersID, id)
		}

		fmt.Println(user.UserName, "登入成功!")

	}

	//將loginResMes序列化
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("json.Marshal fail error :", err)
		return
	}

	//將data賦值給 resMes
	resMes.Data = string(data)

	//將resMes序列化
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal fail error :", err)
		return
	}

	//發送data 將其封裝到writePkg函數
	//使用分層模式(MVC) 先創建一實例 然後讀取
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)

	return
}

func (this *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {

	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail error :", err)
		return
	}

	var resMes message.Message
	resMes.Type = message.RegisterResMesType

	var registerResMes message.RegisterResMes

	err = model.MyUserDao.Register(&registerMes.User)
	if err != nil {
		if err == model.ERROR_USER_EXISTS {
			registerResMes.CodeID = 400
			registerResMes.Error = model.ERROR_USER_EXISTS.Error()
		} else {
			registerResMes.CodeID = 404
			registerResMes.Error = "註冊發生未知錯誤"
		}
	} else {
		registerResMes.CodeID = 200

	}

	data, err := json.Marshal(registerResMes)
	if err != nil {
		fmt.Println("json.Marshal fail error :", err)
		return
	}

	resMes.Data = string(data)

	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal fail error :", err)
		return
	}

	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)

	return
}

//通知其他在線用戶上線
func (this *UserProcess) NotifyOthersOnlineUser(userID string, status int) {

	//遍歷onlineUsers 然後一一發送NotifyUserStatusMes
	for id, up := range userMgr.onlineUsers {
		//過濾自己
		if id == userID {
			continue
		}
		up.NotifyMyOnline(userID, status)

		//開始通知 另外寫個方法
	}
}

func (this *UserProcess) NotifyMyOnline(userID string, status int) {

	//組裝NotifyUserStatusMes消息
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType

	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserID = userID
	notifyUserStatusMes.Status = status

	//notifyUserStatusMes序列化
	data, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("json.Marshal fail error :", err)
		return
	}

	mes.Data = string(data)

	//再對mes序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal fail error :", err)
		return
	}

	//進行發送 創建Transfer實例
	tf := &utils.Transfer{
		Conn: this.Conn,
	}

	tf.WritePkg(data)

	if err != nil {
		fmt.Println("NotifyMyOnline fail error :", err)
		return
	}

}

func (this *UserProcess) Exit(mes *message.Message) (err error) {

	var statusMes message.NotifyUserStatusMes
	err = json.Unmarshal([]byte(mes.Data), &statusMes)
	if err != nil {
		fmt.Println("json.Unmarshal error:", err)
		return
	}

	userMgr.DelOnlineUser(statusMes.UserID)

	this.NotifyOthersOnlineUser(statusMes.UserID, statusMes.Status)

	return
}
