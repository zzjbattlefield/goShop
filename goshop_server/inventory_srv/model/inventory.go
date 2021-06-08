package model

type Inventory struct {
	BaseModel
	Goods   int32 `gorm:"type:int;index"`
	Stocks  int32 `gorm:"type:int"`
	Version int32 `gorm:"type:int"` //分布式锁的乐观锁
}

type StockSellDetail struct {
	OrderSn string          `gorm:"type:varchar(100);"`
	Status  int32           `gorm:"type:int;"` //1代表已扣减 2.已归还
	Detail  GoodsDetailList `gorm:"type:varchar(100);"`
}

func (StockSellDetail) TableName() string {
	return "stockselldetail"
}
