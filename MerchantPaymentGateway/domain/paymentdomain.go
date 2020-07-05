package domain

type PaymentResponse struct {
	TransactionRef  string `json:"transactionRef"`
	ResponseCode    string `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
}
type CardHolder struct {
	Surname string `json:"surname"`
	Name    string `json:"name"`
}

type PaymentRequest struct {
	Card                     *Card           `json:"card"`
	Currency                 string          `json:"currency"`
	Amount                   float64         `json:"amount"`
	BillingAddress           *BillingAddress `json:"billingAddress"`
	MerchantPaymentReference string          `json:"merchantPaymentReference"`
	PaymentDescription       string          `json:"paymentDescription"`
}

type BillingAddress struct {
	Street     string `json:"street"`
	Address    string `json:"address"`
	Country    string `json:"country"`
	PostalCode string `json:"postalCode"`
}

type Card struct {
	CardNumber         string     `json:"cardNumber"`
	ExpiryDate         string     `json:"expiryDate"`
	Cardtype           string     `json:"cardtype"`
	CardToken          string     `json:"cardToken"`
	CardHolder         CardHolder `json:"cardHolder"`
	IssuingBankCountry string     `json:"issuingBankCountry"`
	CVC                string     `json:"cvc"`
}

type Merchant struct {
	Name            string `json:"name"`
	MerchantAccount string `json:"merchantAccount"`
	MerchantCode    string `json:"merchantCode"`
}

type MerchantPaymentRequest struct {
	*Merchant       `json:"merchant"`
	*PaymentRequest `json:"paymentRequest"`
}

type MerchantPaymentResponse struct {
	MerchantPaymentRef string `json:"merchantPaymentReference"`
	TransactionRef     string `json:"transactionRef"`
	ResponseCode       string `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
}
