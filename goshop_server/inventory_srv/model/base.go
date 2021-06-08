package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type GoodsDetail struct {
	Goods int32
	Num   int32
}

type GoodsDetailList []GoodsDetail

type BaseModel struct {
	ID        int32     `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"column:add_time"`
	UpdatedAt time.Time `gorm:"column:update_time"`
	DeletedAt gorm.DeletedAt
	IsDeleted bool
}

func (g GoodsDetailList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

func (g *GoodsDetailList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), g)
}
