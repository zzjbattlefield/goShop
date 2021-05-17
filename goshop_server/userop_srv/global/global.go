package global

import (
	"goshop/userop_srv/config"

	"github.com/go-redsync/redsync/v4/redis"
	"gorm.io/gorm"
)

var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig  config.NacosConfig
	RedisPool    redis.Pool
)
