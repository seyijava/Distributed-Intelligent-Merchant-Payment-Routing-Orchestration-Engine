package main

import (
	"bigdataconcept/fintech/intelligent/payment/routing/config"
	"bigdataconcept/fintech/intelligent/payment/routing/service"

	"github.com/gorilla/mux"
	"net/http"

	"context"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	paymentRoutingEngine, serverConfig := initPaymentRoutingEngine(ctx)
	paymentRoutingEngine.Start(ctx, wg)

	go startServer(serverConfig)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		log.Println("Terminating: Context Cancelled")
	case <-sigterm:
		log.Println("Terminating: via signal")
	}
	cancel()
	wg.Wait()
	if err := paymentRoutingEngine.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
}

func initPaymentRoutingEngine(ctx context.Context) (*service.PaymentRoutingEngine, *config.ServerConfig) {
	if viperConfig, err := config.ReadConfig(); err != nil {
		panic(err)
	} else {
		appConfig := config.InitializeGatewayConfiguration(viperConfig)
		publishKafkaConfig := appConfig.KafkaPublisherConfig
		subKafkaConfig := appConfig.KafkaSubscriberConfig
		ruleConfig := appConfig.RoutingRuleConfig
		routingRule := config.GetRuleType(ruleConfig.RoutingRuleType)
		paymentRoutingService := service.NewPaymentRoutingEngine(publishKafkaConfig, subKafkaConfig, routingRule)
		return paymentRoutingService, &appConfig.ServerConfig
	}
	return nil, nil
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func startServer(serverConfig *config.ServerConfig) {
	router := mux.NewRouter()
	router.HandleFunc("/health", healthCheck)
	http.ListenAndServe(serverConfig.IP+":"+serverConfig.Port, router)
}
