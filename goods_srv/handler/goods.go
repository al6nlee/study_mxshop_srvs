package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"study_mxshop_srvs/goods_srv/global"
	"study_mxshop_srvs/goods_srv/model"
	"study_mxshop_srvs/goods_srv/proto"
)

type GoodsServer struct {
	proto.UnimplementedGoodsServer
}

// 商品接口，使用es实现查询
func buildQuery(req *proto.GoodsFilterRequest, categoryIDs []int) (map[string]interface{}, error) {
	query := map[string]interface{}{
		"bool": map[string]interface{}{
			"must":   []interface{}{},
			"filter": []interface{}{},
		},
	}

	must := query["bool"].(map[string]interface{})["must"].([]interface{})
	filter := query["bool"].(map[string]interface{})["filter"].([]interface{})

	if req.KeyWords != "" {
		query["bool"].(map[string]interface{})["must"] = append(must, map[string]interface{}{
			"match": map[string]interface{}{
				"name": req.KeyWords,
			},
		})
	}
	if req.IsHot {
		query["bool"].(map[string]interface{})["filter"] = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"is_hot": true,
			},
		})
	}
	if req.IsNew {
		query["bool"].(map[string]interface{})["filter"] = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"is_new": true,
			},
		})
	}
	if req.PriceMin > 0 {
		query["bool"].(map[string]interface{})["filter"] = append(filter, map[string]interface{}{
			"range": map[string]interface{}{
				"shop_price": map[string]interface{}{
					"gte": req.PriceMin,
				},
			},
		})
	}
	if req.PriceMax > 0 {
		query["bool"].(map[string]interface{})["filter"] = append(filter, map[string]interface{}{
			"range": map[string]interface{}{
				"shop_price": map[string]interface{}{
					"lte": req.PriceMax,
				},
			},
		})
	}
	if req.Brand > 0 {
		query["bool"].(map[string]interface{})["filter"] = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"brands_id": req.Brand,
			},
		})
	}
	if len(categoryIDs) > 0 {
		query["bool"].(map[string]interface{})["filter"] = append(filter, map[string]interface{}{
			"terms": map[string]interface{}{
				"category_id": categoryIDs,
			},
		})
	}
	return query, nil
}

func getCategoryIDs(topCategory int32) ([]int, error) {
	var categoryIDs []int
	var category model.Category
	if result := global.DB.First(&category, topCategory); result.RowsAffected == 0 {
		return nil, fmt.Errorf("分类不存在")
	}

	switch category.Level {
	case 1:
		global.DB.Model(&model.Category{}).Where("parent_category_id IN (?)",
			global.DB.Model(&model.Category{}).Select("id").Where("parent_category_id = ?", topCategory),
		).Pluck("id", &categoryIDs)
	case 2:
		global.DB.Model(&model.Category{}).Where("parent_category_id = ?", topCategory).Pluck("id", &categoryIDs)
	case 3:
		categoryIDs = append(categoryIDs, int(topCategory))
	}
	return categoryIDs, nil
}

func (s *GoodsServer) GoodsList(ctx context.Context, req *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	// 关键词搜索、查询新品、查询热门商品、通过价格区间筛选， 通过商品分类筛选
	goodsListResponse := &proto.GoodsListResponse{}
	esClient := global.EsClient
	// 分类处理
	var categoryIDs []int
	if req.TopCategory > 0 {
		var err error
		categoryIDs, err = getCategoryIDs(req.TopCategory)
		if err != nil {
			return nil, status.Errorf(codes.NotFound, "商品分类不存在")
		}
	}

	// 构建查询
	query, err := buildQuery(req, categoryIDs)
	if err != nil {
		return nil, err
	}

	// 分页设置
	from := (req.Pages - 1) * req.PagePerNums
	if req.Pages <= 0 {
		from = 0
	}
	if req.PagePerNums <= 0 || req.PagePerNums > 100 {
		req.PagePerNums = 10
	}

	searchQuery := map[string]interface{}{
		"from":  from,
		"size":  req.PagePerNums,
		"query": query,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		return nil, status.Errorf(codes.Internal, "构造查询失败")
	}

	// 执行查询
	res, err := esClient.Search(
		esClient.Search.WithContext(ctx),
		esClient.Search.WithIndex("goods"), // 假设索引名为 goods
		esClient.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "搜索失败")
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, status.Errorf(codes.Internal, "响应错误: %s", res.String())
	}

	// 解析响应
	var esResponse struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source model.Goods `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&esResponse); err != nil {
		return nil, status.Errorf(codes.Internal, "解析响应失败")
	}

	// 填充结果
	goodsListResponse.Total = int32(esResponse.Hits.Total.Value)
	for _, hit := range esResponse.Hits.Hits {
		goodsInfoResponse := ModelToResponse(hit.Source)
		goodsListResponse.Data = append(goodsListResponse.Data, &goodsInfoResponse)
	}

	return goodsListResponse, nil
}

// 商品接口
func (s *GoodsServer) GoodsList1(ctx context.Context, req *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	// 关键词搜索、查询新品、查询热门商品、通过价格区间筛选， 通过商品分类筛选
	goodsListResponse := &proto.GoodsListResponse{}

	// match bool 复合查询
	var goods []model.Goods
	localDB := global.DB.Model(model.Goods{})
	if req.KeyWords != "" {
		localDB = localDB.Where("name LIKE ?", "%"+req.KeyWords+"%")
	}
	if req.IsHot {
		localDB = localDB.Where(model.Goods{IsHot: true})
	}
	if req.IsNew {
		localDB = localDB.Where(model.Goods{IsNew: true})
	}
	if req.PriceMin > 0 {
		localDB = localDB.Where("shop_price >= ?", req.PriceMin)
	}
	if req.PriceMax > 0 {
		localDB = localDB.Where("shop_price <= ?", req.PriceMax)
	}
	if req.Brand > 0 {
		localDB = localDB.Where(model.Goods{BrandsID: req.Brand})
	}

	// 通过 category 去查询商品
	var subQuery string
	// categoryIds := make([]interface{}, 0)
	if req.TopCategory > 0 {
		var category model.Category
		if result := global.DB.First(&category, req.TopCategory); result.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "商品分类不存在")
		}

		if category.Level == 1 {
			subQuery = fmt.Sprintf("select id from category where parent_category_id in (select id from category WHERE parent_category_id=%d)", req.TopCategory)
		} else if category.Level == 2 {
			subQuery = fmt.Sprintf("select id from category WHERE parent_category_id=%d", req.TopCategory)
		} else if category.Level == 3 {
			subQuery = fmt.Sprintf("select id from category WHERE id=%d", req.TopCategory) // 白写
		}
		localDB = localDB.Where(fmt.Sprintf("category_id in (%s)", subQuery))
	}

	var total int64
	localDB.Count(&total)
	goodsListResponse.Total = int32(total)

	if req.Pages == 0 {
		req.Pages = 1
	}
	switch {
	case req.PagePerNums > 100:
		req.PagePerNums = 100
	case req.PagePerNums <= 0:
		req.PagePerNums = 10
	}

	// 查询数据库，获取数据
	result := localDB.Preload("Category").Preload("Brands").Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&goods)
	if result.Error != nil {
		return nil, result.Error
	}

	// 组装成 GoodsInfoResponse
	for _, good := range goods {
		goodsInfoResponse := ModelToResponse(good)
		goodsListResponse.Data = append(goodsListResponse.Data, &goodsInfoResponse)
	}

	return goodsListResponse, nil
}

func ModelToResponse(goods model.Goods) proto.GoodsInfoResponse {
	return proto.GoodsInfoResponse{
		Id:              goods.ID,
		CategoryId:      goods.CategoryID,
		Name:            goods.Name,
		GoodsSn:         goods.GoodsSn,
		ClickNum:        goods.ClickNum,
		SoldNum:         goods.SoldNum,
		FavNum:          goods.FavNum,
		MarketPrice:     goods.MarketPrice,
		ShopPrice:       goods.ShopPrice,
		GoodsBrief:      goods.GoodsBrief,
		ShipFree:        goods.ShipFree,
		GoodsFrontImage: goods.GoodsFrontImage,
		IsNew:           goods.IsNew,
		IsHot:           goods.IsHot,
		OnSale:          goods.OnSale,
		DescImages:      goods.DescImages,
		Images:          goods.Images,
		Category: &proto.CategoryBriefInfoResponse{
			Id:   goods.Category.ID,
			Name: goods.Category.Name,
		},
		Brand: &proto.BrandInfoResponse{
			Id:   goods.Brands.ID,
			Name: goods.Brands.Name,
			Logo: goods.Brands.Logo,
		},
	}
}

func (s *GoodsServer) BatchGetGoods(ctx context.Context, req *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
	goodsListResponse := &proto.GoodsListResponse{}

	var goods []model.Goods
	// 调用where并不会真正执行sql 只是用来生成sql的  当调用find，才会真正的执行sql
	result := global.DB.Where(req.Id).Find(&goods)
	if result.Error != nil {
		return nil, result.Error
	}

	for _, good := range goods {
		goodsInfoResponse := ModelToResponse(good)
		goodsListResponse.Data = append(goodsListResponse.Data, &goodsInfoResponse)
	}

	goodsListResponse.Total = int32(result.RowsAffected)
	return goodsListResponse, nil
}

func (s *GoodsServer) GetGoodsDetail(ctx context.Context, req *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
	var goods model.Goods

	if result := global.DB.Preload("Category").Preload("Brands").First(&goods, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	goodsInfoResponse := ModelToResponse(goods)
	return &goodsInfoResponse, nil
}

func (s *GoodsServer) DeleteGoods(ctx context.Context, req *proto.DeleteGoodsInfo) (*emptypb.Empty, error) {
	if result := global.DB.Delete(&model.Goods{BaseModel: model.BaseModel{ID: req.Id}}, req.Id); result.Error != nil {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	return &emptypb.Empty{}, nil
}

func (s *GoodsServer) CreateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var brand model.Brands
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	// 这里没有看到图片文件是如何上传， 在微服务中 普通的文件上传已经不再使用
	goods := model.Goods{
		Brands:          brand,
		BrandsID:        brand.ID,
		Category:        category,
		CategoryID:      category.ID,
		Name:            req.Name,
		GoodsSn:         req.GoodsSn,
		MarketPrice:     req.MarketPrice,
		ShopPrice:       req.ShopPrice,
		GoodsBrief:      req.GoodsBrief,
		ShipFree:        req.ShipFree,
		Images:          req.Images,
		DescImages:      req.DescImages,
		GoodsFrontImage: req.GoodsFrontImage,
		IsNew:           req.IsNew,
		IsHot:           req.IsHot,
		OnSale:          req.OnSale,
	}

	tx := global.DB.Begin()
	result := tx.Save(&goods)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()

	return &proto.GoodsInfoResponse{
		Id: goods.ID,
	}, nil
}

func (s *GoodsServer) UpdateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*emptypb.Empty, error) {
	var goods model.Goods
	if result := global.DB.First(&goods, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}

	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var brand model.Brands
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	goods.Brands = brand
	goods.BrandsID = brand.ID
	goods.Category = category
	goods.CategoryID = category.ID

	goods.Name = req.Name
	goods.GoodsSn = req.GoodsSn
	goods.MarketPrice = req.MarketPrice
	goods.ShopPrice = req.ShopPrice
	goods.GoodsBrief = req.GoodsBrief
	goods.ShipFree = req.ShipFree
	goods.Images = req.Images
	goods.DescImages = req.DescImages
	goods.GoodsFrontImage = req.GoodsFrontImage
	goods.IsNew = req.IsNew
	goods.IsHot = req.IsHot
	goods.OnSale = req.OnSale

	tx := global.DB.Begin()
	result := tx.Save(&goods)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()

	return &emptypb.Empty{}, nil
}
