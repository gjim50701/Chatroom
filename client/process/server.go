package process

import (
	"bufio"
	"encoding/json"
	"ex17/chatRoom/client/utils"
	"ex17/chatRoom/common/message"
	"fmt"
	"net"
	"os"
)

//ShowMenu顯示登入成功後的介面
func ShowMenu() {

	var key int
	var oneOrAll int
	var userID string

	smsProcess := &SmsProcess{}
	inputReader := bufio.NewReader(os.Stdin)

	fmt.Println("\t 1 顯示在線用戶")
	fmt.Println("\t 2 發送訊息")
	fmt.Println("\t 3 訊息列表")
	fmt.Println("\t 4 登出系統")
	fmt.Println("\t 請選擇(1-4):")

	fmt.Scanf("%d\n", &key)
	switch key {
	case 1:
		outPutOnlineUsers()
	case 2:
		fmt.Println("請選擇 1 向某人發送訊息 2 向群體發送訊息:")
		fmt.Scanf("%d\n", &oneOrAll)
		if oneOrAll == 1 {
			fmt.Println("請輸入想發送對象:")
			fmt.Scanf("%s\n", &userID)
			fmt.Println("請輸入想發送訊息:")
			context, err := inputReader.ReadString('\n')
			if err != nil {
				fmt.Println("輸入內容時有誤 error:", err)
			}
			smsProcess.SendSomeoneMes(userID, context)
		} else if oneOrAll == 2 {
			fmt.Println("請輸入想發送的訊息:")
			context, err := inputReader.ReadString('\n')
			if err != nil {
				fmt.Println("輸入內容時有誤 error:", err)
			}
			smsProcess.SendGroupMes(context)
		}
	case 3:
		fmt.Println("訊息列表")
	case 4:
		exit()
		fmt.Println("已登出系統")
		os.Exit(0)
	default:
		fmt.Println("輸入有誤")
	}

}

func serverProcessMes(conn net.Conn) {

	//創建一個Transfer實例 使它不停讀取服務器
	tf := &utils.Transfer{
		Conn: conn,
	}

	for {
		fmt.Println("客戶端正在讀取服務器發送的訊息...")
		mes, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("tf.ReadPkg() error :", err)
			return
		}

		//如果讀取到了訊息 又到下一個處理環節
		switch mes.Type {
		case message.NotifyUserStatusMesType:

			//取出NotifyUserStatusMes
			var notifyUserStatusMes message.NotifyUserStatusMes
			json.Unmarshal([]byte(mes.Data), &notifyUserStatusMes)
			//把該用戶狀態保存到客戶端map
			updataUserStatus(&notifyUserStatusMes)
		case message.SmsMesType:
			outPutGroupMes(&mes)
		case message.SmsSomeoneMesType:
			outPutSomeoneMes(&mes)
		default:
			fmt.Println("服務器端返回一未知消息類型")
		}
	}
}
