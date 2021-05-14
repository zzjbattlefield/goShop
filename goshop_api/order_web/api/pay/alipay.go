package pay

import (
	"context"
	"fmt"
	"goshop_api/order_web/global"
	"goshop_api/order_web/proto"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
	"go.uber.org/zap"
)

//实例化支付宝客户端
func InitAliPay() (client *alipay.Client, err error) {
	client, err = alipay.New(global.ServeConfig.AliPayInfo.AppId, global.ServeConfig.AliPayInfo.PrivateKey, false)
	if err != nil {
		zap.S().Errorw("实例化支付宝url失败")
		return nil, err
	}
	err = client.LoadAliPayPublicKey(global.ServeConfig.AliPayInfo.AliPublicKey)
	if err != nil {
		zap.S().Errorw("加载支付宝PublicKey失败")
		return nil, err
	}
	return client, nil
}

//支付宝回调通知
func Notify(ctx *gin.Context) {
	client, err := InitAliPay()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
	}
	noti, err := client.GetTradeNotification(ctx.Request)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	fmt.Println("交易状态:", noti)
	//更新交易状态
	_, err = global.OrderClient.UpdateOrderStatus(context.Background(), &proto.OrderStatus{OrderSn: noti.OutTradeNo, Status: string(noti.TradeStatus)})
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	//确认收到消息
	alipay.AckNotification(ctx.Writer)
}
