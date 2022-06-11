package model

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName string
	PasswordDigest string
}
const (
	PassWordCost = 12 // 密码加密难度
)

func(us *User)SetPassword(password string)error{
	bytes,err:=bcrypt.GenerateFromPassword([]byte(password),PassWordCost)
	if err!=nil{
		fmt.Println(err.Error())
		return err
	}
	us.PasswordDigest=string(bytes)
	return nil

}
func(us *User)CheckPassword(password string)bool{
	err:=bcrypt.CompareHashAndPassword([]byte(us.PasswordDigest),[]byte(password))
	return  err==nil
}
