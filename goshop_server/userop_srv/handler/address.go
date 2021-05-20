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

//获取用户地址列表
func (*UseropServer) GetAddressList(ctx context.Context, req *proto.AddressRequest) (*proto.AddressListResponse, error) {
	addressModel := []model.Address{}
	addresListRsp := proto.AddressListResponse{}
	res := global.DB.Where(&model.Address{User: req.UserId}).Find(&addressModel)
	if res.RowsAffected != 0 {
		addresListRsp.Total = int32(res.RowsAffected)
	}
	var addressResponse = make([]*proto.AddressResponse, 0)
	for _, addresInfo := range addressModel {
		addressResponse = append(addressResponse, &proto.AddressResponse{
			Id:           addresInfo.ID,
			UserId:       addresInfo.User,
			City:         addresInfo.City,
			Province:     addresInfo.Province,
			Address:      addresInfo.Address,
			District:     addresInfo.District,
			SignerName:   addresInfo.SignerName,
			SignerMobile: addresInfo.SignerMobile,
		})
	}
	addresListRsp.Data = addressResponse
	return &addresListRsp, nil
}

//新建用户地址
func (*UseropServer) CreateAddress(ctx context.Context, req *proto.AddressRequest) (*proto.AddressResponse, error) {
	addressModel := model.Address{
		User:         req.UserId,
		Province:     req.Province,
		City:         req.City,
		District:     req.District,
		Address:      req.Address,
		SignerName:   req.SignerName,
		SignerMobile: req.SignerMobile,
	}
	if res := global.DB.Create(&addressModel); res.RowsAffected == 0 {
		zap.S().Errorf("插入失败:", res.Error)
		return nil, status.Error(codes.Internal, "插入记录出错")
	}
	return &proto.AddressResponse{Id: addressModel.ID}, nil
}

//删除用户地址
func (*UseropServer) DeleteAddress(ctx context.Context, req *proto.AddressRequest) (*empty.Empty, error) {
	// addressModel := model.Address{
	// 	BaseModel: model.BaseModel{ID: req.Id},
	// 	User:      req.UserId,
	// }
	if res := global.DB.Where("id=? AND user=?", req.Id, req.UserId).Delete(&model.Address{}); res.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "收货地址不存在")
	}
	return &emptypb.Empty{}, nil
}

func (*UseropServer) UpdateAddress(ctx context.Context, req *proto.AddressRequest) (*empty.Empty, error) {
	addressModel := model.Address{}
	if req.Address != "" {
		addressModel.Address = req.Address
	}
	if req.City != "" {
		addressModel.City = req.City
	}
	if req.District != "" {
		addressModel.District = req.District
	}
	if req.Province != "" {
		addressModel.Province = req.Province
	}
	if req.SignerMobile != "" {
		addressModel.SignerMobile = req.SignerMobile
	}
	if req.SignerName != "" {
		addressModel.SignerName = req.SignerName
	}
	addressModel.ID = req.Id
	if req.UserId != 0 {
		addressModel.User = req.UserId
	}

	res := global.DB.Updates(&addressModel)
	if res.RowsAffected == 0 {
		zap.S().Errorw("更新地址失败:", res.Error)
		return nil, status.Error(codes.Internal, "更新失败")
	}
	return &emptypb.Empty{}, nil
}
