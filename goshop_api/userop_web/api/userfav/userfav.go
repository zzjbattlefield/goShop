package userfav

import (
	"context"
	"goshop_api/userop_web/api"
	"goshop_api/userop_web/forms"
	"goshop_api/userop_web/global"
	"goshop_api/userop_web/proto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//获取用户收藏列表
func List(ctx *gin.Context) {
	userId, _ := api.GetUserIdFromCtx(ctx, "userId")
	userFavListRsp := make(map[string]interface{})
	res, err := global.UserFavClient.GetFavList(context.Background(), &proto.UserFavRequest{
		UserId: userId,
	})
	if err != nil {
		zap.S().Errorw("查询用户收藏列表失败", err.Error())
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	userFavListRsp["total"] = res.Total
	ids := make([]int32, 0)
	for _, userFavListInfo := range res.Data {
		ids = append(ids, userFavListInfo.GoodsId)
	}
	//查询相关商品信息
	goodsRes, err := global.GoodsClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{
		Id: ids,
	})
	if err != nil {
		zap.S().Errorw("【获取用户收藏列表】查询批量商品信息失败")
		api.HandleGrpcErrorToHttp(err, ctx)
	}
	userFavList := make([]map[string]interface{}, 0)
	for _, userFavListInfo := range res.Data {
		tmpFavInfo := make(map[string]interface{})
		tmpFavInfo["user_id"] = userFavListInfo.UserId
		tmpFavInfo["goods_id"] = userFavListInfo.GoodsId
		for _, goodsListInfo := range goodsRes.Data {
			if goodsListInfo.Id == userFavListInfo.GoodsId {
				tmpFavInfo["name"] = goodsListInfo.Name
				tmpFavInfo["shop_price"] = goodsListInfo.ShopPrice
			}
		}
		userFavList = append(userFavList, tmpFavInfo)
	}
	userFavListRsp["data"] = userFavList

	ctx.JSON(http.StatusOK, userFavListRsp)
}

//用户是否有收藏本商品
func Detail(ctx *gin.Context) {
	userId, _ := api.GetUserIdFromCtx(ctx, "userId")
	goodsIdStr := ctx.Param("id")
	goodsId, err := strconv.Atoi(goodsIdStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "参数错误")
	}
	_, err = global.UserFavClient.GetUserFavDetail(context.Background(), &proto.UserFavRequest{
		UserId:  userId,
		GoodsId: int32(goodsId),
	})
	if err != nil {
		zap.S().Errorw("查询用户收藏商品失败:", err.Error())
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.Status(http.StatusOK)
}

//新建用户收藏
func New(ctx *gin.Context) {
	userId, _ := api.GetUserIdFromCtx(ctx, "userId")
	userFavForm := forms.UserFavForm{}
	if err := ctx.ShouldBindJSON(&userFavForm); err != nil {
		api.HandleValidatorError(err, ctx)
		return
	}
	_, err := global.UserFavClient.AddUserFav(context.Background(), &proto.UserFavRequest{
		UserId:  userId,
		GoodsId: userFavForm.GoodsId,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.Status(http.StatusOK)
}

func Delete(ctx *gin.Context) {
	userId, _ := api.GetUserIdFromCtx(ctx, "userId")
	goodsIdStr := ctx.Param("id")
	goodsId, _ := strconv.Atoi(goodsIdStr)
	_, err := global.UserFavClient.DeleteUserFav(context.Background(), &proto.UserFavRequest{UserId: userId, GoodsId: int32(goodsId)})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
	}
	ctx.Status(http.StatusOK)
}
