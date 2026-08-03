package conf

import (
	"bytes"
	"flag"
	"fmt"
	"os"

	"github.com/ao9911/go-matrix/cache/redis"
	"github.com/ao9911/go-matrix/database/gorm"
	"github.com/ao9911/go-matrix/log"
	"github.com/ao9911/go-matrix/transport/httpserver"
	"github.com/pelletier/go-toml/v2"
)

var (
	Conf     = &Config{}
	confPath string
)

var envMap = map[string]string{
	"local":   "local",
	"dev":     "dev",
	"test":    "test",
	"staging": "staging",
	"prod":    "prod",
}

type Config struct {
	Env        string
	Log        *log.Config
	HttpServer *httpserver.Config
	Mysql      *gorm.Config `toml:"mysql"`
	Redis      *redis.Config
	Auth       *AuthConfig
	RateLimit  *RateLimitConfig
}

type AuthConfig struct {
	AccessExpire  int64 `toml:"access_expire"`
	RefreshExpire int64 `toml:"refresh_expire"`
}

type RateLimitConfig struct {
	MaxQPS   int `toml:"maxQps"`   // 平均每秒最多请求数
	MaxBurst int `toml:"maxBurst"` // 突发请求数，允许瞬间超过平均每秒请求数的请求数
}

func init() {
	flag.StringVar(&confPath, "conf", "", "config file path")
}

func Init() {
	// TODO 从配置中心获取配置文件
	deployEnv := os.Getenv("DEPLOY_ENV")
	if deployEnv == "" {
		deployEnv = "local"
	}

	if confPath == "" {
		env, ok := envMap[deployEnv]
		if !ok {
			panic(fmt.Sprintf("invalid deploy env: %s", deployEnv))
		}
		confPath = fmt.Sprintf("./conf/%s.toml", env)
	}
	data, err := os.ReadFile(confPath)
	if err != nil {
		panic(fmt.Sprintf("read config file %q: %v", confPath, err))
	}
	if err := toml.NewDecoder(bytes.NewReader(data)).Decode(Conf); err != nil {
		panic(fmt.Sprintf("decode config file %q: %v", confPath, err))
	}
	Conf.Env = deployEnv
}
