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

type RedisConfig struct {
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Port     int    `mapstructure:"port" json:"port" yaml:"port"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	Expire   int    `mapstructure:"expire" json:"expire" yaml:"expire"`
}

type ServerConfig struct {
	Name       string       `mapstructure:"name" json:"name" yaml:"name"`
	Host       string       `mapstructure:"host" yaml:"host"`
	Tags       []string     `mapstructure:"tags" json:"tags" yaml:"tags"`
	MysqlInfo  MysqlConfig  `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	ConsulInfo ConsulConfig `mapstructure:"consul" json:"consul" yaml:"consul"`
	RedisInfo  RedisConfig  `mapstructure:"redis" json:"redis" yaml:"redis"`
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
