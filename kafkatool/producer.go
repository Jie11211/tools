package kafkatool

import (
	"strconv"
	"time"

	"github.com/Shopify/sarama"
)

func (p *Producer) NewMsg(topic, value string) *sarama.ProducerMessage {
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Key:       sarama.StringEncoder(strconv.FormatInt(time.Now().UTC().UnixNano(), 10)),
		Timestamp: time.Now(),
		Value:     sarama.StringEncoder(value),
	}
	return msg
}

func (p *Producer) SendMsg(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	if p.AsyncProducer != nil {
		input := p.AsyncProducer.Input()
		input <- msg
		success := <-p.AsyncProducer.Successes()
		partition = success.Partition
		offset = success.Offset
		err = <-p.AsyncProducer.Errors()
		return
	}
	partition, offset, err = p.SyncProducer.SendMessage(msg)
	return
}
