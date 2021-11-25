package reflect

import (
	"context"
	"encoding/json"
	"fmt"
	v1 "github.com/crossoverJie/ptg/reflect/gen"
	"github.com/crossoverJie/ptg/reflect/gen/user"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"google.golang.org/grpc"
	"log"
	"net"
	"testing"
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
	filename := "gen/test.proto"
	parse, err := NewParse(filename)
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}

	mds, err := parse.MethodDescriptor("order.v1.OrderService", "Create")
	if err != nil {
		panic(err)
	}
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.DialContext(context.Background(), "127.0.0.1:5000", opts...)
	stub := grpcdynamic.NewStub(conn)
	rpc, err := parse.InvokeRpc(context.Background(), stub, mds, data)
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

	fmt.Println("gRPC server started at ", port)
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}

type Order struct {
	v1.UnimplementedOrderServiceServer
}

func (o *Order) Create(ctx context.Context, in *v1.OrderApiCreate) (*v1.Order, error) {

	fmt.Println(in.OrderId)
	return &v1.Order{
		OrderId: in.OrderId,
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
	return &user.User{UserId: in.UserId}, nil
}
