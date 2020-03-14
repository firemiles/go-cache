package relation

import (
	"github.com/firemiles/go-cache/pkg/types"
)

// ReferFunc knows how to get refers from an object. Implementations should be deterministic.
type ReferFunc func(obj interface{}) ([]string, error)

// Cache is a week reference relationship trace memory cache. ,
// you can delete one object which referenced by other without any error.
type Cache interface {
	types.Store
	Referenced(object interface{}) ([]interface{}, error)
	ReferencedKeys(key string) ([]string, error)
	ReferKeys(key string) ([]string, error)
}

type cache struct {
	cacheStorage *threadSafeMap
	keyFunc types.KeyFunc
}

var _ Cache = &cache{}

func NewCache(keyFunc types.KeyFunc, referFunc ReferFunc) Cache {
	c := new(cache)
	c.cacheStorage = NewThreadSafeMap(referFunc)
	c.keyFunc = keyFunc
	return c
}

func (c *cache) Add(obj interface{}) error {
	key, err := c.keyFunc(obj)
	if err != nil {
		return types.KeyError{obj, err}
	}
	return c.cacheStorage.Add(key, obj)
}

func (c *cache) Update(obj interface{}) error {
	key, err := c.keyFunc(obj)
	if err != nil {
		return types.KeyError{obj, err}
	}
	return c.cacheStorage.Update(key, obj)
}

func (c *cache) Delete(obj interface{}) error {
	key, err := c.keyFunc(obj)
	if err != nil {
		return types.KeyError{obj, err}
	}
	return c.cacheStorage.Delete(key)
}

func (c *cache) List() []interface{} {
	return c.cacheStorage.List()
}

func (c *cache) ListKeys() []string {
	return c.cacheStorage.ListKeys()
}

func (c *cache) Get(obj interface{}) (item interface{}, exists bool, err error) {
	key, err := c.keyFunc(obj)
	if err != nil {
		return nil, false, types.KeyError{obj, err}
	}
	return c.GetByKey(key)
}

func (c *cache) GetByKey(key string) (item interface{}, exists bool, err error) {
	item, exists = c.cacheStorage.Get(key)
	return item, exists, nil
}

func (c *cache) Replace(list []interface{}) error {
	items := make(map[string]interface{}, len(list))
	for _, item := range list {
		key, err := c.keyFunc(item)
		if err != nil {
			return types.KeyError{item, err}
		}
		items[key] = item
	}
	return c.cacheStorage.Replace(items)
}

func (c *cache) Referenced(obj interface{}) ([]interface{}, error) {
	key, err := c.keyFunc(obj)
	if err != nil {
		return nil, types.KeyError{obj, err}
	}
	return c.cacheStorage.Referenced(key)
}

func (c *cache) ReferencedKeys(key string) ([]string, error) {
	return c.cacheStorage.ReferencedKeys(key)
}

func (c *cache) ReferKeys(key string) ([]string, error) {
	return c.cacheStorage.ReferKeys(key)
}





