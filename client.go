package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type (
	IClient interface {
		Request() (*Response, error)
	}

	client struct {
		Method      string
		Url         string
		RequestBody string
	}

	Response struct {
		RequestTime  time.Duration
		ResponseSize int64
	}
)

func NewHttpClient(method, url, requestBody string) IClient {
	return &client{
		Method:      method,
		Url:         url,
		RequestBody: requestBody,
	}
}

func NewGrpcClient() IClient {
	return nil
}

func (c *client) Request() (*Response, error) {

	var buf io.Reader
	if len(c.RequestBody) > 0 {
		buf = bytes.NewBufferString(c.RequestBody)
	}
	req, err := http.NewRequest(c.Method, c.Url, buf)
	if err != nil {
		fmt.Println("An error occured doing request", err)
		return nil, err
	}
	req.Header.Add("User-Agent", "ptg")

	httpClient := &http.Client{
		Transport: &http.Transport{
			ResponseHeaderTimeout: time.Millisecond * time.Duration(1000),
		},
	}

	start := time.Now()
	response, err := httpClient.Do(req)
	r := &Response{
		RequestTime:  time.Since(start),
		ResponseSize: 0,
	}
	if err != nil {
		fmt.Println("Request err:", err)
		return nil, err
	}

	if response == nil {
		return nil, errors.New("response is nil")
	}
	defer func() {
		if response != nil && response.Body != nil {
			_ = response.Body.Close()
		}
	}()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("http code not OK: %v", response.StatusCode))
	}
	return r, nil
}
