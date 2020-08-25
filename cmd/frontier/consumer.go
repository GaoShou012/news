package main

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	//sarama_cluster "github.com/bsm/sarama-cluster"
	"github.com/golang/glog"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	//"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wchatv1/config"
	"wchatv1/src/room"
	"wchatv1/utils"
)

var (
	groupId int
)

func main() {
	service := micro.NewService(
		micro.Flags(
			&cli.IntFlag{Name: "group_id", Usage: "group id"},
		),
	)
	service.Init(micro.Action(func(c *cli.Context) error {
		groupId = c.Int("group_id")
		return nil
	}))

	utils.Micro.InitV2()
	utils.Micro.LoadSource()

	if err := utils.Micro.LoadConfig(config.RoomServiceConfig); err != nil {
		panic(err)
	}
	fmt.Println("load config room-service", config.RoomServiceConfig)

	if err := utils.Micro.LoadConfig(config.KafkaClusterConfig); err != nil {
		panic(err)
	}
	fmt.Println("load config kafka-cluster", config.KafkaClusterConfig)

	kafkaAddr := config.KafkaClusterConfig.Addr
	kafkaTopics := []string{config.RoomServiceConfig.Topic}
	//kafkaGroup := fmt.Sprintf("%d", groupId)
	//broker := BrokerKafka{}
	//broker.Init(kafkaAddr, kafkaTopics, kafkaGroup)

	//version, err := sarama.ParseKafkaVersion("2.3.0")
	//if err != nil {
	//	panic(err)
	//}
	//conf := sarama.NewConfig()
	//conf.Version = version
	//addr, topics, group := kafkaAddr, kafkaTopics, kafkaGroup
	//cli, err := sarama.NewConsumerGroup(addr, group, conf)
	//go func() {
	//	ctx, cancel := context.WithCancel(context.Background())
	//	defer cancel()
	//	for {
	//		if err := cli.Consume(ctx, topics, &Consumer{Id: 1}); err != nil {
	//			return
	//		}
	//		if ctx.Err() != nil {
	//			return
	//		}
	//	}
	//}()
	consumerCount := 0
	go func() {
		for {
			fmt.Println(consumerCount)
			time.Sleep(time.Second)
		}
	}()

	//conf := sarama_cluster.NewConfig()
	////conf.ChannelBufferSize = 10000
	//conf.Consumer.Offsets.AutoCommit.Enable = true
	//conf.Consumer.Offsets.AutoCommit.Interval = time.Second
	//conf.Consumer.Offsets.CommitInterval = time.Second
	//conf.Consumer.Return.Errors = true
	//conf.Group.Return.Notifications = true
	//
	//addr, topics := kafkaAddr, kafkaTopics
	//consumer, err := sarama_cluster.NewConsumer(addr, "1", topics, conf)
	//if err != nil {
	//	panic(err)
	//}
	//
	//// trap SIGINT to trigger a shutdown.
	//signals := make(chan os.Signal, 1)
	//signal.Notify(signals, os.Interrupt)
	//
	//// consume errors
	//go func() {
	//	for err := range consumer.Errors() {
	//		log.Printf("Error: %s\n", err.Error())
	//	}
	//}()
	//
	//// consume notifications
	//go func() {
	//	for ntf := range consumer.Notifications() {
	//		log.Printf("Rebalanced: %+v\n", ntf)
	//	}
	//}()
	//
	//// consume messages, watch signals
	//for {
	//	select {
	//	case msg, ok := <-consumer.Messages():
	//		if ok {
	//			//fmt.Fprintf(os.Stdout, "%s/%d/%d\t%s\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
	//			consumerCount++
	//			consumer.MarkOffset(msg, "") // mark message as processed
	//		}
	//	case <-signals:
	//		return
	//	}
	//}

	conf := sarama.NewConfig()
	conf.ChannelBufferSize = 100000

	addr, topic := kafkaAddr, kafkaTopics[0]
	consumer, err := sarama.NewConsumer(addr, conf)
	if err != nil {
		panic(err)
	}
	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		panic(err)
	}
	fmt.Println("topic", topic)
	fmt.Println("partition list", partitionList)

	go func() {
		for partition := range partitionList {
			pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
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

var _ sarama.ConsumerGroupHandler = &Consumer{}

type Consumer struct {
	Id int
}

func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	fmt.Println("i am setup")
	//panic("implement me")
	return nil
}

func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	fmt.Println("cleanup")
	//panic("implement me")
	return nil
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		//log.Printf("Message claimed: key = %s, value = %v, topic = %s, partition = %v, offset = %v", string(message.Key), string(message.Value), message.Topic, message.Partition, message.Offset)

		fmt.Println("consumer id", c.Id)
		room.Agent.BroadcastCache <- message.Value
		fmt.Println(message.Value)
		//if msg, err := room.Codec.Decode(message.Value); err == nil {
		//	room.Agent.Broadcast(msg)
		//} else {
		//	glog.Errorln(err)
		//}

		session.MarkMessage(message, "")
	}
	return nil
}

type BrokerKafka struct {
	addr   []string
	topics []string
	group  string
	cli    sarama.ConsumerGroup
	recv   chan *sarama.ConsumerMessage
}

func (b *BrokerKafka) Init(addr []string, topics []string, group string) error {
	version, err := sarama.ParseKafkaVersion("2.3.0")
	if err != nil {
		return err
	}
	conf := sarama.NewConfig()
	conf.Version = version
	b.addr, b.topics, b.group = addr, topics, group
	b.recv = make(chan *sarama.ConsumerMessage, 100000)
	b.cli, err = sarama.NewConsumerGroup(b.addr, b.group, conf)
	return err
}
func (b *BrokerKafka) Subscribe(consumer sarama.ConsumerGroupHandler) {
	fmt.Println("Subscribe")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for {
		if err := b.cli.Consume(ctx, b.topics, consumer); err != nil {
			return
		}
		if ctx.Err() != nil {
			return
		}
	}
	fmt.Println("End subscribe")
}
