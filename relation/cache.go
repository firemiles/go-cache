/*
 * Copyright (c) 2020 firemiles(miles.dev@outlook.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

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
	keyFunc      types.KeyFunc
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
