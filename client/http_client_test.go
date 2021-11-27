package client

import (
	"fmt"
	"github.com/crossoverJie/ptg/meta"
	"github.com/docker/go-units"
	"testing"
)

func Test_httpClient_Request(t *testing.T) {
	meta.NewResult()
	httpClient := NewClient("GET",
		"http://gobyexample.com", "",
		meta.NewMeta("http://gobyexample.com", "GET", "", "", Http, "",
			"", 1, 1, nil, nil))
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
