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
	err = SaveLog(marshal)
	if err != nil {
		panic(err)
	}
}

func TestLoadLog(t *testing.T) {
	bytes, err := LoadLog()
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
