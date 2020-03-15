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
	"fmt"
	"sync"

	mapset "github.com/deckarep/golang-set"
)

type relation struct {
	referenced mapset.Set
	refers     mapset.Set
}

type threadSafeMap struct {
	lock  sync.RWMutex
	items map[string]interface{}

	// relations maps a key to an relation
	relations map[string]*relation
	referFunc ReferFunc
}

func NewThreadSafeMap(referFunc ReferFunc) *threadSafeMap {
	t := new(threadSafeMap)
	t.items = make(map[string]interface{})
	t.relations = make(map[string]*relation)
	t.referFunc = referFunc
	return t
}

func (t *threadSafeMap) Add(key string, obj interface{}) error {
	return t.Update(key, obj)
}

func (t *threadSafeMap) Update(key string, obj interface{}) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	oldObj := t.items[key]
	t.items[key] = obj
	t.updateRelation(oldObj, obj, key)
	return nil
}

func (t *threadSafeMap) Delete(key string) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	if obj, exists := t.items[key]; exists {
		t.deleteFromRelation(obj, key)
	}
	return nil
}

func (t *threadSafeMap) List() []interface{} {
	t.lock.RLock()
	defer t.lock.RUnlock()

	list := make([]interface{}, 0, len(t.items))
	for _, item := range t.items {
		list = append(list, item)
	}
	return list
}

func (t *threadSafeMap) ListKeys() []string {
	t.lock.RLock()
	defer t.lock.RUnlock()

	list := make([]string, 0, len(t.items))
	for key := range t.items {
		list = append(list, key)
	}
	return list
}

func (t *threadSafeMap) Get(key string) (item interface{}, exists bool) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	item, exists = t.items[key]
	return item, exists
}

func (t *threadSafeMap) Replace(items map[string]interface{}) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.items = items
	t.relations = make(map[string]*relation)
	return nil
}

func (t *threadSafeMap) Referenced(key string) ([]interface{}, error) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	relation, exists := t.relations[key]
	if !exists {
		return nil, fmt.Errorf("relation of key %s not found", key)
	}
	var list []interface{}
	for i := range relation.referenced.Iter() {
		key := i.(string)
		obj, exists := t.items[key]
		if !exists {
			panic(fmt.Errorf("item %s not found", key))
		}
		list = append(list, obj)
	}
	return list, nil
}

func (t *threadSafeMap) ReferencedKeys(key string) ([]string, error) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	relation, exists := t.relations[key]
	if !exists {
		return nil, fmt.Errorf("relation of key %s not found", key)
	}
	var list []string
	for i := range relation.referenced.Iter() {
		key := i.(string)
		list = append(list, key)
	}
	return list, nil

}

func (t *threadSafeMap) ReferKeys(key string) ([]string, error) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	relation, exists := t.relations[key]
	if !exists {
		return nil, fmt.Errorf("relation of key %s no found", key)
	}
	if relation.refers == nil {
		return nil, nil
	}
	var list []string
	for refer := range relation.refers.Iter() {
		r := refer.(string)
		list = append(list, r)
	}
	return list, nil
}

func (t *threadSafeMap) updateRelation(oldObj interface{}, newObj interface{}, key string) {
	if oldObj != nil {
		t.deleteFromRelation(oldObj, key)
	}
	refers, err := t.referFunc(newObj)
	if err != nil {
		panic(fmt.Errorf("unable to calculate refers for key %q: %v", key, err))
	}
	curRelation, exist := t.relations[key]
	if !exist {
		curRelation = new(relation)
		t.relations[key] = curRelation
	}
	if len(refers) == 0 {
		return
	}
	if curRelation.refers == nil {
		curRelation.refers = mapset.NewThreadUnsafeSet()
	}
	curRelation.refers.Clear()
	for _, refKey := range refers {
		curRelation.refers.Add(refKey)
		refRelation, exist := t.relations[refKey]
		if !exist {
			refRelation = new(relation)
			t.relations[refKey] = refRelation
		}
		if refRelation.referenced == nil {
			refRelation.referenced = mapset.NewThreadUnsafeSet()
		}
		refRelation.referenced.Add(key)
	}
}

func (t *threadSafeMap) deleteFromRelation(obj interface{}, key string) {
	refers, err := t.referFunc(obj)
	if err != nil {
		panic(fmt.Errorf("unable to calculate refers for key %q: %v", key, err))
	}
	for _, refKey := range refers {
		relat, exist := t.relations[refKey]
		if !exist {
			continue
		}
		if relat.referenced == nil {
			continue
		}
		relat.referenced.Remove(key)
	}
	delete(t.relations, key)
}
