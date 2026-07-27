package conf

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/pelletier/go-toml/v2"
)

func init() {
	data, err := os.ReadFile("./local.toml")
	if err != nil {
		panic(fmt.Sprintf("read config file %q: %v", "local.toml", err))
	}
	if err := toml.NewDecoder(bytes.NewReader(data)).Decode(Conf); err != nil {
		panic(fmt.Sprintf("decode config file %q: %v", "local.toml", err))
	}
}

// go test -v -run TestLog
func TestLog(t *testing.T) {
	t.Logf("Log.AppName: %s", Conf.Log.AppName)
	t.Logf("Log.LogPath: %s", Conf.Log.LogPath)
	t.Logf("Log.Debug: %v", Conf.Log.Debug)
	t.Logf("Log.MultiFile: %v", Conf.Log.MultiFile)
}

// go test -v -run TestHttpServer
func TestHttpServer(t *testing.T) {
	t.Logf("HttpServer.Addr: %s", Conf.HttpServer.Addr)
	t.Logf("HttpServer.ReadTimeout: %v", Conf.HttpServer.ReadTimeout)
	t.Logf("HttpServer.WriteTimeout: %v", Conf.HttpServer.WriteTimeout)
}

// go test -v -run TestMysql
func TestMysql(t *testing.T) {
	t.Logf("Mysql.DebugMode: %v", Conf.Mysql.DebugMode)
	t.Logf("Mysql.DriverName: %s", Conf.Mysql.DriverName)
	t.Logf("Mysql.DSN: %s", Conf.Mysql.DSN)
	t.Logf("Mysql.MaxIdleConns: %d", Conf.Mysql.MaxIdleConns)
	t.Logf("Mysql.MaxOpenConns: %d", Conf.Mysql.MaxOpenConns)
	t.Logf("Mysql.MaxLifetime: %v", Conf.Mysql.MaxLifetime)
	t.Logf("Mysql.MaxIdleTime: %v", Conf.Mysql.MaxIdleTime)
}

// go test -v -run TestRedis
func TestRedis(t *testing.T) {
	t.Logf("Redis.Addrs: %s", Conf.Redis.Addrs)
	t.Logf("Redis.Password: %s", Conf.Redis.Password)
	t.Logf("Redis.PoolSize: %d", Conf.Redis.PoolSize)
}

func TestAuth(t *testing.T) {
	t.Logf("Auth.Expire: %d", Conf.Auth.Expire)
}
