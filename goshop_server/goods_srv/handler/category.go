package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"goshop/goods_srv/global"
	"goshop/goods_srv/model"
	"goshop/goods_srv/proto"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (g *GoodsServer) GetAllCategorysList(ctx context.Context, req *empty.Empty) (*proto.CategoryListResponse, error) {
	var categorys []model.Category
	global.DB.Where("level=?", 1).Preload("SubCategory.SubCategory").Find(&categorys)
	for _, category := range categorys {
		fmt.Println(category.Name)
	}
	categoryJson, _ := json.Marshal(&categorys)
	return &proto.CategoryListResponse{JsonData: string(categoryJson)}, nil
}

//获取子分类
func (g *GoodsServer) GetSubCategory(ctx context.Context, req *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	var category model.Category
	var subCategorys []model.Category
	var SubCategoryResponse []*proto.CategoryInfoResponse
	if res := global.DB.Where("id=?", req.Id).Find(&category); res.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "分类未找到")
	}
	categoryInfoResponse := proto.CategoryInfoResponse{
		Id:    category.ID,
		Name:  category.Name,
		Level: category.Level,
		IsTab: category.IsTab,
	}
	preloadStr := "SubCategory"
	if categoryInfoResponse.Level == 1 {
		preloadStr = "SubCategory.SubCategory"
	}
	global.DB.Where("parent_category_id=?", categoryInfoResponse.Id).Preload(preloadStr).Find(&subCategorys)
	for _, subCategory := range subCategorys {
		SubCategoryResponse = append(SubCategoryResponse, &proto.CategoryInfoResponse{
			Id:             subCategory.ID,
			Name:           subCategory.Name,
			Level:          subCategory.Level,
			IsTab:          subCategory.IsTab,
			ParentCategory: subCategory.ParentCategoryID,
		})
	}
	return &proto.SubCategoryListResponse{Info: &categoryInfoResponse, SubCategorys: SubCategoryResponse}, nil
}

func (g *GoodsServer) CreateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
	categoryModel := model.Category{
		Name:  req.Name,
		Level: req.Level,
	}
	if categoryModel.Level != 1 {
		categoryModel.ParentCategoryID = req.ParentCategory
	}
	global.DB.Create(&categoryModel)
	return &proto.CategoryInfoResponse{Id: categoryModel.ID}, nil
}

func (s *GoodsServer) DeleteCategory(ctx context.Context, req *proto.DeleteCategoryRequest) (*emptypb.Empty, error) {
	categoryRes := []model.Category{}
	sqlIn := "select id from category where parent_category_id in(select id from category where parent_category_id= ?)"
	subSql := fmt.Sprintf("id = ? or parent_category_id = ? or id in (%s)", sqlIn)
	if result := global.DB.Where(subSql, req.Id, req.Id, req.Id).Find(&categoryRes); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	categoryIdList := make([]int32, 0)
	for _, info := range categoryRes {
		categoryIdList = append(categoryIdList, info.ID)
	}
	if result := global.DB.Where("id in ?", categoryIdList).Delete(&model.Category{}); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.Internal, "删除失败")
	}

	return &emptypb.Empty{}, nil
}

func (s *GoodsServer) UpdateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*emptypb.Empty, error) {
	var category model.Category

	if result := global.DB.First(&category, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.ParentCategory != 0 {
		category.ParentCategoryID = req.ParentCategory
	}
	if req.Level != 0 {
		category.Level = req.Level
	}
	if req.IsTab {
		category.IsTab = req.IsTab
	}

	global.DB.Save(&category)

	return &emptypb.Empty{}, nil
}
