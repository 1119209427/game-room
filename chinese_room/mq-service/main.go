package main

import (
	"mq-service/model"
	"mq-service/service"
)

func main() {
	model.DataBaseInit()
	model.RabbitMq()
	forever := make(chan bool)
	service.CreateTask()
	<-forever
}
