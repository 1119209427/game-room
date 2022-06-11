package handler

import (
	"context"
	"gateway/pb"
)

func UserRegister(client pb.UserServiceClient,userRep *pb.UserRequest)(*pb.UserDetailResponse,error){


	userResp,err:=client.UserRegister(context.Background(),userRep)

	return userResp,err


}
func UserLogin(client pb.UserServiceClient,userRep *pb.UserRequest)(*pb.UserDetailResponse,error){
	userResp,err:=client.UserLogin(context.Background(),userRep)

	return userResp,err

}

