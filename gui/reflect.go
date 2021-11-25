package main

import (
	"fmt"
	"github.com/crossoverJie/ptg/reflect"
	"sync/atomic"
)

var (
	parseContainerMap map[string]*ParseReflectAdapter
	// index->filename
	containerMap map[string]string
	index        int64
)

type ParseReflectAdapter struct {
	parse *reflect.ParseReflect
	index string
}

func (p *ParseReflectAdapter) Parse() *reflect.ParseReflect {
	return p.parse
}
func (p *ParseReflectAdapter) Index() string {
	return p.index
}

func RegisterReflect(filename string) (*ParseReflectAdapter, bool, error) {
	parseAdapter, ok := parseContainerMap[filename]
	if ok {
		return parseAdapter, true, nil
	}
	newParse, err := reflect.NewParse(filename)
	if err != nil {
		return nil, false, err
	}
	if parseContainerMap == nil {
		parseContainerMap = make(map[string]*ParseReflectAdapter)
	}
	if containerMap == nil {
		containerMap = make(map[string]string)
	}

	index := genIndex()
	containerMap[index] = filename
	parseAdapter = &ParseReflectAdapter{
		parse: newParse,
		index: index,
	}
	parseContainerMap[filename] = parseAdapter

	return parseAdapter, false, nil
}

func ClearReflect() {
	parseContainerMap = nil
}

func ResetReflect() {
	var filenameList []string
	for k := range parseContainerMap {
		filenameList = append(filenameList, k)
	}
	ClearReflect()
	for _, filename := range filenameList {
		RegisterReflect(filename)
	}
}

func ParseContainer() map[string]*ParseReflectAdapter {
	return parseContainerMap
}

func genIndex() string {
	return fmt.Sprint(atomic.AddInt64(&index, 1))
}

func GetParseAdapter(index string) *ParseReflectAdapter {
	filename := containerMap[index]
	registerReflect, _, _ := RegisterReflect(filename)
	return registerReflect
}
