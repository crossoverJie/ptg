package main

import (
	"github.com/docker/go-units"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

var target string
var respCh chan *Response
var thread int
var duration int64

var totalRequestTime time.Duration
var totalResponseSize int
var SlowRequestTime time.Duration
var FastRequestTime = time.Minute
var ErrorCount int32

type (
	Model interface {
		Init()
		Run()
		Finish()
		PrintSate()
	}

	Job struct {
		thread   int
		duration int64
		count    int
		target   string
	}

	countModel struct {
		wait   sync.WaitGroup
		count  int
		workCh chan *Job
	}
)

func NewCountModel(count int) Model {
	return &countModel{count: count}
}

func (c *countModel) Init() {
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
func (c *countModel) Run() {
	for i := 0; i < thread; i++ {
		go func() {
			for {
				select {
				case job := <-c.workCh:
					httpClient := NewHttpClient("GET", job.target, "")
					response, err := httpClient.Request()
					c.wait.Done()
					respCh <- response
					if err != nil {
						atomic.AddInt32(&ErrorCount, 1)
						color.Red("request err %v\n", err)
						continue
					}
				}
			}

		}()
	}
	c.wait.Wait()
}
func (c *countModel) Finish() {
	for i := 0; i < c.count; i++ {
		select {
		case res := <-respCh:
			if res != nil {
				totalRequestTime += res.RequestTime
				totalResponseSize += res.ResponseSize
			}
		}
	}
}
func (c *countModel) PrintSate() {
	color.Green("%v requests in %v, %v read\n", c.count, units.HumanDuration(totalRequestTime), units.HumanSize(float64(totalResponseSize)))
	color.Green("Avg Req Time:\t\t%v\n", totalRequestTime/time.Duration(c.count))
	color.Green("Fastest Request:\t%v\n", FastRequestTime)
	color.Green("Slowest Request:\t%v\n", SlowRequestTime)
	color.Green("Number of Errors:\t%v\n", ErrorCount)
}

func main() {
	var count int
	app := &cli.App{Name: "ptg", Usage: "Performance testing tool (Go)",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "thread",
				Usage: "-t 10",
				//Value:       1000,
				DefaultText: "1 thread",
				Aliases:     []string{"t"},
				Required:    true,
				Destination: &thread,
			},
			&cli.Int64Flag{
				Name:        "duration",
				Usage:       "-d 10s",
				DefaultText: "Duration of test in seconds, Default 10s",
				Aliases:     []string{"d"},
				Required:    false,
				Destination: &duration,
			},
			&cli.IntFlag{
				Name:        "request count",
				Usage:       "-c 100",
				DefaultText: "100",
				Aliases:     []string{"c"},
				Required:    false,
				Destination: &count,
			},
			&cli.StringFlag{
				Name:        "target",
				Usage:       "http://gobyexample.com",
				DefaultText: "http://gobyexample.com",
				Aliases:     []string{"tg"},
				Required:    true,
				Destination: &target,
			},
		},
		Action: func(c *cli.Context) error {
			color.White("thread: %v, duration: %v, count %v", thread, duration, count)
			respCh = make(chan *Response, count)
			model := NewCountModel(count)
			model.Init()
			model.Run()
			model.Finish()
			model.PrintSate()

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
