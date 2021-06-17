package handler

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

//订单号生成 年月日时分 纳秒+用户id+随机数
func GenerateOrderSn(userId int32) string {
	rand.Seed(time.Now().UnixNano())
	now := time.Now()
	orderSn := fmt.Sprintf("%d%d%d%d%d%d%d%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Nanosecond(), userId, rand.Intn(100)+10)
	return orderSn
}

//生成rocketMq生产者
func NewProducer() (rocketmq.Producer, error) {
	p, err := rocketmq.NewProducer(producer.WithNameServer([]string{"192.168.58.130:9876"}))
	if err != nil {
		zap.S().Errorf("生成producer失败:%v\n", err.Error())
		return nil, err
	}
	if err = p.Start(); err != nil {
		zap.S().Errorf("启动producer失败:%v\n", err.Error())
		return nil, err
	}
	return p, err
}
