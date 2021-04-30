package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"goshop_api/user_web/forms"
	"goshop_api/user_web/global"
	"goshop_api/user_web/global/response"
	"goshop_api/user_web/middlewares"
	"goshop_api/user_web/models"
	proto "goshop_api/user_web/proto"
)

//将grpc错误码转换成http状态码
func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{"msg": e.Message()})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "内部错误"})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "用户服务不可用"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "内部错误"})
			}
		}
	}
}

func removeTopStruct(fileds map[string]string) map[string]string {
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
		"error": removeTopStruct(erros.Translate(global.Trans)),
	})
}

func GetUserList(ctx *gin.Context) {

	claims, _ := ctx.Get("claims")
	zap.S().Infof("访问用户:%d", claims.(*models.CustomClaims).ID)
	pn := ctx.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSize := ctx.DefaultQuery("psize", "10")
	pSizeInt, _ := strconv.Atoi(pSize)
	//调用接口
	res, err := global.UserClinet.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    uint32(pnInt),
		Psize: uint32(pSizeInt),
	})
	if err != nil {
		zap.S().Errorw(fmt.Sprintf("[GetUserList] 查询 [用户列表失败] :%v", err.Error()))
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	result := make([]response.UserResponse, 0)
	for _, value := range res.Data {
		user := response.UserResponse{
			Id:       value.Id,
			Mobile:   value.Mobile,
			NickName: value.NickName,
			BirthDay: response.JsonTime(time.Unix(int64(value.BirthDay), 0)),
			Gender:   value.Gender,
		}
		result = append(result, user)
	}
	ctx.JSON(http.StatusOK, result)

}

func PasswordLogin(ctx *gin.Context) {
	//表单验证
	passwordLoginForm := forms.PasswordLoginForm{}
	if err := ctx.ShouldBind(&passwordLoginForm); err != nil {
		HandleValidatorError(err, ctx)
		return
	}
	//校验验证码
	if !store.Verify(passwordLoginForm.CaptchaId, passwordLoginForm.Captcha, false) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"captcha": "验证码错误",
		})
		return
	}
	//登录逻辑
	rsp, err := global.UserClinet.GetUserByMobile(context.Background(), &proto.MobileRequest{Mobile: passwordLoginForm.Mobile})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusBadRequest, map[string]string{
					"mobile": "用户不存在",
				})
			default:
				ctx.JSON(http.StatusInternalServerError, map[string]string{
					"mobile": "登陆失败",
					"err":    e.Message(),
				})
			}
			return
		}
	} else {
		if passRsp, passErr := global.UserClinet.CheckPassWord(context.Background(), &proto.PasswordCheckInfo{Password: passwordLoginForm.Password, EncryptedPassword: rsp.Password}); passErr != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]string{
				"password": "登录失败",
				"err":      passErr.Error(),
			})
		} else {
			if passRsp.Success {
				token, err := createToken(rsp)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"msg": "生成token失败",
					})
					return
				}
				ctx.JSON(http.StatusOK, gin.H{
					"id":         rsp.Id,
					"nick_name":  rsp.NickName,
					"token":      token,
					"expired_at": (time.Now().Unix() + 60*60*24) * 1000,
				})
			} else {
				ctx.JSON(http.StatusBadRequest, map[string]string{
					"password": "密码错误",
				})
			}

		}
	}

}

//用户注册
func Register(ctx *gin.Context) {
	registerFrom := forms.RegisterFrom{}
	if err := ctx.ShouldBind(&registerFrom); err != nil {
		HandleValidatorError(err, ctx)
		return
	}
	redisDB := redis.NewClient(&redis.Options{Addr: fmt.Sprintf("%s:%d", global.ServeConfig.RedisInfo.Host, global.ServeConfig.RedisInfo.Port)})
	code, err := redisDB.Get(registerFrom.Mobile).Result()
	if err == redis.Nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": "验证码错误",
		})
	} else if code != registerFrom.Code {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": "验证码错误",
		})
		return
	}

	//拨号连接用户grpc服务器
	user, err := global.UserClinet.CreateUser(context.Background(), &proto.CreateUserInfo{
		Mobile:   registerFrom.Mobile,
		Password: registerFrom.Password,
		NickName: registerFrom.Mobile,
	})
	if err != nil {
		zap.S().Errorf(fmt.Sprintf("[Register] 新建用户失败 :%v", err.Error()))
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	token, err := createToken(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成token失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":         user.Id,
		"nick_name":  user.NickName,
		"token":      token,
		"expired_at": (time.Now().Unix() + 60*60*24) * 1000,
	})
}

//创建token
func createToken(userResponse *proto.UserInfoResponse) (token string, err error) {
	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID:          uint(userResponse.Id),
		NickName:    userResponse.NickName,
		AuthorityId: uint(userResponse.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),                     //签名生效时间
			ExpiresAt: (time.Now().Unix() + 60*60*24) * 1000, //一天过期
			Issuer:    "goShop",
		},
	}
	token, err = j.CreateToken(claims)
	if err != nil {
		return "", err
	}
	return token, nil
}
