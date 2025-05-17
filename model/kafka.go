package model

import "github.com/IBM/sarama"

type KafkaProducer struct {
	Producer sarama.SyncProducer
}

type KafkaConsumer struct {
	Consumer      sarama.Consumer
	ConsumerGroup sarama.ConsumerGroup
}

const Msg_hello_id = 1001
const Topic_hello = "hello"

type BaseMsg struct {
	MsgID   int64  `json:"msgID"`
	MsgData string `json:"msgData"`
}

type ProduceMsg struct {
	ToUid   string `json:"toUid"`
	SendUid string `json:"sendUid"`
	MsgType int    `json:"msgType"`
	Content string `json:"content"`
}
