package initialize

import (
	"encoding/json"
	"fmt"
	"goshop_api/user_web/global"
	"os"
	"strconv"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func InitConfig() {
	configFilePrefix := "config"
	global.IsDebug, _ = strconv.Atoi(os.Getenv("GOSHOP_DEBUG"))
	var configFileName string
	if global.IsDebug == 1 {
		configFileName = fmt.Sprintf("%s-debug.yaml", configFilePrefix)
	} else {
		configFileName = fmt.Sprintf("%s-online.yaml", configFilePrefix)
	}
	v := viper.New()
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		zap.S().Panicf("[InitConfig] 读取配置失败 : %s ", err.Error())
	}
	if err := v.Unmarshal(&global.NacosConfig); err != nil {
		zap.S().Panicf("[InitConfig] 反序列化设置失败 : %s", err.Error())
	}
	sc := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Host,
			Port:   uint64(global.NacosConfig.Port),
		},
	}
	cc := constant.ClientConfig{
		NamespaceId: global.NacosConfig.NameSpace, // 如果需要支持多namespace，我们可以场景多个client,
		// 它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "../tmp/nacos/log",
		CacheDir:            "../tmp/nacos/cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}
	config_client, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		panic(err)
	}
	content, err := config_client.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group})
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(content), &global.ServeConfig)
	zap.S().Info(content)
	zap.S().Info("jwt:", global.ServeConfig.JWTInfo.SigningKey)
	if err != nil {
		zap.S().Fatalf("读取nacos配置失败:%s", err)
	}
	_ = config_client.ListenConfig(vo.ConfigParam{
		DataId: "user-web",
		Group:  "dev",
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("配置文件产生变化")
			fmt.Println("group:" + group + ", dataId:" + dataId + ", data:" + data)
		},
	})
	// zap.S().Info("配置信息:%v", global.ServeConfig)
	// //监听配置文件的变动
	// v.WatchConfig()
	// v.OnConfigChange(func(e fsnotify.Event) {
	// 	zap.S().Infof("配置文件产生变化:%v", e.Name)
	// 	_ = v.ReadInConfig()
	// 	_ = v.Unmarshal(&global.ServeConfig)
	// 	zap.S().Info("配置信息:%v", global.ServeConfig)
	// })
}
