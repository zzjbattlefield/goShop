package global

import (
	"goshop_api/user_web/config"
	proto "goshop_api/user_web/proto"

	ut "github.com/go-playground/universal-translator"
)

var (
	ServeConfig config.ServerConfig
	Trans       ut.Translator
	UserClinet  proto.UserClient
	IsDebug     int
	NacosConfig config.NacosConfig
)
