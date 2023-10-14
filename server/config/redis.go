package config

import "time"

type Redis struct {
	Addr     string `mapstructure:"addr" json:"addr" yaml:"addr"`             // 服务器地址:端口
	Password string `mapstructure:"password" json:"password" yaml:"password"` // 密码
	DB       int    `mapstructure:"db" json:"db" yaml:"db"`                   // redis的哪个数据库
	Size     int    `mapstructure:"size" json:"size" yaml:"size"`             // size
}

type RabbitMQ struct {
	Addr             string        `mapstructure:"addr" json:"addr" yaml:"addr"`                                     // 服务器地址:端口
	Port             string        `mapstructure:"port" json:"port" yaml:"port"`                                     // 端口
	Password         string        `mapstructure:"password" json:"password" yaml:"password"`                         // 密码
	Username         string        `mapstructure:"username" json:"username" yaml:"username"`                         // 账户
	PoolSize         int           `mapstructure:"poolSize" json:"poolSize" yaml:"poolSize"`                         // poolSize
	RetryReConnDelay time.Duration `mapstructure:"retryReConnDelay" json:"retryReConnDelay" yaml:"retryReConnDelay"` // retryReConnDelay
}
