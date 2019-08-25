package model

import (
	"encoding/json"
	"ex17/chatRoom/common/message"
	"fmt"

	"github.com/garyburd/redigo/redis"
)

//在服務器啟動時 就初始化一個userDao實例
//做成全局變量 在需要和redis操作時 即可直接使用
var (
	MyUserDao *UserDao
)

//定義一個userDao的結構體 完成對user結構體的各種操作
type UserDao struct {
	pool *redis.Pool
}

//使用工廠模式 創建一userDao的實例
func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	userDao = &UserDao{
		pool: pool,
	}
	return
}

//根據用戶ID 返回一User實例&error
func (this *UserDao) getUserByID(conn redis.Conn, id string) (user *User, err error) {

	//透過給定的ID 去redis查詢該用戶
	res, err := redis.String(conn.Do("HGet", "users", id))
	if err != nil {
		if err == redis.ErrNil { //沒有找到對應的ID
			err = ERROR_USER_NOEXISTS
		}
		return
	}

	user = &User{}

	//把回傳的res反序列化成User實例
	err = json.Unmarshal([]byte(res), &user)
	if err != nil {
		fmt.Println("json.Unmarshal err:", err)
		return
	}

	return

}

//Login完成對用戶的驗證 id&pwd都正確 返回一個user實例 其中有誤 則返回錯誤訊息
func (this *UserDao) Login(userID, userPwd string) (user *User, err error) {

	//先從userDao的鏈接池裡取出一鏈接
	conn := this.pool.Get()
	defer conn.Close()

	user, err = this.getUserByID(conn, userID)
	if err != nil {
		return
	}

	if user.UserPwd != userPwd {
		err = ERROR_USER_PWD
		return
	}

	return

}

func (this *UserDao) Register(user *message.User) (err error) {

	//先從userDao的鏈接池裡取出一鏈接
	conn := this.pool.Get()
	defer conn.Close()

	_, err = this.getUserByID(conn, user.UserID)
	if err == nil {
		err = ERROR_USER_EXISTS
		return
	}

	//此時該用戶ID尚未在redis註冊過
	data, err := json.Marshal(user)
	if err != nil {
		fmt.Println("json.Marshal err:", err)
		return
	}

	//資料傳入redis
	_, err = conn.Do("HSet", "users", user.UserID, string(data))
	if err != nil {
		fmt.Println("保存註冊用戶錯誤 err:", err)
		return
	}
	return
}
