package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strconv"
	"strings"
)

const (
	MONGO_CONNECTIONURL = "MONGO_CONNECTIONURL"
	MONGO_DB ="MONGO_DB"
	MONGO_COLLECTION = "MONGO_COLLECTION"
	MONGO_USERNAME = "MONGO_USERNAME"
	MONGO_PASSWORD = "MONGO_PASSWORD"


	SERVER_HOST = "SERVER_HOST"
	SERVER_PORT = "SERVER_PORT"
	SERVER_GOROUTINEPOOLSIZE = "SERVER_GOROUTINEPOOLSIZE"

	KAFKA_BROKERS = "KAFKA_BROKERS"
	KAFKA_TOPIC = "KAFKA_TOPIC"
	KAFKA_RETRY = "KAFKA_RETRY"
)


type GetConfig func(*AppConfig,  map[string]interface{}) *AppConfig

func GetKafkaConfig() GetConfig {
	return  func(config *AppConfig, configVal map[string]interface{}) *AppConfig {
		config.KafkaConfig.BrokerAddress = strings.Split(configVal[KAFKA_BROKERS].(string), ",")
		config.KafkaConfig.Topic = configVal[KAFKA_TOPIC].(string)
		retry,_ :=  strconv.Atoi(configVal[KAFKA_RETRY].(string))
		config.KafkaConfig.Retry = retry
		return config
	}
}

func GenerateKafkaConfiguration(v *viper.Viper) map[string]interface{}{
	configurationMap := make(map[string]interface{})
	if brokerUrl := v.Get(KAFKA_BROKERS); brokerUrl == nil{
		configurationMap[KAFKA_BROKERS] = v.Get("kafkaConfig.brokerUrl")
	}else{
		log.Info(brokerUrl)
		configurationMap[KAFKA_BROKERS] = brokerUrl
	}
	if topic := v.Get(KAFKA_TOPIC); topic == nil{
		configurationMap[KAFKA_TOPIC] = v.Get("kafkaConfig.topic")
	}else{
		configurationMap[KAFKA_TOPIC] = topic
	}
	if retry := v.Get(KAFKA_RETRY); retry == nil{
		configurationMap[KAFKA_RETRY] = v.Get("kafkaConfig.retry")
	}else{
		configurationMap[KAFKA_RETRY] = retry
	}

	return configurationMap
}

func GetSeverConfig() GetConfig {
	return  func(config *AppConfig, configVal map[string]interface{}) *AppConfig {
		config.ServerConfig.IP = configVal[SERVER_HOST].(string)
		log.Info("" + configVal[SERVER_PORT].(string))
		port,_ :=  strconv.Atoi(configVal[SERVER_PORT].(string))
		config.ServerConfig.Port = port
		poolSize,_ :=  strconv.Atoi(configVal[SERVER_GOROUTINEPOOLSIZE].(string))
		config.ServerConfig.ProcessorPoolSize = poolSize
		return config
	}
}


func GenerateSeverConfig(v *viper.Viper) map[string]interface{}{
	configurationMap := make(map[string]interface{})
	if severHost := v.Get(SERVER_HOST); severHost == nil{
		configurationMap[SERVER_HOST] = v.Get("serverConfig.host")
	}else{
		configurationMap[SERVER_HOST] = severHost
	}
	if port := v.Get(SERVER_PORT); port == nil{
		configurationMap[SERVER_PORT] = v.Get("serverConfig.port")
	}else{
		configurationMap[SERVER_PORT] = port
	}
	if poolSize := v.Get(SERVER_GOROUTINEPOOLSIZE); poolSize == nil{
		configurationMap[SERVER_GOROUTINEPOOLSIZE] = v.Get("serverConfig.goroutinePoolSize")
	}else{
		configurationMap[SERVER_GOROUTINEPOOLSIZE] = poolSize
	}

	return configurationMap
}


type ServerConfig struct {
	IP string
	Port int
	ProcessorPoolSize int
}


type KafkaConfig struct {
	 BrokerAddress []string
	 Topic string
	 Retry int
}


type MongoDBConfig struct {
	 DataBase string
	 Collection string
	 Username string
	 Password string
	 ConnectionUrl string
}


type AppConfig struct {
	 KafkaConfig
	 ServerConfig
	 MongoDBConfig
}
