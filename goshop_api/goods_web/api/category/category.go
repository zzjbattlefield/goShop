package category

import (
	"context"
	"encoding/json"
	"goshop_api/goods_web/forms"
	"goshop_api/goods_web/global"
	"goshop_api/goods_web/proto"
	"goshop_api/user_web/api"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

func List(ctx *gin.Context) {
	res, err := global.GoodsClinet.GetAllCategorysList(context.Background(), &emptypb.Empty{})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
	}
	zap.S().Info(res.JsonData)
	categoryInfo := make([]interface{}, 0)
	json.Unmarshal([]byte(res.JsonData), &categoryInfo)
	ctx.JSON(http.StatusOK, res.JsonData)
}

func Detail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		zap.S().Errorf("获取参数错误:%s", err.Error())
		api.HandleValidatorError(err, ctx)
		return
	}
	if res, err := global.GoodsClinet.GetSubCategory(context.Background(), &proto.CategoryListRequest{Id: int32(id)}); err != nil {
		zap.S().Errorf("[商品服务][Detail]获取分类子分类错误:%s", err.Error())
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	} else {
		responseMap := make(map[string]interface{})
		subCategory := make([]interface{}, 0)
		for _, categoryInfo := range res.SubCategorys {
			subCategory = append(subCategory, map[string]interface{}{
				"id":     categoryInfo.Id,
				"is_tab": categoryInfo.IsTab,
				"level":  categoryInfo.Level,
				"name":   categoryInfo.Name,
			})
		}
		responseMap["sub_category"] = subCategory
		responseMap["id"] = res.Info.Id
		responseMap["is_tab"] = res.Info.IsTab
		responseMap["level"] = res.Info.Level
		responseMap["name"] = res.Info.Name
		ctx.JSON(http.StatusOK, responseMap)
	}
}

func New(ctx *gin.Context) {
	cateogryForm := forms.CategoryForm{}
	if err := ctx.ShouldBindJSON(&cateogryForm); err != nil {
		api.HandleValidatorError(err, ctx)
		return
	}
	if res, err := global.GoodsClinet.CreateCategory(context.Background(), &proto.CategoryInfoRequest{
		Name:           cateogryForm.Name,
		Level:          cateogryForm.Level,
		IsTab:          *cateogryForm.IsTab,
		ParentCategory: cateogryForm.ParentCategory,
	}); err != nil {
		zap.S().Errorf("[category]-[New] :%s", err)
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	} else {
		response := make(map[string]interface{})
		response["id"] = res.Id
		response["name"] = res.Name
		response["level"] = res.Level
		response["is_tab"] = res.IsTab
		ctx.JSON(http.StatusOK, response)
	}

}

func Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	updateCategoryForm := forms.UpdateCategoryForm{}
	if err := ctx.ShouldBind(&updateCategoryForm); err != nil {
		api.HandleValidatorError(err, ctx)
		return
	}
	CategoryRequest := &proto.CategoryInfoRequest{
		Id:   int32(id),
		Name: updateCategoryForm.Name,
	}
	if updateCategoryForm.IsTab != nil {
		CategoryRequest.IsTab = *updateCategoryForm.IsTab
	}
	if _, err = global.GoodsClinet.UpdateCategory(context.Background(), CategoryRequest); err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	} else {
		ctx.JSON(http.StatusOK, "更新成功")
	}
}

func Delete(ctx *gin.Context) {
	//删除分类时候也删除其下的下属分类
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	if _, err := global.GoodsClinet.DeleteCategory(context.Background(), &proto.DeleteCategoryRequest{Id: int32(id)}); err != nil {
		zap.S().Errorf("[category]-[Delete]:%s", err)
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, "删除成功")
}
