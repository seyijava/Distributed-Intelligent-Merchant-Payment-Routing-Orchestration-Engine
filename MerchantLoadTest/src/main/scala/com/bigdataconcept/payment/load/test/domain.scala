package com.bigdataconcept.payment.load.test



object Domain{

       case class CardHolder(surname: String, name: String)
       case class Card(cardNumber: String,expiryDate: String,cardType: String,cardToken: String, cardHolder: CardHolder,IssuingBankCountry: String,cvc:String)
       case class BillingAddress(street: String,country: String,postalCode:String)
       case class Merchant(name: String, merchantCode: String,merchantAccount: String)
       case class PaymentRequest(merchantPaymentReference: String, card: Card, currency: String, amount: Double,billingAddress: BillingAddress,paymentDescription: String)
       case class MerchantPaymentRequest(paymentRequest: PaymentRequest, merchant: Merchant)
       case class PaymentResponse(transactionRef: String, responseCode: String, responseMessage: String)

}

