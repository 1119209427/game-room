package model

import (
	"fmt"
	"github.com/streadway/amqp"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)


var DB *gorm.DB
func DataBaseInit(){
	dsn := "root:123456@tcp(127.0.0.1:3306)/todo_list?charset=utf8mb4&parseTime=True&loc=Local"
	db,err:=gorm.Open(mysql.Open(dsn),&gorm.Config{})
	if err!=nil{
		panic(err.Error())
	}
	DB=db
	migration()
}
var RQ *amqp.Connection
func RabbitMq(){
	conn:="amqp://guest:guest@localhost:5672/"
	rq,err:=amqp.Dial(conn)
	if err!=nil{
		panic(err.Error())
	}
	RQ=rq
}
func migration(){
	err:=DB.AutoMigrate(&Task{})
	if err!=nil{
		fmt.Println(err.Error())
	}
}
