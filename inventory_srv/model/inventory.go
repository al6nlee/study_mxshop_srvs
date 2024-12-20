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

// Stock 仓库
// type Stock struct {
// 	BaseModel
// 	Name    string
// 	Address string
// }

type Inventory struct {
	BaseModel
	Goods   int32 `gorm:"type:int;index"` // 跨服务了，不能使用外键；另外这个字段使用比较多，所以设置一个索引
	Stocks  int32 `gorm:"type:int;comment:'库存量'"`
	Version int32 `gorm:"type:int"`
}
