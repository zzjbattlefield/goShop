package api

import (
	"fmt"
	"goshop_api/user_web/forms"
	"goshop_api/user_web/global"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

//生成短信验证码
func GenerateSmsCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())
	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

func SendSms(ctx *gin.Context) {
	SendSmsForm := forms.SendSmsForm{}
	if err := ctx.ShouldBind(&SendSmsForm); err != nil {
		HandleValidatorError(err, ctx)
		return
	}
	client, err := dysmsapi.NewClientWithAccessKey("cn-beijing", global.ServeConfig.SmsInfo.AccessKeyId, global.ServeConfig.SmsInfo.AccessKeySecret)
	if err != nil {
		panic(err)
	}
	smsCode := GenerateSmsCode(4)
	zap.S().Infof("发送的验证码:%v", smsCode)
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = "cn-beijing"
	request.QueryParams["PhoneNumbers"] = SendSmsForm.Mobile                      //手机号
	request.QueryParams["SignName"] = "swoole赛事直播"                                //阿里云验证过的项目名 自己设置
	request.QueryParams["TemplateCode"] = global.ServeConfig.SmsInfo.TemplateCode //阿里云的短信模板号 自己设置
	request.QueryParams["TemplateParam"] = "{\"code\":" + smsCode + "}"           //短信模板中的验证码内容 自己生成   之前试过直接返回，但是失败，加上code成功。
	response, err := client.ProcessCommonRequest(request)
	client.DoAction(request, response)
	if err != nil {
		fmt.Print(err.Error())
	}
	//将验证码保存到redis
	redisDB := redis.NewClient(&redis.Options{Addr: fmt.Sprintf("%s:%d", global.ServeConfig.RedisInfo.Host, global.ServeConfig.RedisInfo.Port)})
	redisDB.Set(SendSmsForm.Mobile, smsCode, time.Duration(global.ServeConfig.RedisInfo.Expire)*time.Second)
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "发送成功",
	})
}
