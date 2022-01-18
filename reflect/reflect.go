package reflect

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"strings"
)

type ParseReflect struct {
	filename string
	// serviceInfoMap[sds.GetFullyQualifiedName()] = methodInfos
	serviceInfoMap map[string][]string
	fds            *desc.FileDescriptor
}

func NewParse(filename string) (*ParseReflect, error) {
	p := &ParseReflect{filename: filename}
	err := p.parseProto()
	return p, err
}

// Init: parse proto
func (p *ParseReflect) parseProto() error {
	var parser protoparse.Parser
	fds, err := parser.ParseFiles(p.filename)
	if err != nil {
		return errors.WithStack(err)
	}
	p.fds = fds[0]

	serviceInfoMap := make(map[string][]string)
	for _, sds := range fds[0].GetServices() {

		var methodInfos []string
		for _, mds := range sds.GetMethods() {
			methodInfos = append(methodInfos, mds.GetName())
		}

		serviceInfoMap[sds.GetFullyQualifiedName()] = methodInfos

	}
	p.serviceInfoMap = serviceInfoMap
	return nil
}

func (p *ParseReflect) ServiceInfoMaps() map[string][]string {
	return p.serviceInfoMap
}

func (p *ParseReflect) RequestJSON(serviceName, methodName string) (string, error) {
	_, ok := p.serviceInfoMap[serviceName]
	if !ok {
		return "", errors.New("service not found!")
	}

	sds := p.fds.FindService(serviceName)
	mds := sds.FindMethodByName(methodName)
	messageToMap := convertMessageToMap(mds.GetInputType())
	marshalIndent, err := json.MarshalIndent(messageToMap, "", "\t")
	return string(marshalIndent), err
}

// Get Method desc
func (p *ParseReflect) MethodDescriptor(serviceName, methodName string) (*desc.MethodDescriptor, error) {
	_, ok := p.serviceInfoMap[serviceName]
	if !ok {
		return nil, errors.New("service not found!")
	}
	sds := p.fds.FindService(serviceName)
	return sds.FindMethodByName(methodName), nil
}

// make unary RPC
func (p *ParseReflect) InvokeRpc(ctx context.Context, stub grpcdynamic.Stub, mds *desc.MethodDescriptor, data string, opts ...grpc.CallOption) (proto.Message, error) {

	messages, err := createPayloadsFromJSON(mds, data)
	if err != nil {
		return nil, err
	}
	return stub.InvokeRpc(ctx, mds, messages[0], opts...)
}

// make unary server stream RPC
func (p *ParseReflect) InvokeServerStreamRpc(ctx context.Context, stub grpcdynamic.Stub, mds *desc.MethodDescriptor, data string, opts ...grpc.CallOption) (*grpcdynamic.ServerStream, error) {

	messages, err := createPayloadsFromJSON(mds, data)
	if err != nil {
		return nil, err
	}
	return stub.InvokeRpcServerStream(ctx, mds, messages[0], opts...)
}

// make unary client stream RPC
func (p *ParseReflect) InvokeClientStreamRpc(ctx context.Context, stub grpcdynamic.Stub, mds *desc.MethodDescriptor, opts ...grpc.CallOption) (*grpcdynamic.ClientStream, error) {
	return stub.InvokeRpcClientStream(ctx, mds, opts...)
}

// make unary bidi stream RPC
func (p *ParseReflect) InvokeBidiStreamRpc(ctx context.Context, stub grpcdynamic.Stub, mds *desc.MethodDescriptor, opts ...grpc.CallOption) (*grpcdynamic.BidiStream, error) {
	return stub.InvokeRpcBidiStream(ctx, mds, opts...)
}

func convertMessageToMap(message *desc.MessageDescriptor) map[string]interface{} {
	m := make(map[string]interface{})
	for _, fieldDescriptor := range message.GetFields() {
		fieldName := fieldDescriptor.GetName()
		if fieldDescriptor.IsRepeated() {
			// Array temporary nil
			m[fieldName] = nil
			continue
		}
		switch fieldDescriptor.GetType() {
		case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
			m[fieldName] = convertMessageToMap(fieldDescriptor.GetMessageType())
			continue
		}
		m[fieldName] = fieldDescriptor.GetDefaultValue()
	}
	return m
}

func ParseServiceMethod(svcAndMethod string) (string, string, error) {
	if len(svcAndMethod) == 0 {
		return "", "", errors.New("service not found!")
	}
	if svcAndMethod[0] == '.' {
		svcAndMethod = svcAndMethod[1:]
	}
	if len(svcAndMethod) == 0 {
		return "", "", errors.New("service not found!")
	}
	switch strings.Count(svcAndMethod, "/") {
	case 0:
		pos := strings.LastIndex(svcAndMethod, ".")
		if pos < 0 {
			return "", "", errors.New("service not found!")
		}
		return svcAndMethod[:pos], svcAndMethod[pos+1:], nil
	case 1:
		split := strings.Split(svcAndMethod, "/")
		return split[0], split[1], nil
	default:
		return "", "", errors.New("service not found!")
	}
}

func createPayloadsFromJSON(mds *desc.MethodDescriptor, data string) ([]*dynamic.Message, error) {
	md := mds.GetInputType()
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
