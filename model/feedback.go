package model

import "server_template/mytime"

type Feedback struct {
	ID        uint                  `json:"id" gorm:"primaryKey"`
	Label     string                `form:"label" gorm:"type:varchar(25)"`
	Feedback  string                `form:"feedback" gorm:"type:varchar(2048)"`
	Picture   string                `form:"-" gorm:"type:varchar(1024)"`
	Phone     string                `form:"phone" gorm:"type:varchar(25)"`
	CreatedAt mytime.CustomTime     `json:"created_at"`
	UpdatedAt mytime.CustomTime     `json:"updated_at"`
	DeletedAt mytime.CustomNullTime `json:"deleted_at" gorm:"index"`
}
