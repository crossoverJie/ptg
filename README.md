# ptg
Performance testing tool (Go)


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
btb -t 10 -d 30s -c 100 -tg "http://gobyexample.com"
```