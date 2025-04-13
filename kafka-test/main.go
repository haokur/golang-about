package main

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"os"
	"os/signal"
)

func test() {
	// Kafka broker 地址，可以写多个：[]string{"localhost:9092", "localhost:9093"}
	brokers := []string{"localhost:9092"}

	// 创建配置
	config := sarama.NewConfig()
	config.Version = sarama.V2_6_0_0 // 设置为你的 Kafka 版本

	// 创建客户端
	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		log.Fatalf("❌ 无法连接到 Kafka: %v", err)
	}
	defer client.Close()

	// 获取并打印 metadata，测试连接是否成功
	brokersInfo := client.Brokers()
	if len(brokersInfo) == 0 {
		log.Fatal("⚠️ 没有获取到 Kafka brokers")
	}
	fmt.Println("✅ 成功连接到 Kafka，Brokers:")
	for _, b := range brokersInfo {
		addr := b.Addr()
		fmt.Println(" -", addr)
	}
}

type SimpleConsumerGroupHandler struct{}

func (SimpleConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (SimpleConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (SimpleConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	fmt.Println("开始接收消息...")
	for msg := range claim.Messages() {
		fmt.Printf("收到消息，offset=%d, key=%s, value=%s\n", msg.Offset, msg.Key, msg.Value)
		// 标记已处理
		session.MarkMessage(msg, "")
	}
	return nil
}

// 不支持断点续消费，即不支持服务启动前的旧消息
func testOnlyNewMessage() {
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, nil)
	if err != nil {
		log.Fatalf("创建消费者失败：%v", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition("demo-topic", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatal("订阅分区失败：%v", err)
	}
	defer partitionConsumer.Close()

	fmt.Println("开始接收消息。。。")
	for msg := range partitionConsumer.Messages() {
		fmt.Printf("收到消息，key=%s,value=%s\n", msg.Key, msg.Value)
	}
}

// 支持断点续消费
func testGroupMessage() {
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	config.Consumer.Offsets.Initial = sarama.OffsetOldest // 第一次从最早的消息开始

	brokers := []string{"localhost:9092"}
	groupID := "my-simple-group"
	topics := []string{"demo-topic"}

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		log.Fatalf("创建消费者组失败：%v", err)
	}
	defer consumerGroup.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			err := consumerGroup.Consume(ctx, topics, SimpleConsumerGroupHandler{})
			if err != nil {
				log.Fatalf("消费错误：%v", err)
			}
		}
	}()

	// 监听 Ctrl+C
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
	fmt.Println("退出消费者...")
}

func main() {
	testGroupMessage()
}
