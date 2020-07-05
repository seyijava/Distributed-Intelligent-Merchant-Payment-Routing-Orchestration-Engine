package infracstructure

import (
	"encoding/json"

	"strings"

	"bigdataconcept/fintech/intelligent/payment/routing/gateway/config"
	"bigdataconcept/fintech/intelligent/payment/routing/gateway/domain"
	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"math/rand"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type KafkaPublisher struct {
	config.KafkaConfig
	Producer sarama.SyncProducer
}

func createProducer(brokerList []string, maxRetry int) sarama.SyncProducer {
	config := sarama.NewConfig()
	//config.Producer.RequiredAcks = sarama.WaitForAll //
	config.Producer.Retry.Max = maxRetry
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}
	return producer
}
func NewKafkaPublisher(pubConfig config.KafkaConfig) *KafkaPublisher {
	producer := createProducer(pubConfig.BrokerAddress, pubConfig.Retry)
	return &KafkaPublisher{Producer: producer,KafkaConfig: pubConfig}
}

func (kafkaPublisher *KafkaPublisher) publishMessage(payload []byte) error {

	msg := &sarama.ProducerMessage {
		Topic: kafkaPublisher.KafkaConfig.Topic,
		Value: sarama.StringEncoder(payload),
	}
	_, _, err  := kafkaPublisher.Producer.SendMessage(msg)
	return err
}

func ProcessRequest(merchantRequest *domain.MerchantPaymentRequest, kafkaPublisher *KafkaPublisher) (error, *domain.MerchantPaymentResponse) {
	log.Infof("Publishing Merchant Payment To Kafka Topic [%s] Merchant Payment Reference Number [%s]", kafkaPublisher.KafkaConfig.Topic,merchantRequest.Card.CardNumber)
	merchantPaymentResponse := &domain.MerchantPaymentResponse{}
	if payload, err := json.Marshal(merchantRequest); err != nil {
		log.Error(err)
		return err, nil
	}else{
		if err := kafkaPublisher.publishMessage(payload); err != nil{
			log.Error(err)
			merchantPaymentResponse.MerchantPaymentRef = ""
			merchantPaymentResponse.TransactionRef = "XXXXXXXXXX"
			merchantPaymentResponse.ResponseCode = "99"
			merchantPaymentResponse.ResponseMessage = "Internal Server Error Try again"
		} else {
			paymentReference := strings.ToUpper("PTY/" + merchantRequest.Merchant.MerchantCode + "/" + RandStringBytes(10))
			merchantPaymentResponse.MerchantPaymentRef = merchantRequest.MerchantPaymentReference
			merchantPaymentResponse.TransactionRef = paymentReference
			merchantPaymentResponse.ResponseCode = "00"
			merchantPaymentResponse.ResponseMessage = "Successful"
		}
		//response, _ := json.Marshal(merchantPaymentResponse)
		return nil, merchantPaymentResponse
	}
}

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
