package global

import (
	"gorm.io/gorm"
	"study_mxshop_srvs/order_srv/config"
	"study_mxshop_srvs/order_srv/proto"
)

var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig  config.NacosConfig

	GoodsSrvClient     proto.GoodsClient
	InventorySrvClient proto.InventoryClient
)

func init() {

}
