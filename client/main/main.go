package main

import (
	"ex17/chatRoom/client/process"
	"fmt"
)

var userID string
var userPwd string
var userName string

func main() {

	var key int
	//var loop = true

	for {
		fmt.Println("-----歡迎登入多人聊天系統-----")
		fmt.Println("\t 1 登入聊天室")
		fmt.Println("\t 2 註冊用戶")
		fmt.Println("\t 3 退出系統")
		fmt.Println("\t 請選擇(1-3):")

		fmt.Scanf("%d\n", &key)
		switch key {
		case 1:
			fmt.Println("登入聊天室")
			fmt.Println("請輸入用戶帳號:")
			fmt.Scanf("%s\n", &userID)
			fmt.Println("請輸入用戶密碼:")
			fmt.Scanf("%s\n", &userPwd)

			//完成登入
			//創建一UserProcess的實例
			up := &process.UserProcess{}
			up.Login(userID, userPwd)

		case 2:
			fmt.Println("註冊用戶")
			fmt.Println("請輸入用戶帳號:")
			fmt.Scanf("%s\n", &userID)
			fmt.Println("請輸入用戶密碼:")
			fmt.Scanf("%s\n", &userPwd)
			fmt.Println("請輸入用戶名稱:")
			fmt.Scanf("%s\n", &userName)

			//調用UserProcess 完成註冊要求
			up := &process.UserProcess{}
			up.Register(userID, userPwd, userName)

		case 3:
			fmt.Println("退出系統")
			//loop = false
		default:
			fmt.Println("輸入有誤 請重新輸入")
		}

	}
}
