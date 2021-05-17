package initialize

import (
	"fmt"
	"goshop/userop_srv/global"

	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

func InitRedis() {
	client := goredislib.NewClient(&goredislib.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)
	global.RedisPool = pool
}
