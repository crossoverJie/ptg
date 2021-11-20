package meta

import (
	"github.com/urfave/cli/v2"
	"sync/atomic"
	"time"
)

type Response struct {
	RequestTime  time.Duration
	ResponseSize int
}

func (r *Response) FastRequest() time.Duration {
	if r.RequestTime < GetResult().FastRequestTime() {
		return r.RequestTime
	}
	return GetResult().FastRequestTime()
}
func (r *Response) SlowRequest() time.Duration {
	if r.RequestTime > GetResult().SlowRequestTime() {
		return r.RequestTime
	}
	return GetResult().SlowRequestTime()
}

var result *Result

type Result struct {
	totalRequestTime  time.Duration
	totalResponseSize int
	slowRequestTime   time.Duration
	fastRequestTime   time.Duration
	errorCount        int32
}

func GetResult() *Result {
	return result
}

func NewResult() *Result {
	if result != nil {
		return result
	}
	result = &Result{fastRequestTime: time.Minute}
	return GetResult()
}

func (m *Result) SetTotalRequestTime(req time.Duration) *Result {
	m.totalRequestTime += req
	return m
}

func (m *Result) TotalRequestTime() time.Duration {
	return m.totalRequestTime
}

func (m *Result) SetTotalResponseSize(req int) *Result {
	m.totalResponseSize += req
	return m
}

func (m *Result) TotalResponseSize() int {
	return m.totalResponseSize
}

func (m *Result) SetSlowRequestTime(req time.Duration) *Result {
	GetResult().slowRequestTime = req
	return m
}
func (m *Result) SlowRequestTime() time.Duration {
	return m.slowRequestTime
}

func (m *Result) SetFastRequestTime(req time.Duration) *Result {
	GetResult().fastRequestTime = req
	return m
}
func (m *Result) FastRequestTime() time.Duration {
	return m.fastRequestTime
}

func (m *Result) IncrementErrorCount() *Result {
	atomic.AddInt32(&m.errorCount, 1)
	return m
}

func (m *Result) ErrorCount() int32 {
	return m.errorCount
}

type Meta struct {
	target       string
	respCh       chan *Response
	thread       int
	duration     int64
	method       string
	bodyPath     string
	body         string
	headerSlice  *cli.StringSlice
	headerMap    map[string]string
	protocol     string // http/grpc
	protocolFile string // xx/xx/xx.proto
	fqn          string // fully-qualified method name:[package.Service.Method]
}

var meta *Meta

func NewMeta(target, method, bodyPath, body, protocol, protocolFile, fqn string, thread int, duration int64,
	headerSlice *cli.StringSlice, headerMap map[string]string) *Meta {

	if meta != nil {
		return meta
	}

	meta = &Meta{
		target:       target,
		thread:       thread,
		duration:     duration,
		method:       method,
		bodyPath:     bodyPath,
		body:         body,
		headerSlice:  headerSlice,
		headerMap:    headerMap,
		protocol:     protocol,
		protocolFile: protocolFile,
		fqn:          fqn,
	}
	return meta
}

func GetMeta() *Meta {
	return meta
}

func (m *Meta) Protocol() string {
	return m.protocol
}
func (m *Meta) ProtocolFile() string {
	return m.protocolFile
}
func (m *Meta) Fqn() string {
	return m.fqn
}
func (m *Meta) Target() string {
	return m.target
}
func (m *Meta) Body() string {
	return m.body
}

func (m *Meta) HeaderMap() map[string]string {
	return m.headerMap
}

func (m *Meta) SetRespCh(respCh chan *Response) *Meta {
	m.respCh = respCh
	return meta
}

func (m *Meta) RespCh() chan *Response {
	return meta.respCh
}
