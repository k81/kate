package dynconf

import (
	"time"

	"github.com/coreos/etcd/client"
)

func NewClient(endpoints []string) (client.KeysAPI, error) {
	cfg := client.Config{
		Endpoints:               endpoints,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	c, err := client.New(cfg)
	if err != nil {
		return nil, err
	}

	kapi := client.NewKeysAPI(c)

	return kapi, nil
}
