package model

import "errors"

//根據業務邏輯的需要 自定義一些錯誤

var (
	ERROR_USER_NOEXISTS = errors.New("用戶不存在...")
	ERROR_USER_EXISTS   = errors.New("用戶已存在...")
	ERROR_USER_PWD      = errors.New("密碼不正確...")
)
