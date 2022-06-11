package service

import (
	"encoding/json"
	"fmt"
	"log"
	"mq-service/model"
)

func CreateTask(){
	//初始化管道
	ch,err:=model.RQ.Channel()
	if err!=nil{
		log.Fatalln(err.Error())
	}
	//队列
	queue,err:=ch.QueueDeclare("task_queue",true,false,false,false,nil)
	if err!=nil{
		fmt.Println("创建队列失败:",err)
	}
	//消费呗
	msg,err:=ch.Consume(queue.Name,"",false,false,false,false,nil)
	if err!=nil{
		log.Fatalln(err.Error())
	}
	// 处于一个监听状态，一致监听我们的生产端的生产，所以这里我们要阻塞主进程
	go func() {
		for d:=range msg{
			var task model.Task
			err=json.Unmarshal(d.Body,&task)
			if err != nil {
				panic(err)
			}
			if err=model.DB.Create(&task).Error;err!=nil{
				log.Println("数据库报错信息失败：",err.Error())
			}
			log.Println("Done")
			_ = d.Ack(false)
		}
	}()
}