package model

import "server_template/mytime"

type User struct {
	ID int64 `gorm:"id" json:"id"`
	UserInfoReadOnly
	Password  string                `gorm:"password;type:char(64)"`
	Salt      string                `gorm:"salt;type:char(6)"`
	CreatedAt mytime.CustomTime     `json:"-"`
	UpdatedAt mytime.CustomTime     `json:"-"`
	DeletedAt mytime.CustomNullTime `json:"-" gorm:"index"`
}

type UserInfo struct {
	Nickname string `json:"nickname" gorm:"nickname;type:varchar(255)"`
	Avatar   string `json:"avatar" gorm:"avatar;type:varchar(2048)"`
}

type UserInfoReadOnly struct {
	Email    string `gorm:"email;type:varchar(255)" json:"email"`
	Phone    string `gorm:"phone;type:varchar(25)" json:"phone"`
	Username string `gorm:"column:username;unique;type:varchar(255)" json:"username"`
	Gender   rune   `gorm:"gender;type:char(1)" json:"gender"`
	UserInfo
}

type UseCodeLoginParameters struct {
	Username         string `json:"username" binding:"required"`
	VerificationCode string `json:"code" binding:"required"`
}

type LoginParameters struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegistrationParameters struct {
	Username         string `json:"username" binding:"required"`
	Password         string `json:"password" binding:"required"`
	VerificationCode string `json:"code" binding:"required"`
}

type VerificationParameters struct {
	PhoneOrEmail string `binding:"required" json:"phoneOrEmail"`
	GraphicCode  string `json:"graphicCode"`
	CodeType     string `binding:"required" json:"codeType"`
}

type ResetPasswordParameters struct {
	Username         string `json:"username"  binding:"required"`
	NewPassword      string `json:"newPassword"  binding:"required"`
	VerificationCode string `json:"code"  binding:"required"`
}

type ChangePasswordParameters struct {
	Username    string `json:"username"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}
