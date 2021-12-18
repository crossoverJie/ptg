package io

import (
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"os"
	"os/user"
)

const (
	AppLog       = "log"
	AppSearchLog = "search_log"
	Path         = "path"
	FileName     = "filename"
)

var (
	LogMeta = map[string]map[string]string{
		AppLog: {
			Path:     "/.ptg/",
			FileName: "/.ptg/ptg.log",
		},
		AppSearchLog: {
			Path:     "/.ptg/",
			FileName: "/.ptg/search.log",
		},
	}
)

func SaveLog(logType string, data []byte) error {
	filename, err := initLog(logType)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 0666)
}

func initLog(logType string) (string, error) {
	home, err := user.Current()
	if err != nil {
		return "", err
	}
	m := LogMeta[logType]
	filename := home.HomeDir + m[FileName]

	if !exist(filename) {
		err := os.MkdirAll(home.HomeDir+m[Path], 0777)
		if err != nil {
			return "", err
		}
		_, err = os.Create(filename)
		if err != nil {
			return "", err
		}
	}
	return filename, nil
}

func LoadLog(logType string) ([]byte, error) {
	filename, err := initLog(logType)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadFile(filename)
}

func LoadLogWithStruct() (*Log, error) {
	bytes, err := LoadLog(AppLog)
	var read Log
	err = proto.Unmarshal(bytes, &read)
	if err != nil {
		return nil, err
	}
	return &read, nil
}
func LoadSearchLogWithStruct() (*SearchLogList, error) {
	bytes, err := LoadLog(AppSearchLog)
	var read SearchLogList
	err = proto.Unmarshal(bytes, &read)
	if err != nil {
		return nil, err
	}
	return &read, nil
}

func exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
