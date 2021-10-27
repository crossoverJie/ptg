# ptg
Performance testing tool (Go)

![](pic/show.gif)

# Building

```go
go get github.com/crossoverJie/ptg
```

# Usage

```shell script
NAME:
   ptg - Performance testing tool (Go)

USAGE:
   ___go_build_github_com_crossoverJie_ptg [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --thread value, -t value        -t 10 (default: 1 thread)
   --duration value, -d value      -d 10s (default: Duration of test in seconds, Default 10s)
   --request value, -c value       -c 100 (default: 100)
   --HTTP value, -M value          -m GET (default: GET)
   --bodyPath value, --body value  -bodyPath bodyPath.json
   --header value, -H value        HTTP header to add to request, e.g. "-H Content-Type: application/json"
   --target value, --tg value      http://gobyexample.com (default: http://gobyexample.com)
   --help, -h                      show help (default: false)
```

```shell script
btb -t 20 -d 10  -tg "http://gobyexample.com"
```
Benchmark test for 10 seconds, using 20 goroutines.

output:
```shell script
Requesting: http://gobyexample.com  <---------------> 1 p/s 100.00%

43 requests in 10 seconds, 13.88MB read.
Avg Req Time:           358.512071ms
Fastest Request:        93.518704ms
Slowest Request:        840.680771ms
Number of Errors:       0
```