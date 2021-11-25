package main

import (
	"github.com/crossoverJie/ptg/reflect"
	"sync/atomic"
)

var (
	parseContainer map[string]*reflect.ParseReflect
	containerMap   map[int64]string
	index          int64
)

type ParseReflectAdapter struct {
	parse *reflect.ParseReflect
	index int64
}

func RegisterReflect(filename string) (*reflect.ParseReflect, bool, error) {
	parse, ok := parseContainer[filename]
	if ok {
		return parse, true, nil
	}
	parse, err := reflect.NewParse(filename)
	if err != nil {
		return nil, false, err
	}
	if parseContainer == nil {
		parseContainer = make(map[string]*reflect.ParseReflect)
	}
	parseContainer[filename] = parse

	return parse, false, nil
}

func ClearReflect() {
	parseContainer = nil
}

func ResetReflect() {
	var filenameList []string
	for k := range parseContainer {
		filenameList = append(filenameList, k)
	}
	ClearReflect()
	for _, filename := range filenameList {
		RegisterReflect(filename)
	}
}

func ParseContainer() map[string]*reflect.ParseReflect {
	return parseContainer
}

func genIndex() int64 {
	return atomic.AddInt64(&index, 1)
}
