package model

import (
	"database/sql/driver"
	"encoding/json"
)

// Category 分类表
// 类型， 这个字段是否能为null， 这个字段应该设置为可以为null还是设置为空， 0
// 实际开发过程中，尽量设置为不为null
// https://zhuanlan.zhihu.com/p/73997266
// 这些类型我们使用int32还是int，可以减少proto的类型转换，proto没有int类型
type Category struct {
	BaseModel
	Name  string `gorm:"type:varchar(20);not null" json:"name"`
	Level int32  `gorm:"type:int;not null;default:1;comment:'1为1级类目，2为2级...'" json:"level"`
	IsTab bool   `gorm:"default:false;not null;comment:'能否展示在Tab栏'" json:"is_tab"`

	ParentCategoryID int32       `json:"parent"` // 数据库存储的外键ID
	ParentCategory   *Category   `json:"-"`      // gorm中外键自己指向自己，需要使用指针
	SubCategory      []*Category `gorm:"foreignKey:ParentCategoryID;references:ID" json:"sub_category"`
}

type Brands struct {
	BaseModel
	Name string `gorm:"type:varchar(20);not null;comment:'品牌名称'"`
	Logo string `gorm:"type:varchar(200);default:'';not null;comment:'品牌图标'"`
}

// GoodsCategoryBrand 商品和品牌的对应关系
// 如果品牌与分类没有关系的话，那每次选择品牌的时候，会有N多个品牌，使用分类去过滤品牌就可以解决
// 分类与品牌之间的关系是多对多，这里另起一张表
// CategoryID和BrandsID使用相同的index，gorm就会建成联合的索引
type GoodsCategoryBrand struct {
	BaseModel
	CategoryID int32 `gorm:"type:int;index:idx_category_brand,unique"` // 商品外键
	Category   Category

	BrandsID int32 `gorm:"type:int;index:idx_category_brand,unique"` // 品牌外键
	Brands   Brands
}

// TableName 自定义生成的表名
// 为了让gorm生成GoodsCategoryBrand表的时候不生成下划线的方式
// 如goods_category_brand
func (GoodsCategoryBrand) TableName() string {
	return "goodscategorybrand"
}

// Banner 轮播图
type Banner struct {
	BaseModel
	Image string `gorm:"type:varchar(200);not null"`
	Url   string `gorm:"type:varchar(200);not null;comment:'跳转详情页的地址'"`
	Index int32  `gorm:"type:int;default:1;not null;comment:'轮播图的先后顺序'"`
}

// Goods 商品表结构
type Goods struct {
	BaseModel

	CategoryID int32 `gorm:"type:int;not null"`
	Category   Category
	BrandsID   int32 `gorm:"type:int;not null"`
	Brands     Brands

	OnSale   bool `gorm:"default:false;not null;comment:'是否上架'"`
	ShipFree bool `gorm:"default:false;not null;comment:'是否免运费'"`
	IsNew    bool `gorm:"default:false;not null;comment:'是否新品'"`
	IsHot    bool `gorm:"default:false;not null;comment:'是否热卖商品，可用来做广告位'"`

	Name            string   `gorm:"type:varchar(50);not null"`
	GoodsSn         string   `gorm:"type:varchar(50);not null;comment:'商家的内部编号，尺码颜色等'"`
	ClickNum        int32    `gorm:"type:int;default:0;not null;comment:'点击数'"`
	SoldNum         int32    `gorm:"type:int;default:0;not null;comment:'销售量'"`
	FavNum          int32    `gorm:"type:int;default:0;not null;comment:'收藏数'"`
	MarketPrice     float32  `gorm:"not null;comment:'商品价格'"`
	ShopPrice       float32  `gorm:"not null;comment:'实际价格'"`
	GoodsBrief      string   `gorm:"type:varchar(100);not null;comment:'商品简介'"`
	Images          GormList `gorm:"type:varchar(1000);not null;comment:'商品图片'"`
	DescImages      GormList `gorm:"type:varchar(1000);not null;comment:'详情页商品图片'"`
	GoodsFrontImage string   `gorm:"type:varchar(200);not null;comment:'商品展示封面图'"`
}

// GormList 自定义数据类型
type GormList []string

// Value 传入数据时，将切片变为JSON字符串类型，再存入数据库
func (g GormList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (g *GormList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &g)
}
