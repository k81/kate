package repo

import (
	"encoding/json"
	"errors"
	"reflect"
	"time"

	"github.com/k81/kate/config"
	"github.com/k81/orm"
)

const (
	ConfKeyDbBasic = "db/basic"
	ConfKeyDbPools = "db/pools"
)

type BasicConfig struct {
	DebugSql   bool   `json:"debug_sql"`
	DriverName string `json:"driver_name"`
	DataSource string `json:"data_source"`
}

func NewBasicConfig() (v interface{}, err error) {
	item := config.Global.GetItemByKey(ConfKeyDbBasic)
	if item == nil {
		return nil, errors.New("item not found")
	}

	conf := &BasicConfig{}
	if err = json.Unmarshal([]byte(item.Value), conf); err != nil {
		return
	}

	orm.Debug = conf.DebugSql
	v = conf
	return
}

func GetBasicConfig() *BasicConfig {
	v := config.Get(ConfKeyDbBasic)
	conf, ok := v.(*BasicConfig)
	if !ok {
		logger.Fatal(mctx, "type assert failed", "got", reflect.TypeOf(v), "expect", "*BasicConfig")
	}
	return conf
}

type PoolsConfig struct {
	MaxIdleConns    int           `json:"max_idle_conns"`
	MaxOpenConns    int           `json:"max_open_conns"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime"`
}

func NewPoolsConfig() (v interface{}, err error) {
	item := config.Global.GetItemByKey(ConfKeyDbPools)
	if item == nil {
		return nil, errors.New("item not found")
	}

	conf := &PoolsConfig{}
	if err = json.Unmarshal([]byte(item.Value), conf); err != nil {
		return
	}
	conf.ConnMaxLifetime *= time.Second
	v = conf
	return
}

func GetPoolsConfig() *PoolsConfig {
	v := config.Get(ConfKeyDbPools)
	conf, ok := v.(*PoolsConfig)
	if !ok {
		logger.Fatal(mctx, "type assert failed", "got", reflect.TypeOf(v), "expect", "*PoolsConfig")
	}
	return conf
}

func init() {
	config.Register(&config.Entry{
		Key:     ConfKeyDbBasic,
		NewFunc: NewBasicConfig,
	})
	config.Register(&config.Entry{
		Key:     ConfKeyDbPools,
		NewFunc: NewPoolsConfig,
	})
}
