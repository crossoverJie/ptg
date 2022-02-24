package main

import (
	"errors"
	"github.com/cheggaaa/pb/v3"
	"github.com/crossoverJie/ptg/meta"
	"github.com/crossoverJie/ptg/model"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"
)

// init bind variable
var (
	target string
	//respCh       chan *meta.Response
	thread       int
	duration     int64
	method       string
	bodyPath     string
	body         string
	headerSlice  cli.StringSlice
	headerMap    map[string]string
	protocol     string // http/grpc
	protocolFile string // xx/xx/xx.proto
	fqn          string // fully-qualified method name:[package.Service.Method]
)

var (
	//totalRequestTime  time.Duration
	//totalResponseSize int
	//SlowRequestTime   time.Duration
	//FastRequestTime   = time.Minute
	//ErrorCount        int32
	Bar *pb.ProgressBar
)

const (
	PbTmpl = `{{ green "Requesting:" }} {{string . "target" | blue}}  {{ bar . "<" "-" (cycle . "↖" "↗" "↘" "↙" ) "." ">"}} {{speed . | rndcolor }} {{percent .}}`
	Http   = "http"
	Grpc   = "grpc"
)

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
			&cli.StringFlag{
				Name:        "Request protocol",
				Usage:       "-proto http/grpc",
				DefaultText: "http",
				Aliases:     []string{"proto"},
				Required:    true,
				Destination: &protocol,
			},
			&cli.StringFlag{
				Name:        "protocol buffer file path",
				Usage:       "-pf /file/order.proto",
				DefaultText: "",
				Aliases:     []string{"pf"},
				Required:    false,
				Destination: &protocolFile,
			},
			&cli.StringFlag{
				Name:        "fully-qualified method name",
				Usage:       "-fqn package.Service.Method",
				DefaultText: "",
				Aliases:     []string{"fqn"},
				Required:    false,
				Destination: &fqn,
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
				Name:        "HTTP method",
				Usage:       "-m GET",
				DefaultText: "GET",
				Aliases:     []string{"M"},
				Required:    false,
				Destination: &method,
			},
			&cli.StringFlag{
				Name:        "bodyPath",
				Usage:       "-body bodyPath.json",
				DefaultText: "",
				Aliases:     []string{"body"},
				Required:    false,
				Destination: &bodyPath,
			},
			&cli.StringSliceFlag{
				Name:        "header",
				Aliases:     []string{"H"},
				Usage:       "HTTP header to add to request, e.g. \"-H \"Content-Type: application/json\"\"",
				Required:    false,
				DefaultText: "",
				Destination: &headerSlice,
			},
			&cli.StringFlag{
				Name:        "target",
				Usage:       "http://gobyexample.com/grpc:127.0.0.1:5000",
				DefaultText: "",
				Aliases:     []string{"tg"},
				Required:    true,
				Destination: &target,
			},
		},
		Action: func(c *cli.Context) error {
			color.White("thread: %v, duration: %v, count %v", thread, duration, count)
			runtime.GOMAXPROCS(runtime.NumCPU() + thread)
			// ##########App init##########
			if count == 0 && duration == 0 {
				return errors.New("request count and duration must choose one")
			}

			if count > 0 && duration > 0 {
				return errors.New("request count and duration can only choose one")
			}

			if method == "" {
				method = "GET"
			}
			if bodyPath != "" {
				bytes, err := ioutil.ReadFile(bodyPath)
				if err != nil {
					color.Red("could not read file: %s", bodyPath)
					return err
				}
				body = string(bytes)
			}
			if headerSlice.Value() != nil {
				headerMap = make(map[string]string, len(headerSlice.Value()))
				for _, s := range headerSlice.Value() {
					splitN := strings.SplitN(s, ":", 2)
					headerMap[splitN[0]] = splitN[1]
				}
			}
			meta.NewResult()
			newMeta := meta.NewMeta(target, method, bodyPath, body, protocol, protocolFile, fqn, thread, duration, &headerSlice, headerMap)
			var model model.Model
			if count > 0 {
				respCh := make(chan *meta.Response, count)
				newMeta.SetRespCh(respCh)
				model = NewCountModel(count)
				Bar = pb.ProgressBarTemplate(PbTmpl).Start(count)
			} else {
				// 防止写入 goroutine 阻塞，导致泄露。
				respCh := make(chan *meta.Response, 3*thread)
				newMeta.SetRespCh(respCh)
				model = NewDurationModel(duration)
				Bar = pb.ProgressBarTemplate(PbTmpl).Start(int(duration))
			}
			Bar.Set("target", target).
				SetWidth(120)
			// ##########App init##########

			// shutdown
			signCh := make(chan os.Signal, 1)
			signal.Notify(signCh, os.Interrupt)
			go func() {
				select {
				case <-signCh:
					color.Red("shutdown....")
					model.Shutdown()
				}
			}()

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
