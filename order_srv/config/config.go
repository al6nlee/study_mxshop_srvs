package config

type MysqlConfig struct {
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Port     int    `mapstructure:"port" json:"port" yaml:"port"`
	Name     string `mapstructure:"db" json:"db" yaml:"db"`
	User     string `mapstructure:"user" json:"user" yaml:"user"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host" yaml:"host"`
	Port int    `mapstructure:"port" json:"port" yaml:"port"`
}

type GoodsSrvConfig struct {
	Name string `mapstructure:"name" json:"name" yaml:"name"`
}

type InventorySrvConfig struct {
	Name string `mapstructure:"name" json:"name" yaml:"name"`
}

type ServerConfig struct {
	Name       string       `mapstructure:"name" json:"name" yaml:"name"`
	Host       string       `mapstructure:"host" yaml:"host"`
	Tags       []string     `mapstructure:"tags" json:"tags" yaml:"tags"`
	MysqlInfo  MysqlConfig  `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	ConsulInfo ConsulConfig `mapstructure:"consul" json:"consul" yaml:"consul"`

	GoodsSrvInfo     GoodsSrvConfig     `mapstructure:"goods_srv" json:"goods_srv" yaml:"goods_srv"`
	InventorySrvInfo InventorySrvConfig `mapstructure:"inventory_srv" json:"inventory_srv" yaml:"inventory_srv"`
}

type NacosConfig struct {
	Host      string `mapstructure:"host" yaml:"host"`
	Port      uint64 `mapstructure:"port" yaml:"port"`
	Namespace string `mapstructure:"namespace" yaml:"namespace"`
	User      string `mapstructure:"user" yaml:"user"`
	Password  string `mapstructure:"password" yaml:"password"`
	DataId    string `mapstructure:"dataid" yaml:"dataid"`
	Group     string `mapstructure:"group" yaml:"group"`
}
