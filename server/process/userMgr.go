package process

import "fmt"

//由於userMgr實例在服務器中只有一個 且多處會用到 因此定為全局變量
var (
	userMgr *UserMgr
)

type UserMgr struct {
	onlineUsers map[string]*UserProcess
}

//完成對userMgr的初始化
func init() {
	userMgr = &UserMgr{
		onlineUsers: make(map[string]*UserProcess, 1024),
	}

}

//完成對onlineUsers添加
func (this *UserMgr) AddOnlineUser(up *UserProcess) {

	this.onlineUsers[up.UserID] = up

}

//登出時除去登入列表名單
func (this *UserMgr) DelOnlineUser(userID string) {

	delete(this.onlineUsers, userID)

}

//查詢時返回當前在線用戶列表
func (this *UserMgr) GetAllOnlineUser() map[string]*UserProcess {

	return this.onlineUsers

}

//根據ID回傳對應值
func (this *UserMgr) GetOnlineUserByID(userID string) (up *UserProcess, err error) {

	up, ok := this.onlineUsers[userID]
	if !ok {
		//查找的ID用戶當前不在線
		err = fmt.Errorf("用戶%s不在線...", userID)
		return
	}

	return

}
