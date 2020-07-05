package config

import (
	"bigdataconcept/fintech/intelligent/payment/routing/rule"
	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
	"time"
)

const (

	DROOL_RULE = "DROOL"
	MACHINE_LEARNING_RULE="ML"

	MONGO_CONNECTIONURL = "MONGO_CONNECTIONURL"
	MONGO_DB ="MONGO_DB"
	MONGO_COLLECTION = "MONGO_COLLECTION"
	MONGO_USERNAME = "MONGO_USERNAME"
	MONGO_PASSWORD = "MONGO_PASSWORD"


	SERVER_HOST = "SERVER_HOST"
	SERVER_PORT = "SERVER_PORT"
	SERVER_GOROUTINEPOOLSIZE = "SERVER_GOROUTINEPOOLSIZE"

	KAFKA_BROKERS = "KAFKA_BROKERS"

	KAFKAPUB_RETRY = "KAFKAPUB_RETRY"

	KAFKASUB_CONSUNMERGROUP = "KAFKASUB_CONSUNMERGROUP"
	KAFKASUB_TOPIC = "KAFKASUB_TOPIC"

	ROUTING_RULETYPE= "ROUTING_RULETYPE"



)

type GetConfig func(*AppConfig,  map[string]interface{}) *AppConfig


func GetServerConfig() GetConfig {
	return  func(config *AppConfig, configVal map[string]interface{}) *AppConfig {
		config.ServerConfig.IP = configVal[SERVER_HOST].(string)
		config.ServerConfig.Port = configVal[SERVER_PORT].(string)
		return config
	}
}

func GenerateGetServerConfig(v *viper.Viper) map[string]interface{}{
	configurationMap := make(map[string]interface{})
	if ip := v.Get(SERVER_HOST); ip == nil{
		configurationMap[SERVER_HOST] = v.Get("serverConfig.host")
	}else{

		configurationMap[SERVER_HOST] = ip
	}
	if port := v.Get(SERVER_PORT); port == nil{
		configurationMap[SERVER_PORT] = v.Get("serverConfig.port")
	}else{

		configurationMap[SERVER_PORT] = port
	}
	return configurationMap
}

func GetKafkaPubConfig() GetConfig {
	return  func(config *AppConfig, configVal map[string]interface{}) *AppConfig {
		config.KafkaPublisherConfig.Brokers = strings.Split(configVal[KAFKA_BROKERS].(string), ",")
		return config
	}
}

func GenerateKafkaPubConfig(v *viper.Viper) map[string]interface{}{
	configurationMap := make(map[string]interface{})
	if brokerUrl := v.Get(KAFKA_BROKERS); brokerUrl == nil{
		configurationMap[KAFKA_BROKERS] = v.Get("kafkaConfig.brokerUrl")
	}else{

		configurationMap[KAFKA_BROKERS] = brokerUrl
	}
	return configurationMap
}

func GetKafkaSubConfig() GetConfig {
	return  func(config *AppConfig, configVal map[string]interface{}) *AppConfig {
		config.KafkaSubscriberConfig.Brokers = strings.Split(configVal[KAFKA_BROKERS].(string), ",")
		config.KafkaSubscriberConfig.ConsumerGroup = configVal[KAFKASUB_CONSUNMERGROUP].(string)
		config.KafkaSubscriberConfig.Topic = configVal[KAFKASUB_TOPIC].(string)
		return config
	}
}

func GenerateKafkaSubConfig(v *viper.Viper) map[string]interface{}{
	configurationMap := make(map[string]interface{})
	if brokerUrl := v.Get(KAFKA_BROKERS); brokerUrl == nil{
		configurationMap[KAFKA_BROKERS] = v.Get("kafkaConfig.brokerUrl")
	}else{

		configurationMap[KAFKA_BROKERS] = brokerUrl
	}
	if group := v.Get(KAFKASUB_CONSUNMERGROUP); group == nil{
		configurationMap[KAFKASUB_CONSUNMERGROUP] = v.Get("kafkaSubConfig.group")
	}else{
		configurationMap[KAFKASUB_CONSUNMERGROUP] = group
	}
	if topic := v.Get(KAFKASUB_TOPIC); topic == nil{
		configurationMap[KAFKASUB_TOPIC] = v.Get("kafkaSubConfig.topic")
	}else{
		configurationMap[KAFKASUB_TOPIC] = topic
	}

	return configurationMap
}



func GetRoutingRuleConfig() GetConfig {
	return  func(config *AppConfig, configVal map[string]interface{}) *AppConfig {
		config.RoutingRuleConfig.RoutingRuleType = configVal[ROUTING_RULETYPE].(string)
		return config
	}
}

func GenerateRoutingRuleConfig(v *viper.Viper) map[string]interface{}{
	configurationMap := make(map[string]interface{})
	if ruleType := v.Get(ROUTING_RULETYPE); ruleType == nil{
		configurationMap[ROUTING_RULETYPE] = v.Get("routingRuleConfig.routingRuleType")
	}else{
		log.Info(ruleType)
		configurationMap[ROUTING_RULETYPE] = ruleType
	}
	return configurationMap
}

type AppConfig struct {
	RoutingRuleConfig
	KafkaSubscriberConfig
	KafkaPublisherConfig
	ServerConfig
}

type ServerConfig struct {
	 IP string
	 Port string
}
type RoutingRuleConfig struct {
     RoutingRuleType string

}
type KafkaSubscriberConfig struct {
	// Kafka brokers list.
	Brokers []string
	OverwriteSaramaConfig *sarama.Config
	// Kafka consumer group.
	// When empty, all messages from all partitions will be returned.
	ConsumerGroup string

	// How long after Nack message should be redelivered.
	NackResendSleep time.Duration

	// How long about unsuccessful reconnecting next reconnect will occur.
	ReconnectRetrySleep time.Duration

	//InitializeTopicDetails *sarama.TopicDetail

	Topic string
}


func (c *KafkaSubscriberConfig) SetDefaults() {
	if c.OverwriteSaramaConfig == nil {
		c.OverwriteSaramaConfig = DefaultSaramaSubscriberConfig()
	}
	if c.NackResendSleep == 0 {
		c.NackResendSleep = time.Millisecond * 100
	}
	if c.ReconnectRetrySleep == 0 {
		c.ReconnectRetrySleep = time.Second
	}
}

func DefaultSaramaSubscriberConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Version =   sarama.V0_10_2_0
	config.Consumer.Return.Errors = false
	config.ClientID = "watermill"

	return config
}

type KafkaPublisherConfig struct {
	Brokers []string
	Retries int
}

func GetRuleType(ruleType string) rule.RuleType  {
	if ruleType == DROOL_RULE{
		return rule.DroolBusinessRule
	}
	if ruleType == MACHINE_LEARNING_RULE{
		return rule.MachineLearningRule
	}
	return rule.DroolBusinessRule
}


func ReadConfig() (*viper.Viper, error) {
	viper := viper.New()
	viper.SetConfigName("config")
	log.Println("")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType("yml")
	err := viper.ReadInConfig()
	return viper, err
}

func LoadConfiguration(appConfig *AppConfig, maps map[string]interface{}, config GetConfig) {
	config(appConfig, maps)
}

func InitializeGatewayConfiguration(viperConfig *viper.Viper) *AppConfig {
	appConfig := &AppConfig{}
	kafkaSubConfig := GetKafkaSubConfig()
	kafkaSubConfigMap := GenerateKafkaSubConfig(viperConfig)

	kafkaPubConfig := GetKafkaPubConfig()
	kafkaPubConfigMap := GenerateKafkaPubConfig(viperConfig)

	ruleConfig := GetRoutingRuleConfig()
	ruleConfigMap := GenerateRoutingRuleConfig(viperConfig)

	serverConfig := GetServerConfig()
	serverConfigMap := GenerateGetServerConfig(viperConfig)

	LoadConfiguration(appConfig,kafkaSubConfigMap,kafkaSubConfig)
	LoadConfiguration(appConfig,kafkaPubConfigMap,kafkaPubConfig)
	LoadConfiguration(appConfig,ruleConfigMap,ruleConfig)
	LoadConfiguration(appConfig,serverConfigMap,serverConfig)

	return appConfig
}
