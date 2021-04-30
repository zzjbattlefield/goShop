package global

import (
	"goshop_api/goods_web/config"
	proto "goshop_api/goods_web/proto"

	ut "github.com/go-playground/universal-translator"
)

var (
	ServeConfig config.ServerConfig
	Trans       ut.Translator
	GoodsClinet proto.GoodsClient
	IsDebug     int
	NacosConfig config.NacosConfig
)
