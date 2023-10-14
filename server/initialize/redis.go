package initialize

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func Redis() *redis.Client {
	//redisCfg := global.GVA_CONFIG.Redis

	//client := redis.NewClient(&redis.Options{
	//	Addr:     redisCfg.Addr,
	//	Password: redisCfg.Password, // no password set
	//	DB:       redisCfg.DB,       // use default DB
	//	PoolSize: redisCfg.Size,     // size
	//})

	client := RedisClient()

	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		global.GVA_LOG.Error("redis connect ping failed, err:", zap.Error(err))
	} else {
		global.GVA_LOG.Info("redis connect ping response:", zap.String("pong", pong))
		global.GVA_REDIS = client
	}

	return client
}

/*
// PoolInitRedis redis pool
func PoolInitRedis() *redigo.Pool {
	redisCfg := global.GVA_CONFIG.Redis

	return &redigo.Pool{
		MaxIdle:     2, //空闲数
		IdleTimeout: 240 * time.Second,
		MaxActive:   3, //最大数
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.Dial("tcp", redisCfg.Addr)
			if err != nil {
				return nil, err
			}
			if redisCfg.Password != "" {
				if _, err := c.Do("AUTH", redisCfg.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}*/

func RedisClient() *redis.Client {
	redisCfg := global.GVA_CONFIG.Redis
	return redis.NewClient(&redis.Options{
		//连接信息
		Network:  "tcp",             //网络类型，tcp or unix，默认tcp
		Addr:     redisCfg.Addr,     //主机名+冒号+端口，默认localhost:6379
		Password: redisCfg.Password, //密码
		DB:       0,                 // redis数据库index

		//连接池容量及闲置连接数量
		PoolSize:     redisCfg.Size, // 连接池最大socket连接数，默认为4倍CPU数， 4 * runtime.NumCPU
		MinIdleConns: 100,           //在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；。
		//MaxIdleConns: 50,

		//超时
		//DialTimeout:  5 * time.Second, //连接建立超时时间，默认5秒。
		//ReadTimeout:  3 * time.Second, //读超时，默认3秒， -1表示取消读超时
		//WriteTimeout: 3 * time.Second, //写超时，默认等于读超时
		//PoolTimeout:  4 * time.Second, //当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒。

		//命令执行失败时的重试策略
		//MaxRetries:      0,                      // 命令执行失败时，最多重试多少次，默认为0即不重试
		//MinRetryBackoff: 8 * time.Millisecond,   //每次计算重试间隔时间的下限，默认8毫秒，-1表示取消间隔
		//MaxRetryBackoff: 512 * time.Millisecond, //每次计算重试间隔时间的上限，默认512毫秒，-1表示取消间隔

	})
}
