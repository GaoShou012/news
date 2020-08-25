package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/golang/glog"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wchatv1/config"
	"wchatv1/utils"
)

func main(){
	utils.Micro.InitV2()
	utils.Micro.LoadSource()
	utils.Micro.LoadConfigMust(config.RoomServiceConfig)
	utils.Micro.LoadConfigMust(config.KafkaClusterConfig)

	consumerCount := 0
	go func() {
		for {
			fmt.Println(consumerCount)
			time.Sleep(time.Second)
		}
	}()

	conf := sarama.NewConfig()
	conf.ChannelBufferSize = 100000

	addr, topics := config.KafkaClusterConfig.Addr,config.RoomServiceConfig.Topic
	consumer, err := sarama.NewConsumer(addr, conf)
	if err != nil {
		panic(err)
	}
	partitionList, err := consumer.Partitions(topics)
	if err != nil {
		panic(err)
	}
	fmt.Println("topic", topics)
	fmt.Println("partition list", partitionList)

	go func() {
		for partition := range partitionList {
			pc, err := consumer.ConsumePartition(topics, int32(partition), sarama.OffsetNewest)
			if err != nil {
				glog.Errorln(err)
				return
			}
			go func(pc sarama.PartitionConsumer) {
				defer pc.AsyncClose()
				for {
					<-pc.Messages()
					consumerCount++
					//fmt.Println("msg", msg.Value)
				}
			}(pc)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		switch s := <-c; s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			glog.Infof("got signal %s; stop server", s)
		case syscall.SIGHUP:
			glog.Infof("got signal %s; go to deamon", s)
			continue
		}
		break
	}
}
