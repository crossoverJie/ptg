package io

import (
	"fmt"
	"github.com/crossoverJie/ptg/reflect/gen/user"
	"github.com/golang/protobuf/proto"
	"testing"
)

func TestSaveLog(t *testing.T) {
	creat := user.UserApiCreate{UserId: 100}
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

	var read user.UserApiCreate
	err = proto.Unmarshal(bytes, &read)
	if err != nil {
		panic(err)
	}
	fmt.Println(read.UserId)
}
