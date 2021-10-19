package main

import (
	"errors"
	"github.com/cheggaaa/pb/v3"
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

var Bar *pb.ProgressBar

const PbTmpl = `{{ green "Requesting:" }} {{string . "target" | blue}}  {{ bar . "<" "-" (cycle . "↖" "↗" "↘" "↙" ) "." ">"}} {{speed . | rndcolor }} {{percent .}}`

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
		start  time.Time
	}
)

func NewCountModel(count int) Model {
	return &countModel{count: count, start: time.Now()}
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
	Bar.Finish()
}
func (c *countModel) PrintSate() {
	color.Green("%v requests in %v, %v read, and cost %v.\n", c.count, units.HumanDuration(totalRequestTime), units.HumanSize(float64(totalResponseSize)), units.HumanDuration(time.Since(c.start)))
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
			if count == 0 && duration == 0 {
				return errors.New("request count and duration must choose one")
			}

			if count > 0 && duration > 0 {
				return errors.New("request count and duration can only choose one")
			}
			respCh = make(chan *Response, count)
			var model Model
			if count > 0 {
				model = NewCountModel(count)
				Bar = pb.ProgressBarTemplate(PbTmpl).Start(count)
			} else {
				model = NewDurationModel(duration)
				Bar = pb.ProgressBarTemplate(PbTmpl).Start(int(duration))
			}
			Bar.Set("my_green_string", "green").Set("my_blue_string", "blue")
			Bar.Set("target", target).
				SetWidth(120)

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
