package kafkatool

import (
	"context"
	"log"

	"github.com/Shopify/sarama"
)

func (c consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error { return nil }

func (c consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (c consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	c.Fn(sess, claim)
	return nil
}

// 消费者组默认方法
func (k *Kafkatool) DefaultConsumerGroupFn(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) {
	for msg := range claim.Messages() {
		log.Default().Println(msg)
		sess.MarkMessage(msg, "")
	}
}

// 消费者组处理
func (k *Kafkatool) ConsumerGroupDo(consumerGroup sarama.ConsumerGroup, topic []string, ctx context.Context, f func(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim)) {
	for {
		handler := consumerGroupHandler{Fn: f}
		err := consumerGroup.Consume(ctx, topic, handler)
		if err != nil {
			panic(any(err))
		}
	}
}
