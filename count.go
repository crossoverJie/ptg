package main

import (
	ptgclient "github.com/crossoverJie/ptg/client"
	"github.com/crossoverJie/ptg/meta"
	"github.com/crossoverJie/ptg/model"
	"github.com/docker/go-units"
	"github.com/fatih/color"
	"os"
	"sync"
	"time"
)

type CountModel struct {
	wait   sync.WaitGroup
	count  int
	workCh chan struct{}
	start  time.Time
	result *meta.Result
	meta   *meta.Meta
}

func NewCountModel(count int) model.Model {
	return &CountModel{count: count, start: time.Now(), result: meta.GetResult(), meta: meta.GetMeta()}
}

func (c *CountModel) Init() {
	c.wait.Add(c.count)
	c.workCh = make(chan struct{}, c.count)
	for i := 0; i < c.count; i++ {
		go func() {
			c.workCh <- struct{}{}
		}()
	}
}
func (c *CountModel) Run() {
	for i := 0; i < thread; i++ {
		httpClient := ptgclient.NewClient(method, target, body, c.meta)
		go func() {
			for {
				select {
				case _, ok := <-c.workCh:
					if !ok {
						return
					}
					response, err := httpClient.Request()
					c.meta.RespCh() <- response
					if err != nil {
						color.Red("request err %v\n", err)
						//atomic.AddInt32(&ErrorCount, 1)
						c.result.IncrementErrorCount()
					}
					Bar.Increment()
					c.wait.Done()
				}
			}

		}()
	}
	c.wait.Wait()
}
func (c *CountModel) Finish() {
	for i := 0; i < c.count; i++ {
		select {
		case res := <-c.meta.RespCh():
			if res != nil {
				meta.GetResult().SetTotalRequestTime(res.RequestTime).
					SetTotalResponseSize(res.ResponseSize)
			}
		}
	}
	Bar.Finish()
	close(c.workCh)
}

func (c *CountModel) Shutdown() {
	close(c.workCh)
	os.Exit(-1)
}

func (c *CountModel) PrintSate() {
	color.Green("%v requests in %v, %v read, and cost %v.\n", c.count, units.HumanDuration(c.result.TotalRequestTime()), units.HumanSize(float64(c.result.TotalResponseSize())), units.HumanDuration(time.Since(c.start)))
	color.Green("Requests/sec:\t\t%.2f\n", float64(c.count)/time.Since(c.start).Seconds())
	color.Green("Avg Req Time:\t\t%v\n", c.result.TotalRequestTime()/time.Duration(c.count))
	color.Green("Fastest Request:\t%v\n", c.result.FastRequestTime())
	color.Green("Slowest Request:\t%v\n", c.result.SlowRequestTime())
	color.Green("Number of Errors:\t%v\n", c.result.ErrorCount())
}
