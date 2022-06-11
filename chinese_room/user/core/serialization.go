package core

import (
	"user/model"
	"user/pb"
)

func BuildUser(item model.User)*pb.UserModel{
	var user pb.UserModel
	user.UserName=item.UserName
	user.ID= uint32(item.ID)
	user.CreatedAt=item.CreatedAt.Unix()
	user.UpdatedAt=item.UpdatedAt.Unix()
	return &user
}
