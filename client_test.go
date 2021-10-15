package main

import (
	"fmt"
	"testing"
)

func Test_client_Request(t *testing.T) {
	httpClient := NewHttpClient("GET", "http://gobyexample.com", "")
	err := httpClient.Request()
	if err != nil {
		fmt.Println(err)
	}
}
