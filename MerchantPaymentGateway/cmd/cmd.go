package cmd

import (
	"bigdataconcept/fintech/intelligent/payment/routing/gateway/config"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

func StartMerchantGateway(router *mux.Router, serverConfig config.ServerConfig) *http.Server  {
	srv := &http.Server{
		Addr:         serverConfig.IP + ":" + strconv.Itoa(serverConfig.Port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler: router, // Pass our instance of gorilla/mux in.
	}
	// Run our server in a goroutine so that it doesn't block.
	go func() {
		log.Infof("Merchant Gateway Service Started on IP %s @  Port %v",serverConfig.IP,serverConfig.Port)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
return srv

}
