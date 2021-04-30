package banner

import (
	"context"
	"goshop_api/goods_web/api"
	"goshop_api/goods_web/forms"
	"goshop_api/goods_web/global"
	"goshop_api/goods_web/proto"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

func List(ctx *gin.Context) {
	res, err := global.GoodsClinet.BannerList(context.Background(), &emptypb.Empty{})
	if err != nil {
		zap.S().Info("[banner]-[List]:", err.Error())
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	bannerInfo := make([]map[string]interface{}, 0)
	for _, bannerRes := range res.Data {
		bannerInfo = append(bannerInfo, map[string]interface{}{
			"url":   bannerRes.Url,
			"id":    bannerRes.Id,
			"image": bannerRes.Image,
			"index": bannerRes.Index,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data":  bannerInfo,
		"total": res.Total,
	})

}

func Delete(ctx *gin.Context) {
	id, err := api.GetIntParam("id", ctx)
	if err != nil {
		zap.S().Info(err.Error())
		ctx.Status(http.StatusBadRequest)
		return
	}
	if _, err := global.GoodsClinet.DeleteBanner(context.Background(), &proto.BannerRequest{Id: id}); err != nil {
		zap.S().Info("[banner]-[Delete]", err.Error())
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, "删除成功")
}

func New(ctx *gin.Context) {
	bannerForm := forms.BannerForm{}
	if err := ctx.ShouldBindJSON(&bannerForm); err != nil {
		api.HandleValidatorError(err, ctx)
		return
	}
	if res, err := global.GoodsClinet.CreateBanner(context.Background(), &proto.BannerRequest{
		Index: int32(bannerForm.Index),
		Image: bannerForm.Image,
		Url:   bannerForm.Url,
	}); err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	} else {
		ctx.JSON(http.StatusOK, res.Id)
	}
}

func Update(ctx *gin.Context) {
	id, err := api.GetIntParam("id", ctx)
	if err != nil {
		zap.S().Info(err.Error())
		ctx.Status(http.StatusBadRequest)
		return
	}
	bannerForm := forms.BannerForm{}
	if err := ctx.ShouldBindJSON(&bannerForm); err != nil {
		api.HandleValidatorError(err, ctx)
		return
	}
	if _, err := global.GoodsClinet.UpdateBanner(context.Background(), &proto.BannerRequest{
		Id:    id,
		Index: int32(bannerForm.Index),
		Image: bannerForm.Image,
		Url:   bannerForm.Url,
	}); err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	} else {
		ctx.JSON(http.StatusOK, "更新成功")
	}
}
