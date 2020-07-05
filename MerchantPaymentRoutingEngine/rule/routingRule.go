package rule


import "bigdataconcept/fintech/intelligent/payment/routing/domain"

type RuleType int
const (
	MachineLearningRule RuleType = 1 << iota
	DroolBusinessRule
)

type RoutingRuleService interface {
	ExecutePaymentRoutingRule(merchantPayment *domain.MerchantPaymentRequest) (string, error)
}


func NewRoutingRuleService(rule RuleType) RoutingRuleService {
	switch rule {
	case DroolBusinessRule:
		return NewDroolRoutingRule("cardPaymentRoutingRule.grl")
	case MachineLearningRule:
		return NewMachineLearningRoutingRule("MLPModel path")
	default:
		return NewDroolRoutingRule("cardPaymentRoutingRule.grl")
	}
}




