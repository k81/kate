package app

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/k81/kate/config"
	"github.com/k81/log"
)

const (
	ConfKeyLog = "log"
)

type LogConfig struct {
	Level log.LevelName `json:"level"`
}

func NewLogConfig() (v interface{}, err error) {
	item := config.Global.GetItemByKey(ConfKeyLog)
	if item == nil {
		return nil, errors.New("item not found")
	}

	conf := &LogConfig{}
	if err = json.Unmarshal([]byte(item.Value), conf); err != nil {
		return
	}
	v = conf
	return
}

func OnLogConfigUpdate(v interface{}) {
	conf, ok := v.(*LogConfig)
	if !ok {
		logger.Fatal(mctx, "type assert failed", "got", reflect.TypeOf(v), "expect", "*LogConfig")
	}

	logger.Info(mctx, "log level changed", "old_level", log.GetLevel().String(), "new_level", conf.Level)

	log.SetLevelByName(conf.Level)
}

func GetLogLevel() log.LevelName {
	v := config.Get(ConfKeyLog)
	conf, ok := v.(*LogConfig)
	if !ok {
		logger.Fatal(mctx, "type assert failed", "got", reflect.TypeOf(v), "expect", "*LogConfig")
	}
	return conf.Level
}

func init() {
	config.Register(&config.Entry{
		Key:          ConfKeyLog,
		NewFunc:      NewLogConfig,
		OnUpdateFunc: OnLogConfigUpdate,
	})
}
