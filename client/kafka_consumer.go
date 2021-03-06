package client

import (
	"fmt"
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"os"
	"os/signal"
	"shop/logger"
	"shop/pprofutil"
	wshandler "shop/service"
	"shop/util"
	"sync"
)

var Address = []string{}

func StartKafkaConsumer(kafkaAdress string) {
	defer util.WaitGroup.Done()
	defer util.PrintFuncName()
	topic := []string{"nginx_log"}

	var wg = &sync.WaitGroup{}
	wg.Add(2)
	Address = append(Address, kafkaAdress)

	//广播式消费
	go clusterConsumerWebsocket(wg, Address, topic, "group-1")

	go clusterConsumerRpc(wg, Address, []string{loggerTopic}, "group-2")
	wg.Wait()
}

// 支持brokers cluster的消费者
func clusterConsumerWebsocket(wg *sync.WaitGroup, brokers, topics []string, groupId string) {
	defer wg.Done()

	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	// init consumer
	consumer, err := cluster.NewConsumer(brokers, groupId, topics, config)
	if err != nil {
		logger.Info.Println("%s: sarama.NewSyncProducer err, message=%s \n", groupId, err)
		return
	} else {
		logger.Info.Println("消费者websocket 成功建立")
	}
	defer func() {
		logger.Info.Println("消费者 websocket 完全退出")
	}()

	// trap SIGINT to trigger a shutdown
	//signals := make(chan os.Signal, 1)
	//signal.Notify(signals, os.Interrupt)

	// consume errors
	go func() {
		for err := range consumer.Errors() {
			logger.Info.Println("消费者组%s: 出错，Error: %s\n", groupId, err.Error())
		}
	}()

	// consume notifications
	go func() {
		for ntf := range consumer.Notifications() {
			logger.Info.Println("%s:Rebalanced: %+v \n", groupId, ntf)
		}
	}()

	// consume messages, watch signals
	var successes int
Loop:
	for {
		select {
		case msg, ok := <-consumer.Messages():
			if ok {
				wshandler.WebsocketChan <- string(msg.Value)

				logger.Info.Println("kafka 消费者 websocket 消费消息")
				fmt.Fprintf(
					os.Stdout,
					"消费组ID=%s，主题=%s，分区=%d，offset=%d，key=%s，value=%s\n",
					groupId, msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
				consumer.MarkOffset(msg, "") // mark message as processed
				successes++
			}
		//case <-signals:
		//	break Loop
		case <-util.ChanStop:
			fmt.Fprintf(os.Stdout, "kafka消费者websocket开始执行退出 %s consume %d messages \n", groupId, successes)
			pprofutil.SaveMemProf()
			pprofutil.ProfGc()
			pprofutil.SaveBlockProfile()
			consumer.Close()
			break Loop
		}
	}
}

func clusterConsumerRpc(wg *sync.WaitGroup, brokers, topics []string, groupId string) {
	defer wg.Done()

	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	// init consumer
	consumer, err := cluster.NewConsumer(brokers, groupId, topics, config)
	if err != nil {
		logger.Info.Println("%s: sarama.NewSyncProducer err, message=%s \n", groupId, err)
		return
	} else {
		logger.Info.Println("消费者成功建立")
	}
	defer func() {
		logger.Info.Println("消费者rpc 完全退出")
	}()

	// trap SIGINT to trigger a shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// consume errors
	go func() {
		for err := range consumer.Errors() {
			logger.Info.Println("消费者组%s: 出错，Error: %s\n", groupId, err.Error())
		}
	}()

	// consume notifications
	go func() {
		for ntf := range consumer.Notifications() {
			logger.Info.Println("%s:Rebalanced: %+v \n", groupId, ntf)
		}
	}()

	// consume messages, watch signals
	var successes int
Loop:
	for {
		select {
		case msg, ok := <-consumer.Messages():
			if ok {
				GrpcCall(string(msg.Value))

				//logger.Info.Println("kafka 消费者消费消息")
				fmt.Fprintf(
					os.Stdout,
					"消费组ID=%s，主题=%s，分区=%d，offset=%d，key=%s，value=%s\n",
					groupId, msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
				consumer.MarkOffset(msg, "") // mark message as processed
				successes++
			}
		//case <-signals:
		//	break Loop
		case <-util.ChanStop:
			fmt.Fprintf(os.Stdout, "kafka消费者rpc开始执行退出， %s consume %d messages \n", groupId, successes)
			consumer.Close()
			break Loop
		}
	}
}
