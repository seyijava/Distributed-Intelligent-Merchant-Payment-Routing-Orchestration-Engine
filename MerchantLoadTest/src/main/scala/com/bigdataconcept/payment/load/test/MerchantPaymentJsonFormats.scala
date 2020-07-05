package com.bigdataconcept.payment.load.test

import akka.http.scaladsl.marshallers.sprayjson.SprayJsonSupport
import com.bigdataconcept.payment.load.test.Domain._
import spray.json.DefaultJsonProtocol
trait MerchantPaymentJsonFormats  extends  SprayJsonSupport with DefaultJsonProtocol{

  implicit val cardHolderForamt = jsonFormat2(CardHolder)
  implicit val cardForamt = jsonFormat7(Card)
  implicit val billingFormat = jsonFormat3(BillingAddress)
  implicit val merchantFormat = jsonFormat3(Merchant)
  implicit  val paymentRequest = jsonFormat6(PaymentRequest)
  implicit val merchantRequest = jsonFormat2(MerchantPaymentRequest)
  implicit val paymentResponse = jsonFormat3(PaymentResponse)
}
