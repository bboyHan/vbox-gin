package config

type System struct {
	Env           string `mapstructure:"env" json:"env" yaml:"env"`                // 环境值
	DbType        string `mapstructure:"db-type" json:"db-type" yaml:"db-type"`    // 数据库类型:mysql(默认)|sqlite|sqlserver|postgresql
	OssType       string `mapstructure:"oss-type" json:"oss-type" yaml:"oss-type"` // Oss类型
	RouterPrefix  string `mapstructure:"router-prefix" json:"router-prefix" yaml:"router-prefix"`
	Addr          int    `mapstructure:"addr" json:"addr" yaml:"addr"` // 端口值
	LimitCountIP  int    `mapstructure:"iplimit-count" json:"iplimit-count" yaml:"iplimit-count"`
	LimitTimeIP   int    `mapstructure:"iplimit-time" json:"iplimit-time" yaml:"iplimit-time"`
	UseMultipoint bool   `mapstructure:"use-multipoint" json:"use-multipoint" yaml:"use-multipoint"` // 多点登录拦截
	UseRedis      bool   `mapstructure:"use-redis" json:"use-redis" yaml:"use-redis"`                // 使用redis
	UseMongo      bool   `mapstructure:"use-mongo" json:"use-mongo" yaml:"use-mongo"`                // 使用mongo
	UseMQ         bool   `mapstructure:"use-mq" json:"use-mq" yaml:"use-mq"`                         // 使用mq
	UseMQTask     bool   `mapstructure:"use-mq-task" json:"use-mq-task" yaml:"use-mq-task"`          // vbox mq task
	TimerTask     bool   `mapstructure:"timer-task" json:"timer-task" yaml:"timer-task"`             // 定时任务开关
}
