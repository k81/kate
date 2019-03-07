package profiling

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/k81/kate/config"
)

const (
	ConfKey = "profiling"
)

type Config struct {
	Enabled bool `json:"enabled"`
	Port    int  `json:"port"`
}

func NewConfig() (v interface{}, err error) {
	item := config.Global.GetItemByKey(ConfKey)
	if item == nil {
		return nil, errors.New("item not found")
	}

	conf := &Config{}
	if err = json.Unmarshal([]byte(item.Value), conf); err != nil {
		return
	}
	v = conf
	return
}

func GetConfig() *Config {
	v := config.Get(ConfKey)
	conf, ok := v.(*Config)
	if !ok {
		logger.Fatal(mctx, "type assert failed", "got", reflect.TypeOf(v), "expect", "*Config")
	}
	return conf
}

func init() {
	config.Register(&config.Entry{
		Key:     ConfKey,
		NewFunc: NewConfig,
	})
}
