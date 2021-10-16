package main

import (
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"sync"
	"time"
)

var thread int
var count int
var duration int64
var target string
var wait sync.WaitGroup
var workCh chan *Job
var respCh chan *Response

var totalRequestTime time.Duration
var totalResponseSize int64

func main() {
	app := &cli.App{Name: "ptg", Description: "Performance testing tool (Go)",
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
			wait.Add(count)
			workCh = make(chan *Job, count)
			respCh = make(chan *Response, count)
			addJob(&Job{
				thread:   thread,
				duration: duration,
				count:    count,
				target:   target,
			})
			execJob(thread)
			wait.Wait()
			finishJob()
			printState()

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func addJob(job *Job) {

	for i := 0; i < job.count; i++ {
		go func() {
			workCh <- job
		}()
	}

}

func execJob(thread int) {
	for i := 0; i < thread; i++ {
		go func() {
			for {
				select {
				case job := <-workCh:
					httpClient := NewHttpClient("GET", job.target, "")
					response, err := httpClient.Request()
					respCh <- response
					wait.Done()
					if err != nil {
						color.Red("request err", err)
						return
					}

				}
			}

		}()
	}
}

func finishJob() {
	for i := 0; i < count; i++ {
		select {
		case res := <-respCh:
			if res != nil {
				totalRequestTime += res.RequestTime
				totalResponseSize += res.ResponseSize
			}
		}
	}

}

func printState() {
	color.Green("Avg Req Time:\t\t%v\n", totalRequestTime/time.Duration(count))
}

type Job struct {
	thread   int
	duration int64
	count    int
	target   string
}
