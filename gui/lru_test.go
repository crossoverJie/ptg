package main

import (
	"container/list"
	"fmt"
	"testing"
)

func TestLruCache_Add(t *testing.T) {
	l := list.New()
	l.PushBack(1)
	back := l.PushBack(2)
	l.PushBack(3)
	len := l.Len()

	l.MoveToFront(back)

	for i := 0; i < len; i++ {
		front := l.Front()
		fmt.Println(front)
		l.Remove(front)
	}

}

func TestLruCache_Put(t *testing.T) {
	l := NewLruList(3)
	l.Put(1, 1)
	l.Put(2, 2)
	l.Put(3, 3)
	l.String()
	fmt.Println("=====")
	l.Put(4, 4)
	l.String()
	fmt.Println("=====")
	l.Get(3)
	l.String()
}
func TestLruCache_Put2(t *testing.T) {
	l := NewLruList(3)
	l.Put(1, 1)
	l.Put(2, 2)
	l.Put(3, 3)
	i := l.List()
	fmt.Println(i)
	fmt.Println("=====")
	l.Put(4, 4)
	i = l.List()
	fmt.Println(i)
	fmt.Println("=====")
	l.Get(3)
	i = l.List()
	fmt.Println(i)

}
