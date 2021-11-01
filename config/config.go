package config

import (
	"path/filepath"
	"sync"

	"github.com/BurntSushi/toml"
)

//配置信息
type Config struct {
	Application     application
	MySQL           mysql
	RedisCluster    rediscluster
	API             api
	Log             log
	Url             url
	UserCenter      userCenter
	Email           email
	StaticResources staticresources
	Nsq             nsq
	Encryption      encryption
	Activity        activity
	Wikibit         wikibit
	H5              h5
	GoogleMap       googleMap
	UCloud          ucloud
}

//服务配置
type application struct {
	Mode string `toml:"mode"` //模式
	Host string `toml:"host"` //服务器名
	Name string `toml:"name"` //服务名称
	Port int    `toml:"port"` //端口
}

//mysql配置
type mysql struct {
	DriverName   string `toml:"driver_name"`
	Dsn          string `toml:"dsn"`
	MaxOpenConns int    `toml:"max_open_conns"`
	MaxIdleConns int    `toml:"max_idle_conns"`
}

//redis配置
type rediscluster struct {
	Addr        string `toml:"addr"`
	Password    string `toml:"password"`
	DialTimeout int    `toml:"dial_timeout"`
	PoolSize    int    `toml:"pool_size"`
}

//用户中台配置
type userCenter struct {
	SignUrl string `toml:"sign_url"`
	User    string `toml:"user"`
}

// email
type email struct {
	Host        string `toml:"host"`
	Port        int    `toml:"port"`
	UserName    string `toml:"username"`
	Password    string `toml:"password"`
	ContentType string `toml:"content_type"`
}

//log
type log struct {
	Path string `toml:"path"`
}

//api
type api struct {
	AllowPathPrefixSkipper []string `toml:"allow_path_prefix_skipper"`
	AuthToken              string   `toml:"auth_token"`
}

// url
type url struct {
	Website string `toml:"website"`
}

// 静态资源
type staticresources struct {
	Url string `toml:"url"`
}

//nsq配置
type nsq struct {
	Host string `toml:"host"`
	Port string `toml:"port"`
}

type activity struct {
	StartDate int64 `toml:"start_date"`
	EndDate   int64 `toml:"end_date"`
}

type encryption struct {
	AesSecretKey string `toml:"aes_secret_key"`
}

type wikibit struct {
	Gateway string `toml:"gateway"`
	AppId   string `toml:"app_id"`
	Secret  string `toml:"secret"`
}

type h5 struct {
	Gateway string `toml:"gateway"`
}

type googleMap struct {
	Key string `toml:"key"`
}

type ucloud struct {
	Gateway    string `toml:"gateway"`
	PublicKey  string `toml:"public_key"`
	PrivateKey string `toml:"private_key"`
}

var (
	cfg  *Config
	once sync.Once
)

//加载配置文件
func LoadConfig() {
	once.Do(func() {
		fp, err := filepath.Abs("./config/config.toml")
		if err == nil {
			_, _ = toml.DecodeFile(fp, &cfg)
		}
	})
}

//获取配置对象
func Conf() *Config {
	return cfg
}

//获取redis集群信息
//func (c *Config) RedisClusterConfig() (cs []interface{}, is []int) {
//	cs = []interface{}{
//		c.RedisCluster.Addrs,
//		c.RedisCluster.Password,
//	}
//	is = []int{
//		c.RedisCluster.DialTimeout,
//		c.RedisCluster.PoolSize,
//	}
//	return
//}
