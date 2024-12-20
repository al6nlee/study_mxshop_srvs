package handler

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm/clause"
	"study_mxshop_srvs/inventory_srv/global"
	"study_mxshop_srvs/inventory_srv/model"
	"study_mxshop_srvs/inventory_srv/proto"
	"sync"
)

type InventoryServer struct {
	proto.UnimplementedInventoryServer
}

func (*InventoryServer) SetInv(ctx context.Context, req *proto.GoodsInvInfo) (*emptypb.Empty, error) {
	// 设置库存， 如果我要更新库存
	var inv model.Inventory
	global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv)
	inv.Goods = req.GoodsId
	inv.Stocks = req.Num

	global.DB.Save(&inv)
	return &emptypb.Empty{}, nil
}

func (*InventoryServer) InvDetail(ctx context.Context, req *proto.GoodsInvInfo) (*proto.GoodsInvInfo, error) {
	var inv model.Inventory
	if result := global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "没有库存信息")
	}
	return &proto.GoodsInvInfo{
		GoodsId: inv.Goods,
		Num:     inv.Stocks,
	}, nil
}

func (*InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	// 扣减库存， 本地事务 [1:10,  2:5, 3: 20]
	// 数据库基本的一个应用场景：数据库事务
	// 并发情况之下 可能会出现 超卖 现象 ---> 分布式锁

	tx := global.DB.Begin()
	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory

		mutex := global.Redsync.NewMutex(fmt.Sprintf("goods_%d", goodInfo.GoodsId)) // 锁的粒度尽可能小
		if err := mutex.Lock(); err != nil {
			return nil, status.Errorf(codes.Internal, "获取redis分布式锁异常")
		}

		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() // 回滚之前的操作
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}
		if inv.Stocks < goodInfo.Num {
			tx.Rollback() // 回滚之前的操作
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}
		// 扣减
		inv.Stocks -= goodInfo.Num
		tx.Save(&inv)

		if ok, err := mutex.Unlock(); !ok || err != nil {
			return nil, status.Errorf(codes.Internal, "释放redis分布式锁异常")
		}
	}
	tx.Commit() // 需要自己手动提交操作
	return &emptypb.Empty{}, nil
}

var i = 0

// Sell4 事务、gorm 乐观锁
func (*InventoryServer) Sell4(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	// 扣减库存， 本地事务 [1:10,  2:5, 3: 20]
	// 数据库基本的一个应用场景：数据库事务
	// 并发情况之下 可能会出现 超卖 现象 ---> 分布式锁
	tx := global.DB.Begin()
	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		for {
			if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
				tx.Rollback() // 回滚之前的操作
				return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
			}
			if inv.Stocks < goodInfo.Num {
				tx.Rollback() // 回滚之前的操作
				return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
			}
			// 扣减
			inv.Stocks -= goodInfo.Num
			// 这种写法有瑕疵，为什么？
			// 零值 对于int类型来说 默认值是0 这种会被gorm给忽略掉，Select面对零值的时候，也会强制更新的
			if result := tx.Model(&model.Inventory{}).Select("Stocks", "Version").
				Where("goods = ? and version= ?", goodInfo.GoodsId, inv.Version).
				Updates(model.Inventory{Stocks: inv.Stocks, Version: inv.Version + 1}); result.RowsAffected == 0 {
				i += 1
				zap.S().Infof("库存扣减失败%d", i)
			} else {
				break
			}
		}

	}
	tx.Commit() // 需要自己手动提交操作
	return &emptypb.Empty{}, nil
}

// Sell3 事务 gorm的悲观锁
func (*InventoryServer) Sell3(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	// 扣减库存， 本地事务 [1:10,  2:5, 3: 20]
	// 数据库基本的一个应用场景：数据库事务
	// 并发情况之下 可能会出现 超卖 现象 ---> 分布式锁
	tx := global.DB.Begin()
	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		if result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() // 回滚之前的操作
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}
		if inv.Stocks < goodInfo.Num {
			tx.Rollback() // 回滚之前的操作
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}
		// 扣减
		inv.Stocks -= goodInfo.Num
		tx.Save(&inv)
	}
	tx.Commit() // 需要自己手动提交操作
	return &emptypb.Empty{}, nil
}

var m sync.Mutex

// Sell2 事务、进程锁
func (*InventoryServer) Sell2(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	// 扣减库存， 本地事务 [1:10,  2:5, 3: 20]
	// 数据库基本的一个应用场景：数据库事务
	// 并发情况之下 可能会出现 超卖 现象 ---> 分布式锁
	tx := global.DB.Begin()
	m.Lock() // 获取锁
	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() // 回滚之前的操作
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}
		if inv.Stocks < goodInfo.Num {
			tx.Rollback() // 回滚之前的操作
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}
		// 扣减
		inv.Stocks -= goodInfo.Num
		tx.Save(&inv)
	}
	tx.Commit() // 需要自己手动提交操作
	m.Unlock()  // 释放锁
	return &emptypb.Empty{}, nil
}

// Sell1 引入事务
func (*InventoryServer) Sell1(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	// 扣减库存， 本地事务 [1:10,  2:5, 3: 20]
	// 数据库基本的一个应用场景：数据库事务
	// 并发情况之下 可能会出现 超卖 现象 ---> 分布式锁
	tx := global.DB.Begin()
	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() // 回滚之前的操作
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}
		if inv.Stocks < goodInfo.Num {
			tx.Rollback() // 回滚之前的操作
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}
		// 扣减
		inv.Stocks -= goodInfo.Num
		tx.Save(&inv)
	}
	tx.Commit() // 需要自己手动提交操作
	return &emptypb.Empty{}, nil
}

func (*InventoryServer) Reback(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	// 库存归还：
	// 1：订单超时归还
	// 2. 订单创建失败，归还之前扣减的库存
	// 3. 手动归还
	tx := global.DB.Begin()
	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() // 回滚之前的操作
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}

		// 并发上去后，会出现数据不一致的问题 - 锁，分布式锁
		inv.Stocks += goodInfo.Num
		tx.Save(&inv)
	}
	tx.Commit()
	return &emptypb.Empty{}, nil
}
