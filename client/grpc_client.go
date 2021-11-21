package client

import (
	"context"
	"fmt"
	"github.com/crossoverJie/ptg/meta"
	"github.com/crossoverJie/ptg/reflect"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"google.golang.org/grpc"
	"sync"
	"time"
)

var (
	one sync.Once
	g   *grpcClient
)

type grpcClient struct {
	stub         grpcdynamic.Stub
	mds          *desc.MethodDescriptor
	meta         *meta.Meta
	parseReflect *reflect.ParseReflect
}

func NewGrpcClient(meta *meta.Meta) Client {
	one.Do(func() {
		if g == nil {
			var (
				opts []grpc.DialOption
			)

			g = &grpcClient{
				meta: meta,
			}
			opts = append(opts, grpc.WithInsecure())
			conn, err := grpc.DialContext(context.Background(), meta.Target(), opts...)
			if err != nil {
				panic(fmt.Sprintf("create grpc connection err %v", err))
			}
			g.stub = grpcdynamic.NewStub(conn)

			parse, err := reflect.NewParse(meta.ProtocolFile())
			if err != nil {
				panic(fmt.Sprintf("create grpc parse reflect err %v", err))
			}
			g.parseReflect = parse
			serviceName, methodName, err := reflect.ParseServiceMethod(meta.Fqn())
			if err != nil {
				panic(fmt.Sprintf("parse MethodDescriptor err %v", err))
			}
			g.mds, err = g.parseReflect.MethodDescriptor(serviceName, methodName)
			if err != nil {
				panic(fmt.Sprintf("create grpc MethodDescriptor err %v", err))
			}

		}
	})
	return g
}

func (g *grpcClient) Request() (*meta.Response, error) {
	//fmt.Printf("%p \n", &*g)\
	start := time.Now()
	rpc, err := g.parseReflect.InvokeRpc(context.Background(), g.stub, g.mds, g.meta.Body())
	r := &meta.Response{
		RequestTime: time.Since(start),
	}
	//SlowRequestTime = r.slowRequest()
	//FastRequestTime = r.fastRequest()
	meta.GetResult().SetSlowRequestTime(r.SlowRequest()).
		SetFastRequestTime(r.FastRequest())
	if err != nil {
		return nil, err
	}
	r.ResponseSize = len(rpc.String())
	return r, nil
}
