package main

import (
	"fmt"
	ptgclient "github.com/crossoverJie/ptg/client"
	"github.com/crossoverJie/ptg/meta"
	"github.com/crossoverJie/ptg/model"
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
	result       *meta.Result
	meta         *meta.Meta
}

func NewDurationModel(duration int64) model.Model {
	return &DurationModel{
		duration: duration,
		timeCh:   make(chan struct{}),
		close:    make(chan struct{}),
		result:   meta.GetResult(),
		meta:     meta.GetMeta(),
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
		ptgClient := ptgclient.NewClient(method, target, body, c.meta)
		go func() {
			for {
				if atomic.LoadInt32(&c.shutdown) == 1 {
					return
				}
				response, err := ptgClient.Request()
				atomic.AddInt32(&c.totalRequest, 1)
				c.meta.RespCh() <- response
				if err != nil {
					color.Red("request err %v\n", err)
					c.result.IncrementErrorCount()
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
		case res := <-c.meta.RespCh():
			if res != nil {
				c.result.SetTotalRequestTime(res.RequestTime).
					SetTotalResponseSize(res.ResponseSize)
			}
		}
	}
}

func (c *DurationModel) Shutdown() {
	close(c.close)
}

func (c *DurationModel) PrintSate() {
	fmt.Println("")
	color.Green("%v requests in %v, %v read.\n", c.totalRequest, units.HumanDuration(c.result.TotalRequestTime()), units.HumanSize(float64(c.result.TotalResponseSize())))
	color.Green("Requests/sec:\t\t%.2f\n", float64(c.totalRequest)/float64(duration))
	color.Green("Avg Req Time:\t\t%v\n", c.result.TotalRequestTime()/time.Duration(c.totalRequest))
	color.Green("Fastest Request:\t%v\n", c.result.FastRequestTime())
	color.Green("Slowest Request:\t%v\n", c.result.SlowRequestTime())
	color.Green("Number of Errors:\t%v\n", c.result.ErrorCount())
}
