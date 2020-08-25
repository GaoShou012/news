package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"strconv"
	"time"
	"wchatv1/config"
	"wchatv1/utils"
)

func main(){
	utils.Micro.InitV2()
	utils.Micro.LoadSource()

	if err := utils.Micro.LoadConfig(config.KafkaClusterConfig); err != nil {
		panic(err)
	}
	fmt.Println("kafka-cluster",config.KafkaClusterConfig)
	if err := utils.Micro.LoadConfig(config.RoomServiceConfig); err != nil {
		panic(err)
	}
	fmt.Println("room-service",config.RoomServiceConfig)

	conf := sarama.NewConfig()
	conf.Producer.Retry.Max = 5
	conf.Producer.RequiredAcks = sarama.WaitForAll
	producer , err := sarama.NewAsyncProducer(config.KafkaClusterConfig.Addr,conf)
	if err != nil {
		panic(err)
	}

	strTime := strconv.Itoa(int(time.Now().Unix()))
	for i:=0;i<10000;i++{
		m := &sarama.ProducerMessage{
			Topic:     config.RoomServiceConfig.KafkaGroup,
			Key:       sarama.StringEncoder(strTime),
			Value:     sarama.StringEncoder("hello"),
			Headers:   nil,
			Metadata:  nil,
			Offset:    0,
			Partition: 0,
			Timestamp: time.Time{},
		}
		//producer.SendMessage(m)
		producer.Input() <- m
	}
	for{
		time.Sleep(time.Second)
		fmt.Println(time.Now())
	}
}
