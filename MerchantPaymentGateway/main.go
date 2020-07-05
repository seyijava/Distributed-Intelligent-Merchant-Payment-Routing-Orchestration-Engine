
package main

import (
	"bigdataconcept/fintech/intelligent/payment/routing/gateway/api"
	"bigdataconcept/fintech/intelligent/payment/routing/gateway/cmd"
	"bigdataconcept/fintech/intelligent/payment/routing/gateway/config"
	"bigdataconcept/fintech/intelligent/payment/routing/gateway/infracstructure"
	"context"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"time"
)

func main() {

	startMerchantGatewayService()

}



func startMerchantGatewayService() {
	if viperConfig, err := readConfig(); err != nil {
		panic(err)
	} else {
		appConfig := initializeGatewayConfiguration(viperConfig)
		var wait time.Duration = time.Second * 5
		kafkaPublisher := infracstructure.NewKafkaPublisher(appConfig.KafkaConfig)
		merchantGatewayService := api.NewMerchantGatewayService(appConfig.ServerConfig.ProcessorPoolSize, kafkaPublisher)
		defer merchantGatewayService.RequestProcessorPool.Release()
		router := mux.NewRouter()
		router.HandleFunc("/sendPayment", merchantGatewayService.SendMerchantPaymentRequest)
		srv := cmd.StartMerchantGateway(router, appConfig.ServerConfig)
		c := make(chan os.Signal, 1)
		// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
		// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
		signal.Notify(c, os.Interrupt)

		// Block until we receive our signal.
		<-c

		// Create a deadline to wait for.
		ctx, cancel := context.WithTimeout(context.Background(), wait)
		defer cancel()
		// Doesn't block if no connections, but will otherwise wait
		// until the timeout deadline.
		srv.Shutdown(ctx)
		// Optionally, you could run srv.Shutdown in a goroutine and block on
		// <-ctx.Done() if your application should wait for other services
		// to finalize based on context cancellation.
		log.Println("shutting down")

	}
}


func initializeGatewayConfiguration(viperConfig *viper.Viper) *config.AppConfig {
	appConfig := &config.AppConfig{}
    kafkaConfig := config.GetKafkaConfig()
    kafkaConfigMap := config.GenerateKafkaConfiguration(viperConfig)
    serverConfig := config.GetSeverConfig()
    severConfigMap := config.GenerateSeverConfig(viperConfig)
    loadConfiguration(appConfig,kafkaConfigMap,kafkaConfig)
    loadConfiguration(appConfig,severConfigMap,serverConfig)
	return appConfig
}



func readConfig() (*viper.Viper, error) {
	viper := viper.New()
	viper.SetConfigName("config")
	log.Println("")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType("yml")
	err := viper.ReadInConfig()
	return viper, err
}

func loadConfiguration(appConfig *config.AppConfig, maps map[string]interface{}, config config.GetConfig) {
	config(appConfig, maps)
}

