package initialize

import (
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"study_mxshop_srvs/order_srv/global"
	"study_mxshop_srvs/order_srv/proto"
)

func InitSrvConn() {
	userConn, err := grpc.Dial(
		fmt.Sprintf(
			"consul://%s:%d/%s?wait=14s",
			global.ServerConfig.ConsulInfo.Host,
			global.ServerConfig.ConsulInfo.Port,
			global.ServerConfig.GoodsSrvInfo.Name,
		),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`), // 轮询
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
	}
	// 生成grpc的client调用接口
	global.GoodsSrvClient = proto.NewGoodsClient(userConn)

	// 初始化库存服务连接
	invConn, err := grpc.Dial(
		fmt.Sprintf(
			"consul://%s:%d/%s?wait=14s",
			global.ServerConfig.ConsulInfo.Host,
			global.ServerConfig.ConsulInfo.Port,
			global.ServerConfig.InventorySrvInfo.Name,
		),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【库存服务失败】")
	}
	// 生成grpc的client调用接口
	global.InventorySrvClient = proto.NewInventoryClient(invConn)
}
