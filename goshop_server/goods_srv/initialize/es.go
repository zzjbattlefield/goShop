package initialize

import (
	"context"
	"fmt"
	"goshop/goods_srv/global"
	"goshop/goods_srv/model"
	"log"
	"os"

	"github.com/olivere/elastic/v7"
)

//初始化连接
func InitEs() {
	var err error
	host := fmt.Sprintf("http://%s:%d", global.ServerConfig.EsInfo.Host, global.ServerConfig.EsInfo.Port)
	log := log.New(os.Stdout, "mxshop", log.LstdFlags)
	global.EsClient, err = elastic.NewClient(elastic.SetURL(host), elastic.SetSniff(false), elastic.SetTraceLog(log))
	if err != nil {
		panic(err)
	}
	//新建mapping和index
	exists, err := global.EsClient.IndexExists(model.EsGoods{}.GetIndexName()).Do(context.Background())
	if err != nil {
		panic(err)
	}
	if !exists {
		_, err := global.EsClient.CreateIndex(model.EsGoods{}.GetIndexName()).BodyString(model.EsGoods{}.GetMapping()).Do(context.Background())
		if err != nil {
			panic(err)
		}
	}
}
