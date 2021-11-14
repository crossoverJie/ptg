package main

import (
	"fmt"
	"github.com/docker/go-units"
	"testing"
)

func Test_client_Request(t *testing.T) {
	httpClient := NewClient("GET", "http://gobyexample.com", "")
	request, err := httpClient.Request()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(request)
}
func Test_client_Request2(t *testing.T) {
	httpClient := NewClient("POST", "http://localhost:8080/post", `{"name":"abc"}`)
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

func TestBytes(t *testing.T) {
	body := []byte{3, 14, 159, 2, 65, 35, 9}
	fmt.Println(string(body))

}
