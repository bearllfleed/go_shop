package service

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/bearllflee/go_shop/global"
	"github.com/bearllflee/go_shop/model"
)

type KafkaService struct {
	Brokers  []string
	Config   *sarama.Config
	Producer *model.KafkaProducer
	Consumer *model.KafkaConsumer
}

func NewKafkaService() *KafkaService {
	brokers := global.CONFIG.Kafka.Host + ":" + global.CONFIG.Kafka.Port
	ks := &KafkaService{
		Brokers:  []string{brokers},
		Producer: &model.KafkaProducer{},
		Consumer: &model.KafkaConsumer{},
	}

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	ks.Config = config

	producer, err := sarama.NewSyncProducer(ks.Brokers, config)
	if err != nil {
		panic(err)
	}
	ks.Producer.Producer = producer

	consumer, err := sarama.NewConsumer(ks.Brokers, config)
	if err != nil {
		panic(err)
	}
	ks.Consumer.Consumer = consumer

	return ks
}

func (ks *KafkaService) AddConsumerHandler(groupID string, topics []string, handler sarama.ConsumerGroupHandler) {
	consumerGroup, err := sarama.NewConsumerGroup(ks.Brokers, groupID, ks.Config)
	if err != nil {
		panic(err)
	}
	ks.Consumer.ConsumerGroup = consumerGroup
	for {
		err = consumerGroup.Consume(context.Background(), topics, handler)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (ks *KafkaService) ProduceMsg(topic string, message string) {
	msg := sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	partition, offset, err := ks.Producer.Producer.SendMessage(&msg)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("生产消息：", msg, "分区:", partition, "偏移量:", offset)
}
