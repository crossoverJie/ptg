package io

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"testing"
)

func TestSaveLog(t *testing.T) {
	creat := Log{
		Filenames: []string{"test.proto", "user.proto"},
		Target:    "127.0.0.1:6001",
		Request:   `{"order_id":1123120,"reason_id":null,"remark":"","user_id":null}`,
		Metadata:  `{"lang":"zh"}`,
		Response:  `{"orderId":"1123120"}`,
	}
	marshal, err := proto.Marshal(&creat)
	if err != nil {
		panic(err)
	}
	err = SaveLog(AppLog, marshal)
	if err != nil {
		panic(err)
	}
}

func TestLoadLog(t *testing.T) {
	bytes, err := LoadLog(AppLog)
	if err != nil {
		panic(err)
	}

	var read Log
	err = proto.Unmarshal(bytes, &read)
	if err != nil {
		panic(err)
	}
	fmt.Println(read)
}

func TestLoadLogWithStruct(t *testing.T) {
	withStruct, err := LoadLogWithStruct()
	if err != nil {
		panic(err)
	}
	fmt.Println(withStruct)
}

func TestSaveLog1(t *testing.T) {
	searchLogList := &SearchLogList{}
	searchLogList.SearchLogList = append(searchLogList.SearchLogList, &SearchLog{
		Id: 1,
		Value: &Log{
			Target:   "1212",
			Request:  "4334",
			Metadata: "434",
			Response: "434",
		},
		MethodInfo: "23232",
	})
	marshal, err := proto.Marshal(searchLogList)
	if err != nil {
		panic(err)
	}
	err = SaveLog(AppSearchLog, marshal)
	if err != nil {
		panic(err)
	}
}
