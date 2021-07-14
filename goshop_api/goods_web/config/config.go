package config

type GoodsSrvConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}

type ServerConfig struct {
	Name         string         `mapstructure:"name" json:"name,omitempty"`
	Port         int            `mapstructure:"port" json:"port,omitempty"`
	Host         string         `mapstructure:"host" json:"host,omitempty"`
	Tags         []string       `mapstructure:"tags" json:"tags,omitempty"`
	GoodsSrvInfo GoodsSrvConfig `mapstructure:"goods_srv" json:"goods_srv,omitempty"`
	JWTInfo      JwtConfig      `mapstructure:"jwt" json:"jwt,omitempty"`
	ConsulInfo   ConsulConfig   `mapstructure:"consul" json:"consul,omitempty"`
	JaegerInfo   JaegerConfig   `mapstructure:"jaeger" json:"jaeger,omitempty"`
}

type JaegerConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
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
