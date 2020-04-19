package service

import (
	"fmt"
	"google.golang.org/grpc"
	"net"
	pb "shop/rpc/proto"
	"shop/transformer"
	"shop/util"
)

var GrpcSever *grpc.Server
func StartGrpcService(grpcConf transformer.GrpcConf) {
	defer util.WaitGroup.Done()
	defer util.PrintFuncName()
	// 1. 监听
	addr :=  grpcConf.Addr
	listenSocket, err := net.Listen("tcp", grpcConf.Addr)
	if err != nil {
		fmt.Printf("监听异常：%GrpcSever\n", err)
	}
	fmt.Printf("grpc 服务开始监听的地址和端口：%GrpcSever\n", addr)
	// 2.实例化gRPC
	GrpcSever = grpc.NewServer()

	var u = UserInfoService{}
	u.init()
	defer u.destroy()
	// 3.在gRPC上注册微服务
	// 第二个参数类型需要接口类型的变量
	pb.RegisterUserInfoServiceServer(GrpcSever, &u)
	// 4.启动gRPC服务
	fmt.Println("启动gRPC服务")

	GrpcSever.Serve(listenSocket)

}
