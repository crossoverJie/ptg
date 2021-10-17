package main

import (
	"fmt"
	"github.com/docker/go-units"
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

func TestHuman(t *testing.T) {
	size := units.HumanSize(1024)
	fmt.Println(size)
}
