package redismgr

import (
	"encoding/json"
	"errors"
	"reflect"
	"time"

	"github.com/k81/kate/config"
	"github.com/k81/kate/log"
)

const (
	ConfKeyAddrs = "redis/addrs"
	ConfKeyPools = "redis/pools"
)

// redis服务的接入地址配置
type AddrsConfig struct {
	Addrs []string `json:"addrs"`
}

// redis连接池的参数配置
type PoolsConfig struct {
	MaxIdle        int           `json:"max_idle"`
	MaxActive      int           `json:"max_active"`
	ConnectTimeout time.Duration `json:"connect_timeout"` // connect timeut in seconds
	ReadTimeout    time.Duration `json:"read_timeout"`    // read timeout in seconds
	WriteTimeout   time.Duration `json:"write_timeout"`   // write timeout in seconds
	IdleTimeout    time.Duration `json:"idle_timeout"`    // idle timeout in seconds
	Wait           bool          `json:"wait"`
}

func NewAddrsConfig() (v interface{}, err error) {
	item := config.Global.GetItemByKey(ConfKeyAddrs)
	if item == nil {
		return nil, errors.New("item not found")
	}

	conf := &AddrsConfig{}
	if err = json.Unmarshal([]byte(item.Value), conf); err != nil {
		return
	}
	v = conf
	return
}

func NewPoolsConfig() (v interface{}, err error) {
	item := config.Global.GetItemByKey(ConfKeyPools)
	if item == nil {
		return nil, errors.New("item not found")
	}

	conf := &PoolsConfig{}
	if err = json.Unmarshal([]byte(item.Value), conf); err != nil {
		return
	}
	conf.ConnectTimeout *= time.Second
	conf.ReadTimeout *= time.Second
	conf.WriteTimeout *= time.Second
	conf.IdleTimeout *= time.Second
	v = conf
	return
}

func GetAddrs() []string {
	v := config.Get(ConfKeyAddrs)
	conf, ok := v.(*AddrsConfig)
	if !ok {
		log.Fatal(mctx, "type assert failed", "got", reflect.TypeOf(v), "expect", "*AddrsConfig")
	}
	return conf.Addrs
}

func GetPoolsConfig() *PoolsConfig {
	v := config.Get(ConfKeyPools)
	conf, ok := v.(*PoolsConfig)
	if !ok {
		log.Fatal(mctx, "type assert failed", "got", reflect.TypeOf(v), "expect", "*PoolsConfig")
	}
	return conf
}

func init() {
	config.Register(&config.Entry{
		Key:     ConfKeyAddrs,
		NewFunc: NewAddrsConfig,
	})

	config.Register(&config.Entry{
		Key:          ConfKeyPools,
		NewFunc:      NewPoolsConfig,
		OnUpdateFunc: OnPoolsConfigUpdate,
	})
}
