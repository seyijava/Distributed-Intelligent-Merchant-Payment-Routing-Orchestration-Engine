package service

import (
	"bigdataconcept/fintech/intelligent/payment/routing/config"
	"bigdataconcept/fintech/intelligent/payment/routing/domain"
	"context"
	"encoding/json"
	"errors"
	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"sync"
)

type KafkaSubscriber struct {
	config.KafkaSubscriberConfig
	routingManager *RoutingManager
	closing        chan struct{}
	closed         bool
	subscribersWg  sync.WaitGroup
	groupConsumer  sarama.ConsumerGroup
}

func NewKafkaSubscriber(subscriberConfig config.KafkaSubscriberConfig, routingManager *RoutingManager) *KafkaSubscriber {

	subscriberConfig.SetDefaults()
	kafkaSubscriber := &KafkaSubscriber{
		routingManager:        routingManager,
		closing:               make(chan struct{}),
		KafkaSubscriberConfig: subscriberConfig,
	}
	return kafkaSubscriber

}
func (subscriber *KafkaSubscriber) Subscribe(wg *sync.WaitGroup, ctx context.Context) error {

	if subscriber.closed {
		errors.New("subscriber closed")
	}
	log.Infof("Start Consuming Kafka Topic %v", subscriber.KafkaSubscriberConfig.Topic)
	groupConsumer, err := sarama.NewConsumerGroup(subscriber.KafkaSubscriberConfig.Brokers, subscriber.KafkaSubscriberConfig.ConsumerGroup, subscriber.KafkaSubscriberConfig.OverwriteSaramaConfig)

	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}
	subscriber.groupConsumer = groupConsumer
	err = subscriber.consumeGroupMessage(wg, ctx, groupConsumer, subscriber.KafkaSubscriberConfig.Topic)
	if err != nil {
		panic(err)
		return err
	}
	return nil
}

func (subscriber *KafkaSubscriber) consumeGroupMessage(wg *sync.WaitGroup, ctx context.Context, groupConsumer sarama.ConsumerGroup, topic string) error {
	messageHandler := createMessageHandler(subscriber.routingManager)
	handler := &consumerGroupHandler{ctx: ctx, messageHandler: messageHandler, ready: make(chan bool)}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {

			if err := groupConsumer.Consume(ctx, []string{topic}, handler); err != nil {

				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				log.Error(ctx.Err())
				return
			}
			handler.ready = make(chan bool)
		}
	}()
	<-handler.ready // Await till the consumer has been set up
	log.Println("Payment Routing Engine up and running....")
	return nil
}

type MessageHandler struct {
	routingManger *RoutingManager
}

func createMessageHandler(routingManager *RoutingManager) *MessageHandler {
	return &MessageHandler{routingManager}
}

type consumerGroupHandler struct {
	ctx            context.Context
	messageHandler *MessageHandler
	ready          chan bool
}

func (groupHandler consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	close(groupHandler.ready)
	return nil
}

func (consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (consumerGroupHandler consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {


		if err := consumerGroupHandler.messageHandler.processMessage(message); err != nil {
			log.Error(err)
		}
		sess.MarkMessage(message, "")
	}
	return nil
}

func (messageHandler *MessageHandler) processMessage(kafkaMsg *sarama.ConsumerMessage) error {
	paymentPayload := kafkaMsg.Value
	merchantPaymentRequest := &domain.MerchantPaymentRequest{}
	if err := json.Unmarshal(paymentPayload, merchantPaymentRequest); err == nil {
		log.Infof("Message claimed: value =  timestamp = %v, topic = %s",kafkaMsg.Timestamp,kafkaMsg.Topic)
		err = messageHandler.routingManger.RouteMerchantPayment(merchantPaymentRequest)
		return err
	}
	return nil
}

func (subscriber *KafkaSubscriber) Close() error {
	log.Info("Stopping Kafka Consumer")
	subscriber.closed = true
	return subscriber.groupConsumer.Close()
}
