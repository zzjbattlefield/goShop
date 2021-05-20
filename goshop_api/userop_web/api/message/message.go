package message

import (
	"context"
	"goshop_api/user_web/api"
	"goshop_api/userop_web/forms"
	"goshop_api/userop_web/global"
	"goshop_api/userop_web/models"
	"goshop_api/userop_web/proto"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func List(ctx *gin.Context) {
	var request = proto.MessageRequest{}
	userIdItf, ok := ctx.Get("userId")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "未获得UserId"})
		return
	}
	userId := int32(userIdItf.(uint))
	//获取登录用户身份
	claim, _ := ctx.Get("claims")
	model := claim.(*models.CustomClaims)
	if model.AuthorityId == 1 {
		//普通用户
		request.UserId = userId
	}
	res, err := global.MessageClient.MessageList(context.Background(), &request)
	if err != nil {
		zap.S().Errorw("获取留言列表失败", err.Error())
		api.HandleGrpcErrorToHttp(err, ctx)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"total": res.Total,
		"data":  res.Data,
	})
}

func New(ctx *gin.Context) {
	messageForm := &forms.MessageForm{}
	if err := ctx.ShouldBindJSON(messageForm); err != nil {
		api.HandleValidatorError(err, ctx)
		return
	}
	userIdStr, _ := ctx.Get("userId")
	userId := int32(userIdStr.(uint))
	rsp, err := global.MessageClient.CreateMessage(context.Background(), &proto.MessageRequest{
		UserId:      userId,
		MessageType: messageForm.MessageType,
		Subject:     messageForm.Subject,
		Message:     messageForm.Message,
		File:        messageForm.File,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"id": rsp.Id,
	})
}
