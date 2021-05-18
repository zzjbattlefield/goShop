package handler

import (
	"context"
	"goshop/userop_srv/global"
	"goshop/userop_srv/model"
	"goshop/userop_srv/proto"

	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

//获取用户收藏列表
func (*UseropServer) GetFavList(ctx context.Context, req *proto.UserFavRequest) (*proto.UserFavListResponse, error) {
	favModel := make([]model.UserFav, 0)
	userFavRsp := proto.UserFavListResponse{}
	res := global.DB.Where(&model.UserFav{Goods: req.GoodsId, User: req.UserId}).Find(&favModel)
	if res.RowsAffected != 0 {
		userFavRsp.Total = int32(res.RowsAffected)
	}
	userFavList := make([]*proto.UserFavResponse, 0)
	for _, favInfo := range favModel {
		userFavList = append(userFavList, &proto.UserFavResponse{
			UserId:  favInfo.User,
			GoodsId: favInfo.Goods,
		})
	}
	userFavRsp.Data = userFavList
	return &userFavRsp, nil
}

//添加用户收藏
func (*UseropServer) AddUserFav(ctx context.Context, req *proto.UserFavRequest) (*empty.Empty, error) {
	favModel := model.UserFav{User: req.UserId, Goods: req.GoodsId}
	res := global.DB.Where(&model.UserFav{User: req.UserId, Goods: req.GoodsId}).First(&favModel)
	if res.RowsAffected != 0 {
		zap.S().Errorw("添加用户收藏:商品已存在")
		return nil, status.Error(codes.AlreadyExists, "商品已存在")
	}
	res = global.DB.Create(&favModel)
	if res.RowsAffected == 0 {
		zap.S().Errorw("添加用户收藏:添加失败", res.Error)
		return nil, status.Error(codes.Internal, "添加商品失败")
	}
	return &emptypb.Empty{}, nil
}

//删除用户收藏 因为有唯一索引所以采用硬删除
func (*UseropServer) DeleteUserFav(ctx context.Context, req *proto.UserFavRequest) (*empty.Empty, error) {
	if res := global.DB.Unscoped().Where(&model.UserFav{User: req.UserId, Goods: req.GoodsId}).Delete(&model.UserFav{}); res.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "记录不存在")
	}
	return &emptypb.Empty{}, nil
}

//查看用户是否拥有此收藏
func (*UseropServer) GetUserFavDetail(ctx context.Context, req *proto.UserFavRequest) (*empty.Empty, error) {
	var userfav model.UserFav
	if result := global.DB.Where("goods=? and user=?", req.GoodsId, req.UserId).Find(&userfav); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "收藏记录不存在")
	}
	return &emptypb.Empty{}, nil
}
