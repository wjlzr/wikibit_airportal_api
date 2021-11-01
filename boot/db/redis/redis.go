package redis

import (
	"time"
	"wiki_bit/config"

	"github.com/go-redis/redis"
)

var (
	client *redis.Client
)

//集群
func Connect() {
	client = redis.NewClient(&redis.Options{
		Addr:        config.Conf().RedisCluster.Addr,
		Password:    config.Conf().RedisCluster.Password,
		DialTimeout: time.Second * time.Duration(config.Conf().RedisCluster.DialTimeout),
		PoolSize:    config.Conf().RedisCluster.PoolSize,
	})

	// fmt.Println(config.RedisClusterConfig.Addrs)
	_, err := client.Ping().Result()
	if err != nil {
		panic("redis connect error")
	}
}

// client
func Client() *redis.Client {
	return client
}
