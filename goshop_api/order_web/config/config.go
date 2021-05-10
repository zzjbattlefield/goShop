package config

type OrderSrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

type ServerConfig struct {
	Name          string         `mapstructure:"name" json:"name"`
	Port          int            `mapstructure:"port" json:"port"`
	Host          string         `mapstructure:"host" json:"host"`
	Tags          []string       `mapstructure:"tags" json:"tags"`
	OrderSrvInfo  OrderSrvConfig `mapstructure:"order_srv" json:"order_srv"`
	GoodsSrvInfo  OrderSrvConfig `mapstructure:"goods_srv" json:"goods_srv"`
	InventoryInfo OrderSrvConfig `mapstructure:"inventory_srv" json:"inventory_srv"`
	JWTInfo       JwtConfig      `mapstructure:"jwt" json:"jwt"`
	ConsulInfo    ConsulConfig   `mapstructure:"consul" json:"consul"`
}

type JwtConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type NacosConfig struct {
	Host      string `mapstructure:"host" json:"host"`
	Port      int    `mapstructure:"port" json:"port"`
	NameSpace string `mapstructure:"namespace" json:"namespace"`
	User      string `mapstructure:"user" json:"user"`
	Password  string `mapstructure:"password" json:"password"`
	DataId    string `mapstructure:"dataid" json:"dataid"`
	Group     string `mapstructure:"group" json:"group"`
}
