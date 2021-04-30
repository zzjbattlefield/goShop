package global

import (
	"goshop/order_srv/config"
	"goshop/order_srv/proto"

	"github.com/go-redsync/redsync/v4/redis"
	"gorm.io/gorm"
)

var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig  config.NacosConfig
	RedisPool    redis.Pool

	GoodsSrvClient     proto.GoodsClient
	InventorySrvClient proto.InventoryClient
)
