package global

import (
	"goshop_api/order_web/config"
	proto "goshop_api/order_web/proto"

	ut "github.com/go-playground/universal-translator"
)

var (
	ServeConfig     config.ServerConfig
	Trans           ut.Translator
	GoodsClient     proto.GoodsClient
	OrderClient     proto.OrderClient
	InventoryClient proto.InventoryClient
	IsDebug         int
	NacosConfig     config.NacosConfig
)
