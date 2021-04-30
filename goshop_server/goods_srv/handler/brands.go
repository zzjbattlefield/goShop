package handler

import (
	"context"
	"goshop/goods_srv/global"
	"goshop/goods_srv/model"
	"goshop/goods_srv/proto"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

//获取品牌列表
func (g *GoodsServer) BrandList(ctx context.Context, req *proto.BrandFilterRequest) (*proto.BrandListResponse, error) {
	brandListRes := proto.BrandListResponse{}
	//proto品牌列表Response
	var brandResponses []*proto.BrandInfoResponse
	//数据库品牌表的结构体
	var brands []model.Brands
	//总数
	var total int64
	//分页查询
	result := global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&brands)
	if result.Error != nil {
		return nil, result.Error
	}
	global.DB.Model(&model.Brands{}).Count(&total)
	brandListRes.Total = int32(total)
	for _, brand := range brands {
		brandResponse := proto.BrandInfoResponse{
			Id:   brand.ID,
			Name: brand.Name,
			Logo: brand.Logo,
		}
		brandResponses = append(brandResponses, &brandResponse)
	}
	brandListRes.Data = brandResponses
	return &brandListRes, nil
}

//新建品牌
func (g *GoodsServer) CreateBrand(ctx context.Context, req *proto.BrandRequest) (*proto.BrandInfoResponse, error) {
	//检查品牌名称是否唯一
	brand := model.Brands{
		Name: req.Name,
		Logo: req.Logo,
	}
	if result := global.DB.Where("name=?", req.Name).First(&brand); result.RowsAffected != 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌已存在")
	}

	global.DB.Save(&brand)
	return &proto.BrandInfoResponse{Id: brand.ID}, nil
}

//删除品牌
func (g *GoodsServer) DeleteBrand(ctx context.Context, req *proto.BrandRequest) (*empty.Empty, error) {
	if result := global.DB.Delete(&model.Brands{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌不存在")
	}
	return &empty.Empty{}, nil
}

//更新品牌
func (g *GoodsServer) UpdateBrand(ctx context.Context, req *proto.BrandRequest) (*empty.Empty, error) {
	var brand model.Brands
	if result := global.DB.Where("id = ?", req.Id).First(&brand); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌不存在")
	}
	if req.Name != "" {
		brand.Name = req.Name
	}
	if req.Logo != "" {
		brand.Logo = req.Logo
	}
	global.DB.Save(&brand)
	return &empty.Empty{}, nil
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
