package model

import (
	"context"
	"goshop/goods_srv/global"
	"strconv"

	"gorm.io/gorm"
)

type Category struct {
	BaseModel
	Name             string      `gorm:"type:varchar(20);not null"`
	Level            int32       `gorm:"type:int;not null;default:1"`
	IsTab            bool        `gorm:"default:false;not null"`
	SubCategory      []*Category `gorm:"foreignKey:ParentCategoryID;references:ID"`
	ParentCategoryID int32       `gorm:"type:int;not null"`
	ParentCategory   *Category
}

//品牌表
type Brands struct {
	BaseModel
	Name string `gorm:"type:varchar(20);not null"`
	Logo string `gorm:"type:varchar(255);not null default:''"`
}

//商品品牌关系关系表
type GoodsCategoryBrand struct {
	BaseModel
	CategoryID int32 `gorm:"type:int;index:idx_category_brand,nuique"`
	Category   Category
	BrandsID   int32 `gorm:"type:int;index:idx_category_brand,nuique"`
	Brands     Brands
}

//重载表名
// func (GoodsCategoryBrand) TableName() string {
// 	return "goodscategorybrand"
// }

//轮播图
type Banner struct {
	BaseModel
	Image string `gorm:"type:varchar(200);not null"`
	Url   string `gorm:"type:varchar(200);not null"`
	Index int    `gorm:"type:int;default:1;not null"`
}

//商品信息表
type Goods struct {
	BaseModel

	CategoryID int32 `gorm:"type:int;not null"`
	Category   Category
	BrandsID   int32 `gorm:"type:int;not null"`
	Brands     Brands

	OnSale   bool `gorm:"default:false;not null"`
	ShipFree bool `gorm:"default:false;not null"`
	IsNew    bool `gorm:"default:false;not null"`
	IsHot    bool `gorm:"default:false;not null"`

	Name    string `gorm:"type:varchar(250);not null"`
	GoodsSn string `gorm:"type:varchar(50);not null"`

	ClickNum        int32   `gorm:"type:int;default:0;not null"`
	SoldNum         int32   `gorm:"type:int;default:0;not null"`
	FavNum          int32   `gorm:"type:int;default:0;not null"`
	Stocks          int32   `gorm:"type:int;default:0;not null"`
	MarketPrice     float32 `gorm:"not null"`
	ShopPrice       float32 `gorm:"not null"`
	GoodsBrief      string  `gorm:"type:varchar(50);not null"` //商品简介
	GoodsFrontImage string  `gorm:"type:varchar(200);not null"`

	Images     GormList `gorm:"type:varchar(1000);not null"` //商品图片
	DescImages GormList `gorm:"type:varchar(1000);not null"` //商品详细图片
}

func (g *Goods) AfterCreate(tx *gorm.DB) error {
	esModel := EsGoods{
		ID:          g.ID,
		CategoryID:  g.CategoryID,
		BrandsID:    g.BrandsID,
		OnSale:      g.OnSale,
		ShipFree:    g.ShipFree,
		IsNew:       g.IsNew,
		IsHot:       g.IsHot,
		Name:        g.Name,
		ClickNum:    g.ClickNum,
		SoldNum:     g.SoldNum,
		FavNum:      g.FavNum,
		MarketPrice: g.MarketPrice,
		GoodsBrief:  g.GoodsBrief,
		ShopPrice:   g.ShopPrice,
	}
	_, err := global.EsClient.Index().Index(EsGoods{}.GetIndexName()).BodyJson(esModel).Id(strconv.Itoa(int(esModel.ID))).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}
