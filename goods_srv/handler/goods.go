package handler

import "study_mxshop_srvs/goods_srv/proto"

type GoodsServer struct {
	proto.UnimplementedGoodsServer
}

// // 商品接口
// func (g *GoodsServer) GoodsList(context.Context, *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
//
// }

// // 现在用户提交订单有多个商品，你得批量查询商品的信息吧
// func (g *GoodsServer) BatchGetGoods(context.Context, *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error){}
// func (g *GoodsServer) CreateGoods(context.Context, *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error){}
// func (g *GoodsServer) DeleteGoods(context.Context, *proto.DeleteGoodsInfo) (*emptypb.Empty, error){}
// func (g *GoodsServer) UpdateGoods(context.Context, *proto.CreateGoodsInfo) (*emptypb.Empty, error){}
// func (g *GoodsServer) GetGoodsDetail(context.Context, *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error){}

// mustEmbedUnimplementedGoodsServer()
