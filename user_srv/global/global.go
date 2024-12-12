package global

import (
	"gorm.io/gorm"
	"study_mxshop_srvs/user_srv/config"
)

var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
)

func init() {

}
