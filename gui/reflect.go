package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"github.com/crossoverJie/ptg/reflect"
	"github.com/pkg/errors"
	"sync/atomic"
)

var (
	// filename->*ParseReflectAdapter
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
	containerMap = nil
	index = 0
}

func ReloadReflect(f func(uri fyne.URIReadCloser, err error)) {
	var filenameList []string
	for k := range parseContainerMap {
		filenameList = append(filenameList, k)
	}
	ClearReflect()
	for _, filename := range filenameList {
		f(&ResetUri{Filename: filename}, nil)
	}
}

func ParseContainer() map[string]*ParseReflectAdapter {
	return parseContainerMap
}

func genIndex() string {
	return fmt.Sprint(atomic.AddInt64(&index, 1))
}

func GetParseAdapter(index string) (*ParseReflectAdapter, error) {
	filename := containerMap[index]
	registerReflect, exit, _ := RegisterReflect(filename)
	if !exit {
		return nil, errors.New("proto not register")
	}
	return registerReflect, nil
}
