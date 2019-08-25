package process

import (
	"encoding/binary" //實現簡單的數字與字節序列的轉換方式
	"encoding/json"
	"ex17/chatRoom/client/utils"
	"ex17/chatRoom/common/message"
	"fmt"
	"net"
	"os"
)

type UserProcess struct {
}

func (this *UserProcess) Login(userID, userPwd string) (err error) {

	//連接到服務器端
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial error :", err)
		return
	}

	defer conn.Close()

	//準備透過conn發送消息給服務器
	var mes message.Message
	mes.Type = message.LoginMesType

	//創建loginMes結構體
	var loginMes message.LoginMes
	loginMes.UserID = userID
	loginMes.UserPwd = userPwd

	//將loginMes序列化再給mes.Data
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("json.Marshal error :", err)
		return
	}
	mes.Data = string(data)

	//將mes序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal error :", err)
		return
	}

	//先把data長度發送給服務器(怕丟包)
	//不能直接用conn.Write() 因為要傳的是[]byte (Write(b []byte) (n int, err error))
	var pkgLen uint32
	pkgLen = uint32(len(data))
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[0:4], pkgLen)
	//正式發送長度
	n, err := conn.Write(buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write(buf) error :", err)
		return
	}

	//發送消息本身
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write(data) error :", err)
		return
	}

	//處理服務器端返回的消息
	tf := &utils.Transfer{
		Conn: conn,
	}
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("tf.ReadPkg() error :", err)
		return
	}

	//將mes.Data反序列化成 LoginResMes
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)

	if loginResMes.CodeID == 200 {

		//登入成功

		//初始化curUser
		curUser.Conn = conn
		curUser.UserID = userID
		curUser.UserStatus = message.UserOnline

		//顯示當前在線用戶列表 => 遍歷loginResMes.UserID
		fmt.Println("-----當前在線用戶列表-----")
		for _, v := range loginResMes.UsersID {

			if v == userID {
				continue
			}
			fmt.Printf("@%s\n", v)

			//完成客戶端onlineUsers初始化
			user := &message.User{
				UserID:     v,
				UserStatus: message.UserOnline,
			}
			onlineUsers[v] = user
		}
		fmt.Printf("\n\n")

		//還需在客戶端啟動一協程 該協程保持和服務器端的通訊
		//如果有服務器有數據推送給客戶端 則顯示並顯示在客戶端的終端
		go serverProcessMes(conn)

		//顯示登入成功的介面
		fmt.Printf("-----歡迎用戶%s登入成功-----\n", userID)
		for {
			ShowMenu()
		}

	} else {
		fmt.Println(loginResMes.Error)
	}

	return

}

func (this *UserProcess) Register(userID, userPwd, userName string) (err error) {

	//連接到服務器端
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial error :", err)
		return
	}

	defer conn.Close()

	//準備透過conn發送消息給服務器
	var mes message.Message
	mes.Type = message.RegisterMesType

	//創建RegisterMes結構體
	var registerMes message.RegisterMes
	registerMes.User.UserID = userID
	registerMes.User.UserPwd = userPwd
	registerMes.User.UserName = userName

	data, err := json.Marshal(registerMes)
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
		Conn: conn,
	}

	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("註冊發送訊息錯誤 error :", err)
		return
	}

	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("tf.ReadPkg() error :", err)
		return
	}

	var registerResMes message.RegisterResMes
	err = json.Unmarshal([]byte(mes.Data), &registerResMes)

	if registerResMes.CodeID == 200 {

		fmt.Println("註冊成功 可重新登入...")
		os.Exit(0)
	} else {
		fmt.Println(registerResMes.Error)
		os.Exit(0)
	}

	return

}
