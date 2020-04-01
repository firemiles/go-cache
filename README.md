# go-cache

[![Build Status](https://travis-ci.org/firemiles/go-cache.svg?branch=master)](https://travis-ci.org/firemiles/go-cache)
![GitHub](https://img.shields.io/github/license/firemiles/go-cache)

## Install

```sh
go get -u github.com/firemiles/go-cache
```
## Relation Cache

Relation Cache is a memory cache trace object reference relationship.

### Usage

```go
import "github.com/firemiles/go-cache/relation"

type Object struct {
    ID string
    SubObject []*Object
}

func ObjectKey(obj interface{}) (string, error) {
    o, ok := obj.(*Object)
    if !ok {
        return "", fmt.Error("only support type *Object")
    }
    return o.ID, nil
}

func ObjectRefers(obj interface{}) ([]string, error) {
    o, ok := obj.(*Object)
    if !ok {
        return nil, fmt.Error("only support type *Object")
    }
    var list []string
    for _, sub := range o.SubObject {
        list = append(list, sub.ID)
    }
    return list, nil
}

cache := relation.NewCache(ObjectKey, ObjectRefers)
subObj1 := &Object {
    ID: "sub_object1"
}
obj1 := &Object {
    ID: "object1",
    SubObject: []*Object{subObj1}
}

cache.Add(obj1)
cache.Add(subObj1)
reference, __ := cache.References(subObj1)
fmt.Printf("reference=%v", reference)

// reference = []interface{obj1}
```
