package main

import (
	"fmt"
	"testing"
)

func Test_client_Request(t *testing.T) {
	httpClient := NewHttpClient("GET", "http://gobyexample.com", "")
	request, err := httpClient.Request()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(request)
}
