package service

import (
	"bigdataconcept/fintech/intelligent/payment/routing/config"
	"bigdataconcept/fintech/intelligent/payment/routing/rule"
	"context"
	"sync"
)

func NewPaymentRoutingEngine(kafkaPubConfig config.KafkaPublisherConfig, subscriberConfig config.KafkaSubscriberConfig, ruleType rule.RuleType) *PaymentRoutingEngine {
	routingManager := NewRoutingManager(kafkaPubConfig, ruleType)
	subscriber := NewKafkaSubscriber(subscriberConfig, routingManager)
	return &PaymentRoutingEngine{subscriber: subscriber}
}

type PaymentRoutingEngine struct {
	subscriber *KafkaSubscriber
}

func (paymentRoutingEngine *PaymentRoutingEngine) Close() error {
	return paymentRoutingEngine.subscriber.Close()
}

func (paymentRoutingEngine *PaymentRoutingEngine) Start(ctx context.Context, wg *sync.WaitGroup) error {

	return paymentRoutingEngine.subscriber.Subscribe(wg, ctx)
}
