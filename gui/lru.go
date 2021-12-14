package main

import (
	"container/list"
	"fmt"
	"sync"
)

type LruCache struct {
	size     int
	values   *list.List
	cacheMap map[interface{}]*list.Element
	lock     sync.Mutex
}

func NewLruList(size int) *LruCache {
	values := list.New()

	return &LruCache{
		size:     size,
		values:   values,
		cacheMap: make(map[interface{}]*list.Element, size),
	}
}

func (l *LruCache) Put(k, v interface{}) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.values.Len() == l.size {
		back := l.values.Back()
		l.values.Remove(back)
		delete(l.cacheMap, k)
	}

	front := l.values.PushFront(v)
	l.cacheMap[k] = front
}

func (l *LruCache) Get(k interface{}) (interface{}, bool) {
	v, ok := l.cacheMap[k]
	if ok {
		l.values.MoveToFront(v)
		return v, true
	} else {
		return nil, false
	}
}

func (l *LruCache) Size() int {
	return l.values.Len()
}
func (l *LruCache) String() {
	for i := l.values.Front(); i != nil; i = i.Next() {
		fmt.Print(i.Value, "\t")
	}
}
func (l *LruCache) List() []interface{} {
	var data []interface{}
	for i := l.values.Front(); i != nil; i = i.Next() {
		data = append(data, i.Value)
	}
	return data
}

func (l *LruCache) Clear() {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.values = list.New()
	l.cacheMap = make(map[interface{}]*list.Element, l.size)

}
