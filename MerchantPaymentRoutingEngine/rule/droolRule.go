package rule

import (
	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	_ "github.com/hyperjumptech/grule-rule-engine/engine"
	log "github.com/sirupsen/logrus"

	//log "github.com/sirupsen/logrus"
	"bigdataconcept/fintech/intelligent/payment/routing/domain"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
	//"strings"
)


type CardRuleFacts struct {
	 CardNumber string
	 CardType string
     PSPChannel string
	 CardBin string
	IssueCountry string


}
type DroolRoutingRule struct {
	workingMemory *ast.WorkingMemory
    RuleConfigurationFile string
	knowledgeBase *ast.KnowledgeBase
	ruleDataContext ast.DataContext
}


func NewDroolRoutingRule(ruleConfigurationPath string) *DroolRoutingRule{
	droolRouting := &DroolRoutingRule{}
	lib := ast.NewKnowledgeLibrary()
	ruleBuilder := builder.NewRuleBuilder(lib)
	ruleFile := pkg.NewFileResource(ruleConfigurationPath)
	err := ruleBuilder.BuildRuleFromResource("CardPaymentRoutingRule", "0.0.1" ,ruleFile)
	if err != nil {
		panic(err)
	}
	knowledgeBase := lib.NewKnowledgeBaseInstance("CardPaymentRoutingRule", "0.0.1")
	droolRouting.knowledgeBase = knowledgeBase
    return droolRouting
}

func (droolRoutingRule *DroolRoutingRule) ExecutePaymentRoutingRule(merchantPayment *domain.MerchantPaymentRequest) (string,error)  {
	cardBin := merchantPayment.Card.CardNumber[0:6]
	cardRuleFacts := &CardRuleFacts{CardNumber: merchantPayment.Card.CardNumber,CardType: merchantPayment.Card.CardType,PSPChannel: "",CardBin: cardBin,IssueCountry: merchantPayment.Card.IssuingBankCountry}
	dataCtx := ast.NewDataContext()
	err := dataCtx.Add("CardInfo", cardRuleFacts)
	if err != nil {
		panic(err)
	}
	ruleEngine := engine.NewGruleEngine()

	if err := ruleEngine.Execute(dataCtx,droolRoutingRule.knowledgeBase); err != nil{
		panic(err)
		return "",err
	}else{
		log.Infof("Payment Service Provider Channel Topic [%s]" ,cardRuleFacts.PSPChannel)
		return cardRuleFacts.PSPChannel, nil
	}
}