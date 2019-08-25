package message

//先定義一個用戶的結構體
type User struct {
	UserID     string `json:"userID"`
	UserPwd    string `json:"userPwd"`
	UserName   string `json:"userName"`
	UserStatus int    `json:"userStatus"` //用戶狀態
}
