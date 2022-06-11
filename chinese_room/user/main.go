package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"user/core"
	"user/model"
	"user/pb"
)

func main(){
	model.DataBaseInit()


	listen, err := net.Listen("tcp", ":8090")
	defer listen.Close()
	if err!=nil{
		fmt.Println("设置监听失败",err)
	}
	fmt.Println("服务启动")
	//初始化服务
	grpcServer:=grpc.NewServer()
	//注册服务
	pb.RegisterUserServiceServer(grpcServer,new(core.UserService))
	//设置监听


		if err = grpcServer.Serve(listen); err != nil {
			log.Fatalf("failed to serve: %v", err)



	}

}
