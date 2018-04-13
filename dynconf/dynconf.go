package dynconf

import (
	"context"
	"sync"
	"time"

	"github.com/coreos/etcd/client"
)

var (
	mu            sync.Mutex
	defaultClient client.KeysAPI
)

var (
	serverAddrs []string
)

type Item struct {
	Key     string
	Value   string
	Version uint64
}

type ItemHandle struct {
	ctx    context.Context
	client client.KeysAPI
	Key    string
}

func Init(endpoints []string) error {
	serverAddrs = endpoints
	return doInit()
}

func doInit() error {
	var err error
	defaultClient, err = NewClient(serverAddrs)
	if err != nil {
		return err
	}
	return nil
}

func GetAll(ctx context.Context, key string) (map[string]*Item, error) {
	mu.Lock()
	defer mu.Unlock()

	if defaultClient == nil {
		err := doInit()
		if err != nil {
			return nil, err
		}
	}

	opt := &client.GetOptions{
		Recursive: true,
		Sort:      true,
		Quorum:    true,
	}

	resp, err := defaultClient.Get(ctx, key, opt)
	if err != nil {
		return nil, err
	}

	vars := make(map[string]*Item)
	walk(resp.Node, vars)

	return vars, nil
}

func Get(ctx context.Context, key string) (*Item, error) {
	mu.Lock()
	defer mu.Unlock()

	if defaultClient == nil {
		err := doInit()
		if err != nil {
			return nil, err
		}
	}

	opt := &client.GetOptions{
		Recursive: false,
		Sort:      false,
		Quorum:    true,
	}

	resp, err := defaultClient.Get(ctx, key, opt)
	if err != nil {
		return nil, err
	}

	item := &Item{
		Key:     resp.Node.Key,
		Value:   resp.Node.Value,
		Version: resp.Node.ModifiedIndex,
	}

	return item, nil
}

func Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	mu.Lock()
	defer mu.Unlock()

	if defaultClient == nil {
		err := doInit()
		if err != nil {
			return err
		}
	}

	opt := &client.SetOptions{
		TTL: ttl,
	}
	_, err := defaultClient.Set(ctx, key, value, opt)
	if err != nil {
		return err
	}
	return nil
}

func NewItemHandle(ctx context.Context, key string) (*ItemHandle, error) {
	var (
		client client.KeysAPI
		err    error
	)

	if client, err = NewClient(serverAddrs); err != nil {
		return nil, err
	}

	return &ItemHandle{
		ctx:    ctx,
		client: client,
		Key:    key,
	}, nil
}

func walk(node *client.Node, vars map[string]*Item) {
	if node != nil {
		key := node.Key
		if !node.Dir {
			vars[key] = &Item{
				Key:     node.Key,
				Value:   node.Value,
				Version: node.ModifiedIndex,
			}
		} else {
			for _, node := range node.Nodes {
				walk(node, vars)
			}
		}
	}
}

func (ith *ItemHandle) GetAll() (map[string]*Item, error) {
	opt := &client.GetOptions{
		Recursive: true,
		Sort:      true,
		Quorum:    true,
	}

	resp, err := ith.client.Get(ith.ctx, ith.Key, opt)
	if err != nil {
		return nil, err
	}

	vars := make(map[string]*Item)
	walk(resp.Node, vars)

	return vars, nil
}

func (ith *ItemHandle) Get() (*Item, error) {
	opt := &client.GetOptions{
		Recursive: true,
		Sort:      true,
		Quorum:    true,
	}

	resp, err := ith.client.Get(ith.ctx, ith.Key, opt)
	if err != nil {
		return nil, err
	}

	item := &Item{
		Key:     resp.Node.Key,
		Value:   resp.Node.Value,
		Version: resp.Node.ModifiedIndex,
	}

	return item, nil
}

func (ith *ItemHandle) Set(value string, ttl time.Duration) error {
	opt := &client.SetOptions{
		TTL: ttl,
	}
	_, err := ith.client.Set(ith.ctx, ith.Key, value, opt)
	if err != nil {
		return err
	}
	return nil
}
