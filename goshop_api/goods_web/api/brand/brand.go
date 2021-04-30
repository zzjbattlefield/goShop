package brand

import (
	"context"
	"goshop_api/goods_web/api"
	"goshop_api/goods_web/global"
	"goshop_api/goods_web/proto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func BrandList(ctx *gin.Context) {
	if rsp, err := global.GoodsClinet.BrandList(context.Background(), &proto.BrandFilterRequest{Pages: 15, PagePerNums: 1}); err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	} else {
		brandInfo := make([]map[string]interface{}, 0)
		for _, brand := range rsp.Data {
			brandInfo = append(brandInfo, map[string]interface{}{
				"id":   brand.Id,
				"name": brand.Name,
				"logo": brand.Logo,
			})
		}
		ctx.JSON(http.StatusOK, gin.H{"data": brandInfo, "total": rsp.Total})
	}

}
