package model

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        int32     `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"column:add_time"`
	UpdatedAt time.Time `gorm:"column:update_time"`
	DeletedAt gorm.DeletedAt
	IsDeleted bool
}

const (
	LEAVING_MESSAGES = iota + 1
	COMPLAINT
	INQUIRY
	POST_SALE
	WANT_TO_BUY
)

type LeavingMessages struct {
	BaseModel

	User        int32  `gorm:"type:int;index"`
	MessageType int32  `gorm:"type:int comment '留言类型: 1(留言),2(投诉),3(询问),4(售后),5(求购)'"`
	Subject     string `gorm:"type:varchar(100)"`

	Message string // 默认就是text
	File    string `gorm:"type:varchar(200)"`
}

func (LeavingMessages) TableName() string {
	return "leavingmessages"
}

type Address struct {
	BaseModel

	User         int32  `gorm:"type:int;index"`
	Province     string `gorm:"type:varchar(10)"`
	City         string `gorm:"type:varchar(10)"`
	District     string `gorm:"type:varchar(20)"`
	Address      string `gorm:"type:varchar(100)"`
	SignerName   string `gorm:"type:varchar(20)"`
	SignerMobile string `gorm:"type:varchar(11)"`
}

type UserFav struct {
	// 用户对商品的收藏，user与good联合唯一索引
	BaseModel

	User  int32 `gorm:"type:int;index:idx_user_goods,unique"`
	Goods int32 `gorm:"type:int;index:idx_user_goods,unique"`
}

func (UserFav) TableName() string {
	return "userfav"
}
