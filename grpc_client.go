package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"google.golang.org/grpc"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var errNoMethodNameSpecified = errors.New("no method name specified")

var (
	one    sync.Once
	client *grpcClient
)

type grpcClient struct {
	stub    grpcdynamic.Stub
	mtd     *desc.MethodDescriptor
	message []*dynamic.Message
}

func NewGrpcClient() Client {
	one.Do(func() {
		if client == nil {
			var (
				opts       []grpc.DialOption
				importPath []string
			)

			client = &grpcClient{}
			opts = append(opts, grpc.WithInsecure())
			conn, err := grpc.DialContext(context.Background(), target, opts...)
			if err != nil {
				panic(fmt.Sprintf("create grpc connection err %v", err))
			}
			client.stub = grpcdynamic.NewStub(conn)

			dir := filepath.Dir(protocolFile)
			importPath = append(importPath, dir)
			client.mtd, err = GetMethodDescFromProto(fqn, protocolFile, importPath)
			if err != nil {
				panic(fmt.Sprintf("create grpc MethodDescriptor err %v", err))
			}

			client.message, err = getDataForCall(client.mtd)
			if err != nil {
				panic(fmt.Sprintf("create grpc message err %v", err))
			}

		}
	})
	return client
}

func (g *grpcClient) Request() (*Response, error) {
	//fmt.Printf("%p \n", &*g)\
	start := time.Now()
	rpc, err := g.stub.InvokeRpc(context.Background(), g.mtd, g.message[0])
	r := &Response{
		RequestTime: time.Since(start),
	}
	SlowRequestTime = r.slowRequest()
	FastRequestTime = r.fastRequest()
	if err != nil {
		return nil, err
	}
	r.ResponseSize = len(rpc.String())
	return r, nil
}

func GetMethodDescFromProto(call, proto string, imports []string) (*desc.MethodDescriptor, error) {
	p := &protoparse.Parser{ImportPaths: imports}

	filename := proto
	if filepath.IsAbs(filename) {
		filename = filepath.Base(proto)
	}

	fds, err := p.ParseFiles(filename)
	if err != nil {
		return nil, err
	}

	fileDesc := fds[0]

	files := map[string]*desc.FileDescriptor{}
	files[fileDesc.GetName()] = fileDesc

	return getMethodDesc(call, files)
}

func getMethodDesc(call string, files map[string]*desc.FileDescriptor) (*desc.MethodDescriptor, error) {
	svc, mth, err := parseServiceMethod(call)
	if err != nil {
		return nil, err
	}

	dsc, err := findServiceSymbol(files, svc)
	if err != nil {
		return nil, err
	}
	if dsc == nil {
		return nil, fmt.Errorf("cannot find service %q", svc)
	}

	sd, ok := dsc.(*desc.ServiceDescriptor)
	if !ok {
		return nil, fmt.Errorf("cannot find service %q", svc)
	}

	mtd := sd.FindMethodByName(mth)
	if mtd == nil {
		return nil, fmt.Errorf("service %q does not include a method named %q", svc, mth)
	}

	return mtd, nil
}

func parseServiceMethod(svcAndMethod string) (string, string, error) {
	if len(svcAndMethod) == 0 {
		return "", "", errNoMethodNameSpecified
	}
	if svcAndMethod[0] == '.' {
		svcAndMethod = svcAndMethod[1:]
	}
	if len(svcAndMethod) == 0 {
		return "", "", errNoMethodNameSpecified
	}
	switch strings.Count(svcAndMethod, "/") {
	case 0:
		pos := strings.LastIndex(svcAndMethod, ".")
		if pos < 0 {
			return "", "", newInvalidMethodNameError(svcAndMethod)
		}
		return svcAndMethod[:pos], svcAndMethod[pos+1:], nil
	case 1:
		split := strings.Split(svcAndMethod, "/")
		return split[0], split[1], nil
	default:
		return "", "", newInvalidMethodNameError(svcAndMethod)
	}
}

func newInvalidMethodNameError(svcAndMethod string) error {
	return fmt.Errorf("method name must be package.Service.Method or package.Service/Method: %q", svcAndMethod)
}

func findServiceSymbol(resolved map[string]*desc.FileDescriptor, fullyQualifiedName string) (desc.Descriptor, error) {
	for _, fd := range resolved {
		if dsc := fd.FindSymbol(fullyQualifiedName); dsc != nil {
			return dsc, nil
		}
	}
	return nil, fmt.Errorf("cannot find service %q", fullyQualifiedName)
}

func getDataForCall(mtd *desc.MethodDescriptor) ([]*dynamic.Message, error) {
	var inputs []*dynamic.Message
	var err error
	if inputs, err = getPayloadMessages([]byte(body), mtd); err != nil {
		return nil, err
	}

	if len(inputs) > 0 {
		unaryInput := inputs[0]

		return []*dynamic.Message{unaryInput}, nil
	}

	return inputs, nil
}

func getPayloadMessages(inputData []byte, mtd *desc.MethodDescriptor) ([]*dynamic.Message, error) {
	var inputs []*dynamic.Message
	data := inputData
	inputs, err := createPayloadsFromJSON(mtd, string(data))
	if err != nil {
		return nil, err
	}

	return inputs, nil
}

func createPayloadsFromJSON(mtd *desc.MethodDescriptor, data string) ([]*dynamic.Message, error) {
	md := mtd.GetInputType()
	var inputs []*dynamic.Message

	if len(data) > 0 {
		inputs = make([]*dynamic.Message, 1)
		inputs[0] = dynamic.NewMessage(md)
		err := jsonpb.UnmarshalString(data, inputs[0])
		if err != nil {
			return nil, errors.New(fmt.Sprintf("create payload json err %v \n", err))
		}
	}

	return inputs, nil
}
