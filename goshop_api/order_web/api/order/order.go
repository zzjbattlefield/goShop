package order

import (
	"context"
	"goshop_api/order_web/api"
	"goshop_api/order_web/api/pay"
	"goshop_api/order_web/forms"
	"goshop_api/order_web/global"
	"goshop_api/order_web/models"
	"goshop_api/order_web/proto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
	"go.uber.org/zap"
)

func List(ctx *gin.Context) {
	userId, _ := ctx.Get("userId")
	claims, _ := ctx.Get("claims")
	model := claims.(*models.CustomClaims)
	request := proto.OrderFilterRequest{}
	//管理员用户返回所有订单列表
	if model.AuthorityId == 1 {
		//普通用户
		request.UserId = int32(userId.(uint))
	}

	//分页
	pages := ctx.DefaultQuery("p", "0")
	pagesInt, _ := strconv.Atoi(pages)
	request.Pages = int32(pagesInt)

	perNums := ctx.DefaultQuery("pnum", "0")
	perNumInt, _ := strconv.Atoi(perNums)
	request.PagePerNums = int32(perNumInt)

	res, err := global.OrderClient.OrderList(context.Background(), &request)
	if err != nil {
		zap.S().Errorw("[List] 获取订单列表失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	rspMap := make(map[string]interface{})
	rspMap["total"] = res.Total

	orderList := make([]interface{}, 0)
	for _, orderInfo := range res.Data {
		tmpMap := map[string]interface{}{
			"id":       orderInfo.Id,
			"status":   orderInfo.Status,
			"pay_type": orderInfo.PayType,
			"user":     orderInfo.UserId,
			"post":     orderInfo.Post,
			"total":    orderInfo.Total,
			"address":  orderInfo.Address,
			"name":     orderInfo.Name,
			"mobile":   orderInfo.Mobile,
			"order_sn": orderInfo.OrderSn,
			"add_time": orderInfo.AddTime,
		}
		orderList = append(orderList, tmpMap)
	}
	rspMap["data"] = orderList
	ctx.JSON(http.StatusOK, rspMap)
}

func New(ctx *gin.Context) {
	orderForm := forms.CreateOrderForm{}
	if err := ctx.ShouldBindJSON(&orderForm); err != nil {
		api.HandleValidatorError(err, ctx)
		return
	}
	uid, _ := ctx.Get("userId")
	rsp, err := global.OrderClient.CreateOrder(context.Background(), &proto.OrderRequest{
		UserId:  int32(uid.(uint)),
		Address: orderForm.Address,
		Name:    orderForm.Name,
		Mobile:  orderForm.Mobile,
		Post:    orderForm.Post,
	})
	if err != nil {
		zap.S().Errorf("新建订单失败:%v", err.Error())
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	//生成支付宝支付url
	var p = alipay.TradePagePay{}
	p.NotifyURL = global.ServeConfig.AliPayInfo.NotifyURL
	p.ReturnURL = global.ServeConfig.AliPayInfo.ReturnUrl
	p.Subject = rsp.OrderSn
	p.OutTradeNo = rsp.OrderSn
	p.TotalAmount = strconv.FormatFloat(float64(rsp.Total), 'f', 2, 64)
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"
	url, err := createAliPayUrl(p)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"id":          rsp.Id,
		"ali_pay_url": url,
	})
}

func Detail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "url地址出错",
		})
	}
	orderReq := &proto.OrderRequest{Id: int32(idInt)}
	userId, _ := ctx.Get("userId")

	claims, _ := ctx.Get("claims")
	model := claims.(*models.CustomClaims)
	if model.AuthorityId == 1 {
		orderReq.UserId = int32(userId.(uint))
	}
	orderInfoRes, err := global.OrderClient.OrderDetail(context.Background(), orderReq)
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	//生成支付宝支付url
	var p = alipay.TradePagePay{}
	p.NotifyURL = global.ServeConfig.AliPayInfo.NotifyURL
	p.ReturnURL = global.ServeConfig.AliPayInfo.ReturnUrl
	p.Subject = orderInfoRes.OrderInfo.OrderSn
	p.OutTradeNo = orderInfoRes.OrderInfo.OrderSn
	p.TotalAmount = strconv.FormatFloat(float64(orderInfoRes.OrderInfo.Total), 'f', 2, 64)
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"
	url, err := createAliPayUrl(p)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"orderInfo": orderInfoRes.OrderInfo,
		"goodsInfo": orderInfoRes.Goods,
		"url":       url,
	})
}

//创建支付宝支付url
func createAliPayUrl(tradePagePay alipay.TradePagePay) (string, error) {
	client, err := pay.InitAliPay()
	if err != nil {
		return "", err
	}
	url, err := client.TradePagePay(tradePagePay)
	if err != nil {
		zap.S().Errorw("生成支付url失败")
		return "", err
	}
	return url.String(), nil
}
