package kafkatool

import "github.com/Shopify/sarama"

type Kafkatool struct {
	Addr   []string
	Config *sarama.Config
}

type Producer struct {
	AsyncProducer sarama.AsyncProducer
	SyncProducer  sarama.SyncProducer
}

type consumerGroupHandler struct {
	Fn func(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim)
}
