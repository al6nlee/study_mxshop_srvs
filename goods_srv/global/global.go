package global

import (
	"github.com/elastic/go-elasticsearch/v8"
	"gorm.io/gorm"
	"study_mxshop_srvs/goods_srv/config"
)

var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig  config.NacosConfig
	EsClient     *elasticsearch.Client
)
