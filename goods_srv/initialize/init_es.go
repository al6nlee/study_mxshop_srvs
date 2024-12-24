package initialize

import (
	"bytes"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"study_mxshop_srvs/goods_srv/global"
	"study_mxshop_srvs/goods_srv/model"
)

func InitEs() {
	// 初始化连接
	host := fmt.Sprintf("http://%s:%d", global.ServerConfig.EsInfo.Host, global.ServerConfig.EsInfo.Port)
	var err error
	global.EsClient, err = elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{host},
	})
	if err != nil {
		panic(err)
	}

	// 新建mapping和index
	res, err := global.EsClient.Indices.Exists([]string{model.EsGoods{}.GetIndexName()})
	if err != nil {
		log.Fatalf("检查索引异常: %s", err)
	}
	defer res.Body.Close()

	// 如果索引不存在则新建索引
	if res.StatusCode == 404 {
		mapping := model.EsGoods{}.GetMapping()

		createRes, err := global.EsClient.Indices.Create(model.EsGoods{}.GetIndexName(), global.EsClient.Indices.Create.WithBody(bytes.NewReader([]byte(mapping))))
		if err != nil {
			log.Fatalf("创建索引异常: %s", err)
		}
		defer createRes.Body.Close()

		if createRes.IsError() {
			log.Fatalf("创建索引失败: %s", createRes)
		}
		log.Println("创建索引成功")
	}
}
