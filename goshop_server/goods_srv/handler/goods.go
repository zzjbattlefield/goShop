package handler

import (
	"context"
	"fmt"
	"goshop/goods_srv/global"
	"goshop/goods_srv/model"
	"goshop/goods_srv/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GoodsServer struct {
	proto.UnimplementedGoodsServer
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

//获取商品列表页
func (g *GoodsServer) GoodsList(ctx context.Context, req *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	goodsListResponse := &proto.GoodsListResponse{}
	var goodsModel []model.Goods
	localDB := global.DB.Model(&goodsModel)
	if req.IsHot {
		localDB.Where("is_hot=?", req.IsHot)
	}
	if req.IsNew {
		localDB = localDB.Where(&model.Goods{IsNew: true})
	}
	if req.PriceMin > 0 {
		localDB = localDB.Where("shop_price>=?", req.PriceMin)
	}
	if req.PriceMax > 0 {
		localDB = localDB.Where("shop_price<=?", req.PriceMax)
	}

	if req.Brand > 0 {
		localDB = localDB.Where("brand_id=?", req.Brand)
	}
	if req.KeyWords != "" {
		localDB = localDB.Where(" goods.name like ?", "%"+req.KeyWords+"%")
	}

	if req.TopCategory > 0 {
		//先查询是否存在分类
		var categoryModel model.Category
		if res := global.DB.Where("id=?", req.TopCategory).First(&categoryModel); res.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "分类不存在")
		}
		var subSql string
		if categoryModel.Level == 1 {
			subSql = fmt.Sprintf("SELECT id FROM category WHERE parent_category_id in (SELECT id FROM category WHERE parent_category_id = %d)", req.TopCategory)
		} else if categoryModel.Level == 2 {
			subSql = fmt.Sprintf("SELECT id FROM category WHERE parent_category_id = %d", req.TopCategory)
		} else if categoryModel.Level == 3 {
			subSql = fmt.Sprintf("SELECT id FROM category WHERE id = %d", req.TopCategory)
		}

		localDB.Where(fmt.Sprintf("category_id in (%s)", subSql))
	}

	var count int64
	var goodsInfoResponse []*proto.GoodsInfoResponse
	localDB.Count(&count)
	goodsListResponse.Total = int32(count)
	localDB.Joins("Category").Joins("Brands").Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&goodsModel)
	for _, goodsInfo := range goodsModel {
		info := ModelToResponse(goodsInfo)
		goodsInfoResponse = append(goodsInfoResponse, &info)
	}
	goodsListResponse.Data = goodsInfoResponse
	return goodsListResponse, nil
}

//批量查询商品
func (g *GoodsServer) BatchGetGoods(ctx context.Context, req *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
	var goodsModel []model.Goods
	var goodsListResponse proto.GoodsListResponse
	res := global.DB.Where("id IN (?)", req.Id).Find(&goodsModel)
	goodsListResponse.Total = int32(res.RowsAffected)
	for _, good := range goodsModel {
		goodResponse := ModelToResponse(good)
		goodsListResponse.Data = append(goodsListResponse.Data, &goodResponse)
	}
	return &goodsListResponse, nil
}

func (g *GoodsServer) CreateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var brand model.Brands
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}
	//这里没有看到图片文件是如何上传， 在微服务中 普通的文件上传已经不再使用
	goods := model.Goods{
		// Brands:          brand,
		BrandsID: brand.ID,
		// Category:        category,
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
		Stocks:          req.Stocks,
		IsNew:           req.IsNew,
		IsHot:           req.IsHot,
		OnSale:          req.OnSale,
	}

	//srv之间互相调用了
	global.DB.Save(&goods)
	return &proto.GoodsInfoResponse{
		Id: goods.ID,
	}, nil
}

func (g *GoodsServer) DeleteGoods(ctx context.Context, req *proto.DeleteGoodsInfo) (*emptypb.Empty, error) {
	if result := global.DB.Delete(&model.Goods{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	return &emptypb.Empty{}, nil
}

func (g *GoodsServer) UpdateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*emptypb.Empty, error) {
	var goods model.Goods

	if result := global.DB.First(&goods, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}

	if req.CategoryId != 0 {
		var category model.Category
		if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
			return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
		}
		goods.CategoryID = category.ID
	}

	if req.BrandId != 0 {
		var brand model.Brands
		if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
			return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
		}
		goods.BrandsID = brand.ID
	}

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

	global.DB.Updates(&goods)
	return &emptypb.Empty{}, nil
}

//获的商品详情
func (g *GoodsServer) GetGoodsDetail(ctx context.Context, req *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
	var goodsModel model.Goods
	var goodsInfoResponse proto.GoodsInfoResponse
	res := global.DB.Joins("Brands").Joins("Category").Where("goods.id = ?", req.Id).First(&goodsModel)
	if res.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "商品未找到")
	}
	goodsInfoResponse = ModelToResponse(goodsModel)
	return &goodsInfoResponse, nil
}
