package api

import (
	"goshop_api/goods_web/global"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//将grpc错误码转换成http状态码
func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			zap.S().Error(err.Error())
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

func RemoveTopStruct(fileds map[string]string) map[string]string {
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
		"error": RemoveTopStruct(erros.Translate(global.Trans)),
	})
}

func GetIntParam(key string, ctx *gin.Context) (int32, error) {
	keyStr := ctx.Param(key)
	if keyint, err := strconv.ParseInt(keyStr, 10, 64); err != nil {
		return 0, err
	} else {
		return int32(keyint), nil
	}

}
