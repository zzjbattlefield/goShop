package global

import (
	"goshop_api/userop_web/config"
	proto "goshop_api/userop_web/proto"

	ut "github.com/go-playground/universal-translator"
)

var (
	ServeConfig   config.ServerConfig
	Trans         ut.Translator
	MessageClient proto.MessageClient
	AddressClient proto.AddressClient
	UserFavClient proto.UserFavClient

	GoodsClient proto.GoodsClient
	IsDebug     int
	NacosConfig config.NacosConfig
)
