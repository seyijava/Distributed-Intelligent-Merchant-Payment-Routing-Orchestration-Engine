package rule



//Machine Learning Service
//We can invoke the service model directly or call external API for the ML

import (

	"bigdataconcept/fintech/intelligent/payment/routing/domain"
)

type MachineLearningRoutingRule struct {

}


func NewMachineLearningRoutingRule(modelPath string) *MachineLearningRoutingRule{
	  return &MachineLearningRoutingRule{}
}

func (mlRoutingRule *MachineLearningRoutingRule) ExecutePaymentRoutingRule(merchantPayment *domain.MerchantPaymentRequest) (string,error) {

	//call to Machine Learning model to evaluate the the routing payment service provider to route
	// to base on scoring model this can be through api to an external ML Services

	return "", nil
}
