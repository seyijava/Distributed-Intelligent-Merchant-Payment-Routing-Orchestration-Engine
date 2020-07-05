package com.bigdataconcept.payment.domain


object PaymentDomain {

   case class CardHolder(surname:String, name:String, email: String)
   case class Card(cardNumber: String, expiryDate: String,expiryMonth: String, cvc: String,cardHolder: CardHolder)
   case class Currency(code: String)
   case class BillingAddress(street: String, postalCode: String,address: String,country: String)
   case class Money(amount: Double)
   case class PaymentRequest(card: Card,currency:String, amount: Double, billingAddress: BillingAddress, merchantPaymentReference: String,paymentDescription:String)
   case class Merchant(name: String, merchantAccount: String, merchantCode: String)
   case class MerchantPaymentRequestCommand(paymentRequest: PaymentRequest, merchant: Merchant)

   case class PaymentResponseMessage(responseCode: String, responseMessage: String, paymentReference: String)
}




object Configs{
    case class PaymentServiceProviderConfig(serviceUrl: String, userName: String,
                                             password: String, merchantAccount: String, merchantId: String, publicKey: String, privateKey: String)

      /* def generateEncodedBasicAuthentication() : String={
         val authString = userName + ":" + password
         val authEncBytes = Base64.encodeBase64(authString.getBytes)
         val authStringEnc = new String(authEncBytes)
           return authStringEnc
        }*/

}


object Adyen{

   case class Card(cardNumber: String, cardHolder: String,  expiryMonth: String, expiryYear: String,cvc: String)
   case class CardHolder(name: String)
   case class PaymentMethod(ty: String,number: String,expiryMonth: String, expiryYear:String,holderName:String, cvc: String)
   case class Amount(value: Double, currency: String)
   case class PaymentRequest(card: Card, amount: Amount, cardHolder: CardHolder,merchantAccount: String,reference: String)
   case class PaymentResponse(pspReference: String, resultCode: String, merchantReference: String, amount: Amount)
}


object PayPal{
   case class Amount(value: Double, currency: String)
   case class PayPalPaymentRequest(invoice_id: String, amount: Amount, final_capture: Boolean)
   case class Link(rel: String, method: String, href: String)
   case class PayPalPaymentResponse(id: String, status: String,links: Array[Link])
}


object Stripe{

   case class Card(brand: String, exp_month: String,exp_year: String,number: String)
   case class Address(city: String, country: String, line1: String, line2: String, postal_code: String, state: String)
   case class BillingDetails(email: String, name: String, address: Address)
   case class PaymentMethod(card :Card, billingDetails: BillingDetails)
   case class Amount(value: Double, currency: String)
   case class PaymentRequest(paymentMethod: PaymentMethod, amount: Amount, merchantAccount: String)
}

object Square{
     case class app_fee_money(amount: Double, currency: String)
     case class amount_money(amount: Double, currency: String)
     case class PaymentRequest(idempotency_key: String, source_id: String, autocomplete: Boolean,customer_id: String,location_id: String,
                               reference_id: String,note: String, app_fee_money: app_fee_money,amount_money: amount_money)

}