package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"strconv"
	"study_mxshop_srvs/goods_srv/global"
	"study_mxshop_srvs/goods_srv/initialize"
	"study_mxshop_srvs/goods_srv/model"
	"time"
)

func syncToEs(goods model.EsGoods) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 将商品数据序列化为 JSON
	body, err := json.Marshal(goods)
	if err != nil {
		return err
	}

	// 执行请求
	res, err := global.EsClient.Index(
		model.EsGoods{}.GetIndexName(),
		bytes.NewReader(body),
		global.EsClient.Index.WithContext(ctx),
		global.EsClient.Index.WithDocumentID(strconv.Itoa(int(goods.ID))),
		global.EsClient.Index.WithRefresh("wait_for"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// 检查响应状态
	if res.IsError() {
		return fmt.Errorf("error indexing document ID %d: %s", goods.ID, res.String())
	}
	return nil
}

// 将mysql中的数据同步至es中
func main() {
	initialize.InitConfig()
	initialize.InitDB()
	initialize.InitEs()

	var goods []model.Goods
	global.DB.Find(&goods)
	for _, g := range goods {
		esModel := model.EsGoods{
			ID:          g.ID,
			CategoryID:  g.CategoryID,
			BrandsID:    g.BrandsID,
			OnSale:      g.OnSale,
			ShipFree:    g.ShipFree,
			IsNew:       g.IsNew,
			IsHot:       g.IsHot,
			Name:        g.Name,
			ClickNum:    g.ClickNum,
			SoldNum:     g.SoldNum,
			FavNum:      g.FavNum,
			MarketPrice: g.MarketPrice,
			GoodsBrief:  g.GoodsBrief,
			ShopPrice:   g.ShopPrice,
		}

		if err := syncToEs(esModel); err != nil {
			log.Printf("从MySQL导入到ES操作失败，ID: %d, ERR: %v", g.ID, err)
		} else {
			log.Printf("从MySQL导入到ES操作成功，ID %d", g.ID)
		}
		// 注意 一定要将docker启动es的java_ops的内存设置大一些 否则运行过程中会出现 bad request错误
	}
}

func main2() {
	dsn := "root:admin123@tcp(172.26.25.139:30006)/mxshop_goods_srv?charset=utf8mb4&parseTime=True&loc=Local"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // 禁用彩色打印
		},
	)

	// 全局模式
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	_ = db.AutoMigrate(&model.Category{}, model.Brands{}, model.GoodsCategoryBrand{}, model.Banner{}, model.Goods{})
}
