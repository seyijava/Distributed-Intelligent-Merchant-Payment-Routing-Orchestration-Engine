package service
import (
	"bigdataconcept/fintech/intelligent/payment/routing/config"
	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
)

type KafkaPublisher struct {
	config.KafkaPublisherConfig
	producer sarama.SyncProducer

}

func NewKafkaPublisher(kafkaConfig config.KafkaPublisherConfig) *KafkaPublisher {

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = kafkaConfig.Retries                   // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(kafkaConfig.Brokers, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}
	return &KafkaPublisher{KafkaPublisherConfig:kafkaConfig,producer:producer}
}



func (publisher *KafkaPublisher) PublishMessage(payload []byte, topic string) error  {

	msg := &sarama.ProducerMessage{Topic: topic, Value: sarama.StringEncoder(payload)}
	publisher.producer.SendMessage(msg)

	return nil
}