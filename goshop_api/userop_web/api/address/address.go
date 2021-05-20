package address

import (
	"context"
	"goshop_api/userop_web/api"
	"goshop_api/userop_web/forms"
	"goshop_api/userop_web/global"
	"goshop_api/userop_web/proto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func List(ctx *gin.Context) {
	userId, _ := ctx.Get("userId")
	rsp, err := global.AddressClient.GetAddressList(context.Background(), &proto.AddressRequest{
		UserId: int32(userId.(uint)),
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	addressList := make([]map[string]interface{}, 0)
	for _, addressInfo := range rsp.Data {
		addressList = append(addressList, map[string]interface{}{
			"user_id":     addressInfo.UserId,
			"province":    addressInfo.Province,
			"city":        addressInfo.City,
			"district":    addressInfo.District,
			"address":     addressInfo.Address,
			"sign_name":   addressInfo.SignerName,
			"sign_mobile": addressInfo.SignerMobile,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"total": rsp.Total,
		"data":  addressList,
	})
}

func Delete(ctx *gin.Context) {
	userId, _ := ctx.Get("userId")
	addressIdStr := ctx.Param("id")
	addressid, err := strconv.Atoi(addressIdStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "参数错误")
		return
	}
	_, err = global.AddressClient.DeleteAddress(context.Background(), &proto.AddressRequest{
		UserId: int32(userId.(uint)),
		Id:     int32(addressid),
	})

	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.Status(http.StatusOK)
}

func New(ctx *gin.Context) {
	userId, _ := api.GetUserIdFromCtx(ctx, "userId")
	addressForm := forms.AddressForm{}
	if err := ctx.BindJSON(&addressForm); err != nil {
		api.HandleValidatorError(err, ctx)
		return
	}
	global.AddressClient.CreateAddress(context.Background(), &proto.AddressRequest{
		UserId:       userId,
		Province:     addressForm.Province,
		City:         addressForm.City,
		District:     addressForm.District,
		Address:      addressForm.Address,
		SignerName:   addressForm.SignerName,
		SignerMobile: addressForm.SignerMobile,
	})
}

func Update(ctx *gin.Context) {
	updateForm := forms.AddressForm{}
	if err := ctx.ShouldBindJSON(&updateForm); err != nil {
		api.HandleValidatorError(err, ctx)
		return
	}
	idStr := ctx.Param("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"id": "参数错误",
		})
		return
	}

	userId, _ := api.GetUserIdFromCtx(ctx, "userId")
	_, err = global.AddressClient.UpdateAddress(context.Background(), &proto.AddressRequest{
		Id:           int32(idInt),
		UserId:       userId,
		Province:     updateForm.Province,
		City:         updateForm.City,
		District:     updateForm.District,
		Address:      updateForm.Address,
		SignerName:   updateForm.SignerName,
		SignerMobile: updateForm.SignerMobile,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.Status(http.StatusOK)
}
