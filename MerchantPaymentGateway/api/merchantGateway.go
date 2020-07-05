package api

import (
	"bigdataconcept/fintech/intelligent/payment/routing/gateway/domain"
	"bigdataconcept/fintech/intelligent/payment/routing/gateway/infracstructure"
	"encoding/json"
	"fmt"
	"github.com/panjf2000/ants/v2"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type RequestResponse struct {
	incomingRequest *domain.MerchantPaymentRequest
	outgoingResponse        chan *domain.MerchantPaymentResponse
}

type MerchantGatewayService struct {
	RequestProcessorPool *ants.PoolWithFunc
}

func NewMerchantGatewayService(poolSize int, publisher *infracstructure.KafkaPublisher) *MerchantGatewayService {
	pool, _ := ants.NewPoolWithFunc(poolSize, func(payload interface{}) {
		requestResponse, ok := payload.(*RequestResponse)
		if !ok {
			return
		}
		if err, reply := infracstructure.ProcessRequest(requestResponse.incomingRequest, publisher); err != nil {
			log.Error(err)
			requestResponse.outgoingResponse <- &domain.MerchantPaymentResponse{ResponseMessage: "Internal Server Error Try again",ResponseCode: "99",TransactionRef: "XXXXXXXXXX"}
		} else {
			requestResponse.outgoingResponse <- reply
		}
	})
	return &MerchantGatewayService{RequestProcessorPool: pool}
}

func (merchantGatewayService *MerchantGatewayService) SendMerchantPaymentRequest(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	fmt.Sprint(string(requestBody))
	if err != nil {
		log.Error(err)
		http.Error(w, "request error", http.StatusInternalServerError)
	}
	var merchantPayment = &domain.MerchantPaymentRequest{}
	if err := json.Unmarshal(requestBody, merchantPayment); err != nil {
		log.Error(err)
		http.Error(w, "request error", http.StatusInternalServerError)
	}
	defer r.Body.Close()
	requestResponse := &RequestResponse{incomingRequest: merchantPayment, outgoingResponse: make(chan *domain.MerchantPaymentResponse)}
	if err := merchantGatewayService.RequestProcessorPool.Invoke(requestResponse); err != nil {
		log.Error(err)
		http.Error(w, "throttle limit error", http.StatusInternalServerError)
		return
	}
	var response = <-requestResponse.outgoingResponse
	 responseMsg, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseMsg)
}
