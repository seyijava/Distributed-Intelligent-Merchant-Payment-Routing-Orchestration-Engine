package service

import (
	"bigdataconcept/fintech/intelligent/payment/routing/config"
	"bigdataconcept/fintech/intelligent/payment/routing/domain"
	"bigdataconcept/fintech/intelligent/payment/routing/rule"
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

type RoutingManager struct {
	 *KafkaPublisher
	  RoutingRule rule.RuleType
}

func NewRoutingManager(kafkaConfig config.KafkaPublisherConfig,ruleType rule.RuleType)  *RoutingManager {
	 kafkaPublisher := NewKafkaPublisher(kafkaConfig)

	 routingManager := &RoutingManager{KafkaPublisher: kafkaPublisher,RoutingRule: ruleType}
	 return routingManager
}
func (routingManager *RoutingManager) RouteMerchantPayment(merchantPaymentRequest *domain.MerchantPaymentRequest) error {

	routingService := rule.NewRoutingRuleService(routingManager.RoutingRule)
	if channel, err := routingService.ExecutePaymentRoutingRule(merchantPaymentRequest); err != nil{
		return err
	}else{
		if payload, err := json.Marshal(merchantPaymentRequest); err == nil{
			log.Infof("Payment Service Provider Channle %s", channel)
		routingManager.KafkaPublisher.PublishMessage(payload,channel)
		}else{
			return err
		}
	}
	return nil
}