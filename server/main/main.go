package main

import (
	"ex17/chatRoom/server/model"
	"fmt"
	"net"
	"time"
)

//處理和客戶端的通訊
func process01(conn net.Conn) {

	defer conn.Close()

	//調用總控 先創建實例
	processor := &Processor{
		Conn: conn,
	}
	err := processor.process02()
	if err != nil {
		fmt.Println("processor.process2() error :", err)
		return
	}

}

func init() {

	//初始化鏈接池 並完成對userDao的初始化 (注意初始化的順序)
	initPool("localhost:6379", 16, 0, 300*time.Second)
	initUserDao()
}

//完成對userDao的初始化
func initUserDao() {
	//pool是在redis.go定的全局變量
	model.MyUserDao = model.NewUserDao(pool)
}

func main() {

	fmt.Println("服務器[新結構]在8889端口接聽...")
	listen, err := net.Listen("tcp", "127.0.0.1:8889")
	if err != nil {
		fmt.Println("net.Listen error :", err)
		return
	}

	defer listen.Close()

	//監聽成功就等待客戶端鏈接服務器
	for {
		fmt.Println("等待客戶端連接...")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept() error :", err)
		}

		//鏈接成功就啟動一協程和客戶端保持通訊
		go process01(conn)
	}
}
