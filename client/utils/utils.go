package utils

import (
	"encoding/binary"
	"encoding/json"
	"ex17/chatRoom/common/message"
	"fmt"
	"net"
)

//將方法關聯到結構體中

type Transfer struct {
	Conn net.Conn
	Buf  [8096]byte //到時當切片用
}

func (this *Transfer) ReadPkg() (mes message.Message, err error) {

	//buf := make([]byte, 8096)
	fmt.Println("讀取客戶端發送數據...")

	//conn.Read 在conn沒有被關閉的情況下 會產生阻塞 此時如果客戶端關閉了conn 就不會產生阻塞
	_, err = this.Conn.Read(this.Buf[:4])
	if err != nil {
		//err = errors.New("read pkg header error")
		return
	}

	//根據buf[:4] 轉成一個uint32類型 為了知道要讀多少個字節
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(this.Buf[:4])

	//根據pkgLen讀取消息內容
	n, err := this.Conn.Read(this.Buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		//err = errors.New("read pkg body error")
		return
	}

	//把pkg先反序列化成message.Message類型
	err = json.Unmarshal(this.Buf[:pkgLen], &mes) //記得加&
	if err != nil {
		fmt.Println("json.Unmarshal error :", err)
		return
	}

	return
}

func (this *Transfer) WritePkg(data []byte) (err error) {

	//發送一長度給對方
	var pkgLen uint32
	pkgLen = uint32(len(data))
	//var buf [4]byte
	binary.BigEndian.PutUint32(this.Buf[0:4], pkgLen)
	//正式發送長度
	n, err := this.Conn.Write(this.Buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write(buf) error :", err)
		return
	}

	//發送消息本身
	n, err = this.Conn.Write(data)
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Write(data) error :", err)
		return
	}

	return
}
