package message

const (
	LoginMesType            = "LoginMes"
	LoginResMesType         = "LoginResMes"
	RegisterMesType         = "RegisterMes"
	RegisterResMesType      = "RegisterResMes"
	NotifyUserStatusMesType = "NotifyUserStatusMes"
	SmsMesType              = "SmsMes"
	SmsSomeoneMesType       = "SmsSomeoneMes"
)

//定義用戶狀態常量
const (
	UserOnline = iota
	UserOffline
)

type Message struct {
	Type string `json:"type"` //消息類型
	Data string `json:"data"` //消息內容

}

//定義消息 後面可再追加...

type LoginMes struct {
	UserID   string `json:"userID"`
	UserPwd  string `json:"userPwd"`
	UserName string `json:"userName"`
}

type LoginResMes struct {
	CodeID  int      `json:"codeID"`  //返回狀態碼  500:尚未註冊 200:登入成功
	UsersID []string `json:"usersID"` //保存用戶在線的id切片
	Error   string   `json:"error"`   //返回錯誤訊息
}

type RegisterMes struct {
	User User `json:"user"`
}

type RegisterResMes struct {
	CodeID int    `json:"codeID"` //返回狀態碼  400:已有用戶 200:註冊成功
	Error  string `json:"error"`
}

//為了配合服務器端推送用戶狀態變化的消息
type NotifyUserStatusMes struct {
	UserID string `json:"userID"`
	Status int    `json:"status"`
}

type SmsMes struct {
	Context string        `json:"context"`
	User    `json:"user"` //使用匿名結構體 繼承
}

type SmsSomeoneMes struct {
	SmsMes
	SomeoneID string `json:"someoneID"`
}
