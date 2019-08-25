package model

import (
	"ex17/chatRoom/common/message"
	"net"
)

type CurUser struct {
	Conn net.Conn
	message.User
}
