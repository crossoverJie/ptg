package main

import (
	"fmt"
	"github.com/docker/go-units"
	"github.com/fatih/color"
	"sync/atomic"
	"time"
)

type DurationModel struct {
	duration     int64
	timeCh       chan struct{}
	close        chan struct{}
	totalRequest int32
	shutdown     int32
}

func NewDurationModel(duration int64) Model {
	return &DurationModel{
		duration: duration,
		timeCh:   make(chan struct{}),
		close:    make(chan struct{}),
	}
}

func (c *DurationModel) Init() {
	go func() {
		for {
			select {
			case <-time.Tick(1 * time.Second):
				if atomic.LoadInt32(&c.shutdown) == 1 {
					return
				}
				Bar.Increment()
			}
		}
	}()
	go func() {
		select {
		case <-time.After(time.Second * time.Duration(c.duration)):
			close(c.close)
			return
		}
	}()

}
func (c *DurationModel) Run() {
	for i := 0; i < thread; i++ {
		go func() {
			for {
				if atomic.LoadInt32(&c.shutdown) == 1 {
					return
				}
				httpClient := NewHttpClient(method, target, body)
				response, err := httpClient.Request()
				atomic.AddInt32(&c.totalRequest, 1)
				respCh <- response
				if err != nil {
					color.Red("request err %v\n", err)
					atomic.AddInt32(&ErrorCount, 1)
				}
			}
		}()
	}
}
func (c *DurationModel) Finish() {
	for {
		select {
		case _, ok := <-c.close:
			if !ok {
				Bar.SetCurrent(duration)
				Bar.Finish()
				atomic.StoreInt32(&c.shutdown, 1)
				return
			}
		case res := <-respCh:
			if res != nil {
				totalRequestTime += res.RequestTime
				totalResponseSize += res.ResponseSize
			}
		}
	}
}

func (c *DurationModel) Shutdown() {
	close(c.close)
}

func (c *DurationModel) PrintSate() {
	fmt.Println("")
	color.Green("%v requests in %v, %v read.\n", c.totalRequest, units.HumanDuration(totalRequestTime), units.HumanSize(float64(totalResponseSize)))
	color.Green("Avg Req Time:\t\t%v\n", totalRequestTime/time.Duration(c.totalRequest))
	color.Green("Fastest Request:\t%v\n", FastRequestTime)
	color.Green("Slowest Request:\t%v\n", SlowRequestTime)
	color.Green("Number of Errors:\t%v\n", ErrorCount)
}
