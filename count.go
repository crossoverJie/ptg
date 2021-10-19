package main

import (
	"github.com/docker/go-units"
	"github.com/fatih/color"
	"sync"
	"sync/atomic"
	"time"
)

type CountModel struct {
	wait   sync.WaitGroup
	count  int
	workCh chan *Job
	start  time.Time
}

func NewCountModel(count int) Model {
	return &CountModel{count: count, start: time.Now()}
}

func (c *CountModel) Init() {
	c.wait.Add(c.count)
	c.workCh = make(chan *Job, c.count)
	for i := 0; i < c.count; i++ {
		go func() {
			c.workCh <- &Job{
				thread:   thread,
				duration: duration,
				count:    c.count,
				target:   target,
			}
		}()
	}
}
func (c *CountModel) Run() {
	for i := 0; i < thread; i++ {
		go func() {
			for {
				select {
				case job := <-c.workCh:
					httpClient := NewHttpClient("GET", job.target, "")
					response, err := httpClient.Request()
					respCh <- response
					if err != nil {
						color.Red("request err %v\n", err)
						atomic.AddInt32(&ErrorCount, 1)
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
		case res := <-respCh:
			if res != nil {
				totalRequestTime += res.RequestTime
				totalResponseSize += res.ResponseSize
			}
		}
	}
	Bar.Finish()
}
func (c *CountModel) PrintSate() {
	color.Green("%v requests in %v, %v read, and cost %v.\n", c.count, units.HumanDuration(totalRequestTime), units.HumanSize(float64(totalResponseSize)), units.HumanDuration(time.Since(c.start)))
	color.Green("Avg Req Time:\t\t%v\n", totalRequestTime/time.Duration(c.count))
	color.Green("Fastest Request:\t%v\n", FastRequestTime)
	color.Green("Slowest Request:\t%v\n", SlowRequestTime)
	color.Green("Number of Errors:\t%v\n", ErrorCount)
}
