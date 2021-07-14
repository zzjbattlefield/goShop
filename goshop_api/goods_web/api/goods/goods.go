package goods

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"goshop_api/goods_web/forms"
	"goshop_api/goods_web/global"
	"goshop_api/goods_web/proto"
)

//将grpc错误码转换成http状态码
func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{"msg": e.Message()})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "内部错误"})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "服务不可用"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "内部错误"})
			}
		}
	}
}

func removeTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

func HandleValidatorError(err error, ctx *gin.Context) {
	erros, ok := err.(validator.ValidationErrors)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	ctx.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopStruct(erros.Translate(global.Trans)),
	})
}

func List(ctx *gin.Context) {
	goodsRequest := proto.GoodsFilterRequest{}
	priceMax := ctx.DefaultQuery("pmax", "0")
	priceMaxInt, _ := strconv.Atoi(priceMax)
	goodsRequest.PriceMax = int32(priceMaxInt)

	priceMin := ctx.DefaultQuery("pmin", "0")
	priceMinInt, _ := strconv.Atoi(priceMin)
	goodsRequest.PriceMin = int32(priceMinInt)

	isHot := ctx.DefaultQuery("ih", "0")
	if isHot == "1" {
		goodsRequest.IsHot = true
	}
	isNew := ctx.DefaultQuery("in", "0")
	if isNew == "1" {
		goodsRequest.IsNew = true
	}
	isTab := ctx.DefaultQuery("it", "0")
	if isTab == "1" {
		goodsRequest.IsTab = true
	}
	categoryId := ctx.DefaultQuery("c", "0")
	categoryIdInt, _ := strconv.Atoi(categoryId)
	goodsRequest.TopCategory = int32(categoryIdInt)

	pages := ctx.DefaultQuery("p", "0")
	pagesInt, _ := strconv.Atoi(pages)
	goodsRequest.Pages = int32(pagesInt)

	perNums := ctx.DefaultQuery("pnum", "0")
	perNumsInt, _ := strconv.Atoi(perNums)
	goodsRequest.PagePerNums = int32(perNumsInt)

	keywords := ctx.DefaultQuery("q", "")
	goodsRequest.KeyWords = keywords

	brandId := ctx.DefaultQuery("b", "0")
	brandIdInt, _ := strconv.Atoi(brandId)
	goodsRequest.Brand = int32(brandIdInt)

	context := context.WithValue(context.Background(), "ginContext", ctx)
	r, err := global.GoodsClinet.GoodsList(context, &goodsRequest)
	if err != nil {
		zap.S().Errorw("[list]查询[商品列表失败]")
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	reMap := make(map[string]interface{})
	reMap["Total"] = r.Total
	goodsList := make([]map[string]interface{}, 0)
	for _, value := range r.Data {
		goodsList = append(goodsList, map[string]interface{}{
			"id":          value.Id,
			"name":        value.Name,
			"goods_brief": value.GoodsBrief,
			"desc":        value.GoodsDesc,
			"ship_free":   value.ShipFree,
			"images":      value.Images,
			"desc_images": value.DescImages,
			"front_image": value.GoodsFrontImage,
			"shop_price":  value.ShopPrice,
			"category": map[string]interface{}{
				"id":   value.Category.Id,
				"name": value.Category.Name,
			},
			"brand": map[string]interface{}{
				"id":   value.Brand.Id,
				"name": value.Brand.Name,
				"logo": value.Brand.Logo,
			},
			"is_hot":  value.IsHot,
			"is_new":  value.IsNew,
			"on_sale": value.OnSale,
		})
	}
	reMap["Data"] = goodsList
	ctx.JSON(http.StatusOK, reMap)
}

func New(ctx *gin.Context) {
	goodsForm := forms.GoodsForm{}
	//参数校验
	if err := ctx.ShouldBindJSON(&goodsForm); err != nil {
		HandleValidatorError(err, ctx)
		return
	}
	goodsClient := global.GoodsClinet
	rsp, err := goodsClient.CreateGoods(context.Background(), &proto.CreateGoodsInfo{
		Name:            goodsForm.Name,
		GoodsSn:         goodsForm.GoodsSn,
		Stocks:          goodsForm.Stocks,
		MarketPrice:     goodsForm.MarketPrice,
		ShopPrice:       goodsForm.ShopPrice,
		GoodsBrief:      goodsForm.GoodsBrief,
		ShipFree:        *goodsForm.ShipFree,
		Images:          goodsForm.Images,
		DescImages:      goodsForm.DescImages,
		GoodsFrontImage: goodsForm.FrontImage,
		CategoryId:      goodsForm.CategoryId,
		BrandId:         goodsForm.Brand,
	})
	if err != nil {
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	//TODO 商品的库存 - 分布式
	ctx.JSON(http.StatusOK, rsp)
}

//获取商品详情
func Detail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	// idStr := ctx.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	res, err := global.GoodsClinet.GetGoodsDetail(context.Background(), &proto.GoodInfoRequest{Id: int32(id)})
	if err != nil {
		HandleValidatorError(err, ctx)
		return
	}
	//TODO 商品库存单独获取
	ctx.JSON(http.StatusOK, res)
}

//商品删除操作
func Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		zap.S().Info(err.Error())
		ctx.Status(http.StatusNotFound)
		return
	}
	_, err = global.GoodsClinet.DeleteGoods(context.Background(), &proto.DeleteGoodsInfo{Id: int32(id)})
	if err != nil {
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.Status(http.StatusOK)
}

func Stocks(ctx *gin.Context) {
	idStr := ctx.Param("id")
	_, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		zap.S().Info(err.Error())
		ctx.Status(http.StatusNotFound)
		return
	}
	//TODO 获取商品库存
	return
}

func UpdateStatus(ctx *gin.Context) {
	statusForm := forms.GoodsStatusForm{}
	if err := ctx.ShouldBindJSON(&statusForm); err != nil {
		HandleValidatorError(err, ctx)
		return
	}
	idStr := ctx.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 32)
	if _, err := global.GoodsClinet.UpdateGoods(context.Background(), &proto.CreateGoodsInfo{
		Id:     int32(id),
		IsHot:  *statusForm.IsHot,
		IsNew:  *statusForm.IsNew,
		OnSale: *statusForm.OnSale,
	}); err != nil {
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "修改成功",
	})
}

func Update(ctx *gin.Context) {
	goodsForm := forms.GoodsForm{}
	if err := ctx.ShouldBindJSON(&goodsForm); err != nil {
		HandleValidatorError(err, ctx)
		return
	}

	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if _, err = global.GoodsClinet.UpdateGoods(context.Background(), &proto.CreateGoodsInfo{
		Id:              int32(i),
		Name:            goodsForm.Name,
		GoodsSn:         goodsForm.GoodsSn,
		Stocks:          goodsForm.Stocks,
		MarketPrice:     goodsForm.MarketPrice,
		ShopPrice:       goodsForm.ShopPrice,
		GoodsBrief:      goodsForm.GoodsBrief,
		ShipFree:        *goodsForm.ShipFree,
		Images:          goodsForm.Images,
		DescImages:      goodsForm.DescImages,
		GoodsFrontImage: goodsForm.FrontImage,
		CategoryId:      goodsForm.CategoryId,
		BrandId:         goodsForm.Brand,
	}); err != nil {
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}
