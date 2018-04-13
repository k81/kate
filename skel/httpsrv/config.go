package httpsrv

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/k81/kate/config"
	"github.com/k81/kate/log"
)

const (
	ConfKeyHttp = "http"
)

type HttpConfig struct {
	Port           int           `json:"port"`
	RequestTimeout time.Duration `json:"request_timeout_secs"`
	ReadTimeout    time.Duration `json:"read_timeout_secs"`
	WriteTimeout   time.Duration `json:"write_timeout_secs"`
	MaxHeaderBytes int           `json:"max_header_bytes"`
	MaxBodyBytes   int64         `json:"max_body_bytes"`
}

func NewHttpConfig() (v interface{}, err error) {
	item := config.Global.GetItemByKey(ConfKeyHttp)
	if item == nil {
		return nil, errors.New("item not found")
	}

	conf := &HttpConfig{}
	if err = json.Unmarshal([]byte(item.Value), conf); err != nil {
		return
	}
	conf.RequestTimeout *= time.Second
	conf.ReadTimeout *= time.Second
	conf.WriteTimeout *= time.Second
	v = conf
	return
}

func GetHttpConfig() *HttpConfig {
	v := config.Get(ConfKeyHttp)
	conf, ok := v.(*HttpConfig)
	if !ok {
		log.Fatal(mctx, "type assert failed", "got", reflect.TypeOf(v), "expect", "*HttpConfig")
	}
	return conf
}

func GetListenAddr() string {
	conf := GetHttpConfig()
	return fmt.Sprint("0.0.0.0:", conf.Port)
}

func init() {
	config.Register(&config.Entry{
		Key:     ConfKeyHttp,
		NewFunc: NewHttpConfig,
	})
}
