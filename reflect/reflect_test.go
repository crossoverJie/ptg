package reflect

import (
	"context"
	"encoding/json"
	"fmt"
	v1 "github.com/crossoverJie/ptg/reflect/gen"
	"github.com/crossoverJie/ptg/reflect/gen/user"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"strings"
	"testing"
	"time"
)

func TestParseProto(t *testing.T) {
	filename := "gen/test.proto"
	parse, err := NewParse(filename)
	if err != nil {
		panic(err)
	}
	maps := parse.ServiceInfoMaps()
	fmt.Println(maps)
}

func TestRequestJSON(t *testing.T) {
	filename := "gen/test.proto"
	parse, err := NewParse(filename)
	if err != nil {
		panic(err)
	}
	json, err := parse.RequestJSON("order.v1.OrderService", "Create")
	if err != nil {
		panic(err)
	}
	fmt.Println(json)
}

func TestParseReflect_InvokeRpc(t *testing.T) {
	data := `{"order_id":20,"user_id":[20],"remark":"Hello","reason_id":[10]}`
	metaStr := `{"lang":"zh"}`
	var m map[string]string
	err := json.Unmarshal([]byte(metaStr), &m)
	if err != nil {
		panic(err)
	}
	filename := "gen/test.proto"
	parse, err := NewParse(filename)
	if err != nil {
		panic(err)
	}

	mds, err := parse.MethodDescriptor("order.v1.OrderService", "Create")
	if err != nil {
		panic(err)
	}
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.DialContext(context.Background(), "127.0.0.1:6001", opts...)
	stub := grpcdynamic.NewStub(conn)

	// metadata
	// create a new context with some metadata
	//md := metadata.Pairs("name", "v1", "k1", "v2", "k2", "v3")
	md := metadata.New(m)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	rpc, err := parse.InvokeRpc(ctx, stub, mds, data)
	if err != nil {
		panic(err)
	}
	fmt.Println(rpc.String())
	fmt.Println("=========")
	//marshal ,_:= proto.Marshal(rpc)
	marshalIndent, _ := json.MarshalIndent(rpc, "", "\t")
	fmt.Println(string(marshalIndent))
}

func TestServer(t *testing.T) {
	port := 6001
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	v1.RegisterOrderServiceServer(grpcServer, &Order{})
	//reflection.Register(grpcServer)

	fmt.Println("gRPC server started at ", port)
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}

type Order struct {
	v1.UnimplementedOrderServiceServer
}

func (o *Order) Create(ctx context.Context, in *v1.OrderApiCreate) (*v1.Order, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.DataLoss, "failed to get metadata")
	}
	fmt.Println(md)

	time.Sleep(200 * time.Millisecond)
	fmt.Println(in.OrderId)
	return &v1.Order{
		OrderId: in.OrderId,
		Reason:  nil,
	}, nil
}

func (o *Order) Close(ctx context.Context, req *v1.CloseApiCreate) (*v1.Order, error) {
	log.Println(req)
	time.Sleep(200 * time.Millisecond)
	return &v1.Order{
		OrderId: 1000,
		Reason:  nil,
	}, nil
}

func TestParseServiceMethod(t *testing.T) {
	s, m, err := ParseServiceMethod("order.v1.OrderService.Create")
	fmt.Println(s, m, err)
}

func TestUserServer(t *testing.T) {
	port := 7001
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	user.RegisterUserServiceServer(grpcServer, &User{})

	fmt.Println("gRPC user server started at ", port)
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}

type User struct {
	user.UnimplementedUserServiceServer
}

func (*User) Create(ctx context.Context, in *user.UserApiCreate) (*user.User, error) {
	time.Sleep(200 * time.Millisecond)
	return &user.User{UserId: in.UserId}, nil
}

func TestCommon(t *testing.T) {
	x := "order.v1.OrderService.Detail-2"
	fmt.Println(strings.Split(x, "-")[1])
}
