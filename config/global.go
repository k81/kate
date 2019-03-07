package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/k81/kate/dynconf"
	"github.com/k81/kate/utils"
	"github.com/k81/log"
)

var CONF_ROOT string

var CACHE_FILE_PREFIX string = `gcached_`

var CACHE_FILE_PATH string

// 解析配置的callback
type NewFunc func() (interface{}, error)

// 配置更新的callback
type OnUpdateFunc func(v interface{})

// 模块的配置描述条目
type Entry struct {
	Key          string
	NewFunc      NewFunc
	OnUpdateFunc OnUpdateFunc
	Restart      bool
	IsPrefix     bool
	Optional     bool

	loaded bool
	value  atomic.Value
}

type UpdatedEntry struct {
	entry *Entry
	value interface{}
}

/*************************************************
Description: 判断key是否与当前配置描述条目匹配
Input:		 key
Output:
Return:		 匹配，true；不匹配，false
Others:
*************************************************/
func (entry *Entry) match(key string) bool {
	if !entry.IsPrefix {
		return key == entry.Key
	}
	return strings.HasPrefix(key, entry.Key)
}

// 全局配置
type globalConfig struct {
	items    atomic.Value
	entryMap map[string]*Entry
	watching bool
	stopChan chan bool
}

var Global *globalConfig

/*************************************************
Description: 初始化
Input:
Output:
Return:
Others:
*************************************************/
func init() {
	Global = &globalConfig{
		entryMap: make(map[string]*Entry),
		stopChan: make(chan bool, 1),
	}
}

/*************************************************
Description: 注册托管的配置描述条目
Input:
	entry	 模块的配置描述条目
Output:
Return:
Others:
*************************************************/
func (g *globalConfig) Add(entry *Entry) {
	if entry.NewFunc == nil {
		log.Fatal(mctx, "add entry", "key", entry.Key, "error", "entry.NewFunc == nil")
	}

	e := &Entry{
		Key:          entry.Key,
		NewFunc:      entry.NewFunc,
		OnUpdateFunc: entry.OnUpdateFunc,
		Restart:      entry.Restart,
		IsPrefix:     entry.IsPrefix,
		Optional:     entry.Optional,
	}
	g.entryMap[entry.Key] = e
}

/*************************************************
Description: 获取key对应的解析后配置信息
Input:
	key		 key
Output:
Return:		 和key关联的配置信息
Others:
*************************************************/
func (g *globalConfig) GetEntryData(key string) interface{} {
	entry := g.entryMap[key]
	if entry == nil {
		return nil
	}
	return entry.value.Load()
}

/*************************************************
Description: 设置key对应的解析后配置信息
Input:
	key		 key
	v		 经过解析的配置信息
Output:
Return:
Others:
*************************************************/
func (g *globalConfig) SetEntryData(key string, v interface{}) {
	entry := g.entryMap[key]
	if entry == nil {
		return
	}
	entry.value.Store(v)
}

/*************************************************
Description: 根据key查找原始的配置项
Input:
	key		 原始配置项的key
Output:
Return:		 配置项
Others:
*************************************************/
func (g *globalConfig) GetItemByKey(key string) *dynconf.Item {
	items := g.items.Load().(map[string]*dynconf.Item)
	item := items[key]
	if item == nil {
		return nil
	}

	clone := &dynconf.Item{
		Key:     item.Key,
		Value:   item.Value,
		Version: item.Version,
	}
	return clone
}

/*************************************************
Description: 获取所有前缀为prefix的原始配置项
Input:
	prefix	 要获取的key的前缀
Output:
Return:		 所有前缀为prefix的item列表
Others:
*************************************************/
func (g *globalConfig) GetItemsByPrefix(prefix string) map[string]*dynconf.Item {
	items := g.items.Load().(map[string]*dynconf.Item)

	vars := make(map[string]*dynconf.Item)

	for _, item := range items {
		if strings.HasPrefix(item.Key, prefix) {
			vars[item.Key] = &dynconf.Item{
				Key:     item.Key,
				Value:   item.Value,
				Version: item.Version,
			}
		}
	}

	return vars
}

/*************************************************
Description: 查找与key匹配的配置描述条目
Input:
	key		 key
Output:
Return:		 如果找到，entry；未找到，nil
Others:
*************************************************/
func (g *globalConfig) matchEntry(key string) *Entry {
	for _, entry := range g.entryMap {
		if entry.match(key) {
			return entry
		}
	}
	return nil
}

/*************************************************
Description: 从配置后端加载所有前缀为prefix的原始配置项
Input:
	prefix	 原始配置项的前缀
Output:
Return:		 成功时，返回所有符合的配置项，且error==nil; 失败时，返回nil，且error!=nil
Others:
*************************************************/
func (g *globalConfig) loadItemsFromServer(prefix string) (map[string]*dynconf.Item, error) {
	items, err := dynconf.GetAll(mctx, prefix)
	if err != nil {
		return nil, err
	}

	trimed := make(map[string]*dynconf.Item, len(items))
	for _, item := range items {
		item.Key = strings.TrimPrefix(item.Key, prefix)
		item.Value = strings.TrimSpace(item.Value)
		trimed[item.Key] = item
	}

	return trimed, nil
}

func (g *globalConfig) parseItems(items map[string]*dynconf.Item) error {
	for _, item := range items {
		entry := g.matchEntry(item.Key)
		if entry == nil {
			continue
		}

		entry.loaded = true
		log.Info(mctx, "global config", "key", item.Key, "version", item.Version, "value", item.Value)
	}

	g.items.Store(items)

	for key, entry := range g.entryMap {
		if !entry.loaded {
			if !entry.Optional {
				log.Error(mctx, "global config not found", "key", key)
				return fmt.Errorf("config key not found: %s", key)
			}
			continue
		}

		conf, err := entry.NewFunc()
		if err != nil {
			log.Error(mctx, "global config parse failed", "key", key, "error", err)
			return fmt.Errorf("config key parse failed: %s", key)
		}
		entry.value.Store(conf)
	}

	return nil
}

func (g *globalConfig) initFromServer() (err error) {
	var items map[string]*dynconf.Item

	if err = dynconf.Init(Local.EtcdAddrs); err != nil {
		return
	}

	if items, err = g.loadItemsFromServer(CONF_ROOT); err != nil {
		return
	}

	if err = g.parseItems(items); err != nil {
		return
	}

	return
}

func (g *globalConfig) initFromCacheFile() (err error) {
	var items map[string]*dynconf.Item

	if items, err = g.loadItemsFromCacheFile(); err != nil {
		log.Error(mctx, "load global config from cache file", "filename", CACHE_FILE_PATH, "error", err)
		return
	}

	if err = g.parseItems(items); err != nil {
		return
	}

	return
}

/*************************************************
Description: 加载所有注册的配置信息，并解析，然后分别存储在对应的key下面
Input:
Output:
Return:
Others:
*************************************************/
func (g *globalConfig) init() {
	CONF_ROOT = fmt.Sprint(Local.EtcdPrefix, "conf/")
	CACHE_FILE_PATH = path.Join(path.Dir(configFilePath), fmt.Sprint(CACHE_FILE_PREFIX, app, ".yaml"))

	var loaded bool

	if !Local.UseCacheOnly {
		if err := g.initFromServer(); err != nil {
			log.Error(mctx, "load config from server", "error", err)
		} else {
			loaded = true
			_ = g.saveItemsToCacheFile()
		}
	}

	if !loaded {
		log.Info(mctx, "try to use cached version")

		if err := g.initFromCacheFile(); err != nil {
			log.Fatal(mctx, "load config from cache file", "error", err)
		}
	}

	_ = g.saveItemsToCacheFile()

	if !Local.UseCacheOnly && Local.WatchEnabled {
		g.startWatch()
	}
}

/*************************************************
Description: 执行配置更新，先解析，然后回调通知对应模块
Input:
	entry	 变更key对应的配置描述条目
Output:
Return:
Others:
*************************************************/
func (g *globalConfig) onUpdate(entry *Entry, value interface{}) {
	if entry.Restart {
		log.Info(mctx, "should restart to make config update to take effect", "key", entry.Key)

		self := syscall.Getpid()
		log.Info(mctx, "send SIGUSR2 to self", "pid", self)
		if err := syscall.Kill(self, syscall.SIGUSR2); err != nil {
			log.Fatal(mctx, "send SIGUSR2 to do a graceful shutdown", "error", err)
		}
		select {}
	}

	entry.value.Store(value)

	if entry.OnUpdateFunc != nil {
		entry.OnUpdateFunc(value)
	}
}

/*************************************************
Description: 定期检测配置更新
Input:
Output:
Return:
Others:
*************************************************/
func (g *globalConfig) watchLoop() {
	defer func() {
		if r := recover(); r != nil {
			log.Error(mctx, "got panic", "error", r, "stack", utils.GetPanicStack())
		}
		log.Info(mctx, "global config watch stopped")
	}()

	for {
		select {
		case <-g.stopChan:
			return
		case <-time.After(Local.WatchInterval):
			{
				newItems, err := g.loadItemsFromServer(CONF_ROOT)
				if err != nil {
					log.Error(mctx, "reload global config", "conf_root", CONF_ROOT, "error", err)
					continue
				}

				oldItems := g.items.Load().(map[string]*dynconf.Item)

				var (
					updated        []string
					updatedEntries []*UpdatedEntry
					valid          = true
				)

				for _, item := range newItems {
					oldItem := oldItems[item.Key]
					if oldItem == nil {
						// added item
						log.Info(mctx, "global config item added", "key", item.Key, "value", item.Value, "version", item.Version)

						updated = append(updated, item.Key)
						continue
					}

					if item.Version != oldItem.Version {
						// modified item
						log.Info(mctx, "global config item modified", "key", item.Key, "value", item.Value, "version", item.Version)

						updated = append(updated, item.Key)
					}
				}

				for _, item := range oldItems {
					if newItems[item.Key] == nil {
						// deleted item, only config type is_prefix=true can be deleted
						// module use prefix config, should reload all subitems buy prefix on update.
						log.Info(mctx, "global config item deleted", "key", item.Key, "value", item.Value, "version", item.Version)

						entry := g.matchEntry(item.Key)
						if entry == nil {
							continue
						}

						if !entry.IsPrefix {
							log.Warning(mctx, "detected deletion for non-prefix config item, not allowed", "key", item.Key)

							newItems[item.Key] = item
							continue
						}

						updated = append(updated, item.Key)
					}
				}

				if len(updated) <= 0 {
					// no watched entry updated
					continue
				}

				g.items.Store(newItems)

				updatedEntries = make([]*UpdatedEntry, 0, len(updated))
				for _, updatedKey := range updated {
					entry := g.matchEntry(updatedKey)
					if entry == nil {
						// ignore entries not watched
						continue
					}

					conf, err := entry.NewFunc()
					if err != nil {
						log.Error(mctx, "parse config entry", "entry_key", entry.Key, "error", err)
						valid = false
						break
					}

					updatedEntries = append(updatedEntries, &UpdatedEntry{
						entry: entry,
						value: conf,
					})
				}

				if valid && len(updatedEntries) > 0 {
					// save current global config copy in cache file
					_ = g.saveItemsToCacheFile()

					for _, u := range updatedEntries {
						g.onUpdate(u.entry, u.value)
					}
				}

			}
		}
	}
}

/*************************************************
Description: 开启配置更新监控
Input:
Output:
Return:
Others:
*************************************************/
func (g *globalConfig) startWatch() {
	if g.watching {
		log.Fatal(mctx, "global config already watching")
	}

	g.watching = true
	go g.watchLoop()
}

/*************************************************
Description: 停止配置更新监控
Input:
Output:
Return:
Others:
*************************************************/
func (g *globalConfig) stopWatch() {
	close(g.stopChan)
}

type VersionedValue struct {
	Value   string `json:"value"`
	Version uint64 `json:"version"`
}

func (g *globalConfig) saveItemsToCacheFile() error {
	log.Info(mctx, "saving global config to cache file")

	cacheFileNameTmp := fmt.Sprint(CACHE_FILE_PATH, ".tmp")

	items := g.items.Load().(map[string]*dynconf.Item)
	cached := make(map[string]VersionedValue, len(items))

	var buf bytes.Buffer
	for _, item := range items {
		buf.Reset()
		var v VersionedValue

		v.Version = item.Version

		err := json.Compact(&buf, []byte(item.Value))
		if err != nil {
			log.Error(mctx, "compact json", "json_data", item.Value, "error", err)
			v.Value = item.Value
		} else {
			v.Value = buf.String()
		}
		cached[item.Key] = v
	}

	err := SaveConfig(cacheFileNameTmp, cached)
	if err != nil {
		log.Error(mctx, "save config item to cache file", "file", cacheFileNameTmp, "error", err)
		return err
	}

	err = os.Rename(cacheFileNameTmp, CACHE_FILE_PATH)
	if err != nil {
		log.Error(mctx, "rename cache config file", "from_file", cacheFileNameTmp, "to_file", CACHE_FILE_PATH, "error", err)
	}
	return err
}

func (g *globalConfig) loadItemsFromCacheFile() (map[string]*dynconf.Item, error) {
	var cached map[string]VersionedValue

	err := LoadConfig(CACHE_FILE_PATH, &cached)
	if err != nil {
		return nil, err
	}

	items := make(map[string]*dynconf.Item, len(cached))
	for k, v := range cached {
		items[k] = &dynconf.Item{
			Key:     k,
			Value:   v.Value,
			Version: v.Version,
		}
	}

	return items, nil
}
