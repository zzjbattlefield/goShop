package shop_cart

import (
	"context"
	"goshop_api/order_web/api"
	"goshop_api/order_web/forms"
	"goshop_api/order_web/global"
	"goshop_api/order_web/proto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//获取购物车列表
func List(ctx *gin.Context) {
	userId, ok := ctx.Get("userId")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "缺少用户id",
		})
		return
	}
	rsp, err := global.OrderClient.CartItemList(context.Background(), &proto.UserInfo{Id: int32(userId.(uint))})
	if err != nil {
		zap.S().Errorw("[List] 查询购物车列表失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ids := make([]int32, 0)
	for _, cartInfo := range rsp.Data {
		ids = append(ids, cartInfo.GoodsId)
	}
	if len(ids) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"total": 0,
		})
		return
	}
	//请求商品服务获取商品信息
	goodsResp, err := global.GoodsClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{
		Id: ids,
	})
	if err != nil {
		zap.S().Errorw("[List] 查询购物车商品详情失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	reMap := gin.H{
		"total": rsp.Total,
	}
	goodList := make([]interface{}, 0)
	for _, item := range rsp.Data {
		for _, good := range goodsResp.Data {
			if good.Id == item.GoodsId {
				tmpMap := map[string]interface{}{
					"id":         item.Id,
					"goods_id":   item.GoodsId,
					"good_name":  good.Name,
					"good_price": good.ShopPrice,
					"good_image": good.Images,
					"nums":       item.Num,
					"check":      item.Checked,
				}
				goodList = append(goodList, tmpMap)
			}
		}
	}
	reMap["data"] = goodList
	ctx.JSON(http.StatusOK, reMap)
}

//删除购物车商品
func Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	userId, _ := ctx.Get("userId")
	_, err = global.OrderClient.DeleteCartItem(context.Background(), &proto.CartItemRequest{UserId: int32(userId.(uint)), GoodsId: int32(id)})
	if err != nil {
		zap.S().Errorw("[Delete] 删除购物车商品失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.Status(http.StatusOK)
}

//添加商品到购物车
func New(ctx *gin.Context) {
	itemForm := forms.ShopCartItemForm{}
	if err := ctx.ShouldBindJSON(&itemForm); err != nil {
		api.HandleValidatorError(err, ctx)
		return
	}
	//添加前先检查商品是否存在
	_, err := global.GoodsClient.GetGoodsDetail(context.Background(), &proto.GoodInfoRequest{Id: itemForm.GoodsId})
	if err != nil {
		zap.S().Errorw("[New] 查询购物车商品详情失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	//库存是否充足
	invInfo, err := global.InventoryClient.InvDetail(context.Background(), &proto.GoodsInvInfo{GoodsId: itemForm.GoodsId})
	if err != nil {
		zap.S().Errorw("[New] 查询购物车商品库存信息失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	if itemForm.Nums >= invInfo.Num {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"nums": "库存不足",
		})
		return
	}
	userId, _ := ctx.Get("userId")
	rsp, err := global.OrderClient.CreateCartItem(context.Background(), &proto.CartItemRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: itemForm.GoodsId,
		Num:     itemForm.Nums,
	})
	if err != nil {
		zap.S().Errorw("[New] 商品添加购物车失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	reMap := map[string]interface{}{
		"id": rsp.Id,
	}
	ctx.JSON(http.StatusOK, reMap)
}

//更新购物车商品(数量和选中状态)
func Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	itemForm := forms.ShopCartItemUpdateForm{}
	if err := ctx.ShouldBindJSON(&itemForm); err != nil {
		api.HandleValidatorError(err, ctx)
		return
	}
	userId, _ := ctx.Get("userId")
	cartItemReq := proto.CartItemRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: int32(id),
		Num:     itemForm.Num,
	}
	if itemForm.Check != nil {
		cartItemReq.Checked = *itemForm.Check
	}
	if _, err := global.OrderClient.UpdateCartItem(context.Background(), &cartItemReq); err != nil {
		zap.S().Errorw("[Update] 更新购物车商品失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.Status(http.StatusOK)

}
