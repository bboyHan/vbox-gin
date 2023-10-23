package config

import "time"

type RabbitMQ struct {
	Addr             string        `mapstructure:"addr" json:"addr" yaml:"addr"`                                     // 服务器地址:端口
	Port             string        `mapstructure:"port" json:"port" yaml:"port"`                                     // 端口
	Password         string        `mapstructure:"password" json:"password" yaml:"password"`                         // 密码
	Username         string        `mapstructure:"username" json:"username" yaml:"username"`                         // 账户
	PoolSize         int           `mapstructure:"poolSize" json:"poolSize" yaml:"poolSize"`                         // poolSize
	RetryReConnDelay time.Duration `mapstructure:"retryReConnDelay" json:"retryReConnDelay" yaml:"retryReConnDelay"` // retryReConnDelay
}
