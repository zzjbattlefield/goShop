package global

import (
	"goshop/goods_srv/config"

	"gorm.io/gorm"
)

var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig  config.NacosConfig
)

// func init() {
// 	dsn := "root:123456@tcp(192.168.1.150:3306)/mxshop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"
// 	newLogger := logger.New(
// 		log.New(os.Stdout, "\r\n", log.LstdFlags),
// 		logger.Config{
// 			SlowThreshold: time.Second,
// 			LogLevel:      logger.Info,
// 			Colorful:      true,
// 		},
// 	)
// 	// var err error
// 	// DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
// 	// 	NamingStrategy: schema.NamingStrategy{
// 	// 		SingularTable: true,
// 	// 	},
// 	// 	Logger: newLogger,
// 	// })
// 	// if err != nil {
// 	// 	panic(err)
// 	// }
// }
