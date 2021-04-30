package handler

import (
	"context"
	"goshop/goods_srv/global"
	"goshop/goods_srv/model"
	"goshop/goods_srv/proto"

	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

//获取轮播图列表
func (g *GoodsServer) BannerList(ctx context.Context, req *empty.Empty) (*proto.BannerListResponse, error) {
	var bannerListResponse proto.BannerListResponse
	var bannerResponses []*proto.BannerResponse
	var banners []model.Banner
	res := global.DB.Find(&banners)
	if res.Error != nil {
		return nil, res.Error
	}
	bannerListResponse.Total = int32(res.RowsAffected)

	for _, banner := range banners {
		zap.S().Info(banner)
		bannerInfo := proto.BannerResponse{
			Id:    banner.ID,
			Index: int32(banner.Index),
			Image: banner.Image,
			Url:   banner.Url,
		}
		bannerResponses = append(bannerResponses, &bannerInfo)
	}
	bannerListResponse.Data = bannerResponses
	return &bannerListResponse, nil
}

//创建新的Banner
func (g *GoodsServer) CreateBanner(ctx context.Context, req *proto.BannerRequest) (*proto.BannerResponse, error) {
	var banner model.Banner
	var BannerResponse proto.BannerResponse
	banner.Image = req.Image
	banner.Index = int(req.Index)
	banner.Url = req.Url
	res := global.DB.Save(&banner)
	if res.Error != nil {
		return nil, res.Error
	}
	BannerResponse.Id = banner.ID
	return &BannerResponse, nil
}

//删除Banner
func (g *GoodsServer) DeleteBanner(ctx context.Context, req *proto.BannerRequest) (*empty.Empty, error) {
	if result := global.DB.Delete(&model.Banner{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "banner未找到")
	}
	return &empty.Empty{}, nil
}

//更新banner
func (g *GoodsServer) UpdateBanner(ctx context.Context, req *proto.BannerRequest) (*empty.Empty, error) {
	var bannerModel model.Banner
	if res := global.DB.Where("id=?", req.Id).First(&model.Banner{}); res.RowsAffected == 0 {
		return &emptypb.Empty{}, status.Errorf(codes.NotFound, "banner未找到")
	}
	if req.Image != "" {
		bannerModel.Image = req.Image
	}
	if req.Index != 0 {
		bannerModel.Index = int(req.Index)
	}
	if req.Url != "" {
		bannerModel.Url = req.Url
	}
	res := global.DB.Model(&bannerModel).Where("id=?", req.Id).Updates(bannerModel)
	zap.S().Info(res.RowsAffected)
	return &emptypb.Empty{}, nil
}
