package main

import (
	"encoding/json"
	"fmt"
	"github.com/crossoverJie/ptg/gui/io"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHistory_Search(t *testing.T) {
	history := NewHistory(10, nil, nil, nil, nil, nil)
	history.Put(1, &HistoryValue{
		Id: 1,
		Value: &io.Log{
			Target:   "127.0.0.1:6001",
			Request:  "{\"order_id\":0,\"reason_id\":null,\"remark\":\"\",\"user_id\":null}",
			Metadata: "{\"name\":\"abc\"}",
			Response: "",
		},
		MethodInfo: "order.v1.OrderService.Create",
	})
	history.Put(2, &HistoryValue{
		Id: 2,
		Value: &io.Log{
			Target:   "127.0.0.1:6002",
			Request:  "{\"order_id\":99999,\"reason_id\":null,\"remark\":\"\",\"user_id\":null}",
			Metadata: "{\"name\":\"zhangsan\"}",
			Response: "",
		},
		MethodInfo: "order.v1.OrderService.Close",
	})
	history.Put(3, &HistoryValue{
		Id: 3,
		Value: &io.Log{
			Target:   "127.0.0.1:6003",
			Request:  "{\"order_id\":99999,\"reason_id\":null,\"remark\":\"\",\"user_id\":null}",
			Metadata: "{\"name\":\"zhangsan\"}",
			Response: "",
		},
		MethodInfo: "order.v1.OrderService.List",
	})

	search := history.SearchResult("6001")
	for _, value := range search {
		marshal, _ := json.Marshal(value)
		fmt.Println(string(marshal))
		assert.Equal(t, value.Id, 1)
	}
	search = history.SearchResult("9999")
	for _, value := range search {
		marshal, _ := json.Marshal(value)
		fmt.Println(string(marshal))
	}
}
