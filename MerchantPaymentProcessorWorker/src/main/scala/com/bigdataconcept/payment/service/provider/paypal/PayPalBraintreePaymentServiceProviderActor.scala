package com.bigdataconcept.payment.service.provider.paypal

import akka.actor.{Actor, ActorLogging, Props}
import akka.stream.ActorMaterializer
import com.bigdataconcept.payment.domain.Configs.PaymentServiceProviderConfig
import com.bigdataconcept.payment.domain.PaymentDomain
import com.bigdataconcept.payment.domain.PaymentDomain.{MerchantPaymentRequestCommand, PaymentResponseMessage}
import com.bigdataconcept.payment.service.provider.IPaymentServiceProvider
import com.braintreegateway.PaymentMethodRequest
import com.braintreegateway.exceptions.{AuthenticationException, AuthorizationException}

import scala.concurrent.Future
import scala.util.{Failure, Success}

/**
 * @author Oluwaseyi Otun
 *         PayPal Braintree Payment Service Provider Implementation
 *         Payment Service Provider Actor. Any outbound payment for PayPal Braintree Service Provider is forwarded to this
 *         actor for processing
 */

object PayPalBraintreePaymentServiceProviderActor{
      def props( config: PaymentServiceProviderConfig) : Props = Props(new PayPalBraintreePaymentServiceProviderActor(config))
}
class PayPalBraintreePaymentServiceProviderActor(config: PaymentServiceProviderConfig)  extends  Actor with ActorLogging with IPaymentServiceProvider {


  implicit val materializer = ActorMaterializer()

  implicit val executionContext = context.system.dispatcher

  implicit val sys = context.system

  override def receive: Receive = {
    case cmd: MerchantPaymentRequestCommand => sendPaymentRequestToPaymentServiceProvider(cmd)
  }

  override def sendPaymentRequestToPaymentServiceProvider(merchantPaymentRequestCommand: PaymentDomain.MerchantPaymentRequestCommand): Unit = {
    log.info("Sending Payment To PayPal Payment Reference {}", merchantPaymentRequestCommand.paymentRequest.merchantPaymentReference)
    val payPalPaymentResponse = makeAPIToPayPal(merchantPaymentRequestCommand)
    payPalPaymentResponse.mapTo[PaymentResponseMessage].onComplete {
      case Success(response) => log.info("Response Message [{}]" , response)
      case Failure(ex) =>
    }
  }


  def makeAPIToPayPal(merchantPaymentRequestCmd: PaymentDomain.MerchantPaymentRequestCommand): Future[PaymentResponseMessage] = {
    log.info("Sending Payment To BrainTree Gateway Payment Reference {}", merchantPaymentRequestCmd.paymentRequest.merchantPaymentReference)
    var responseGateway: Future[PaymentResponseMessage]  = null
    val merchantPaymentReference = merchantPaymentRequestCmd.paymentRequest.merchantPaymentReference
    val request = new PaymentMethodRequest()
    request.cvv(merchantPaymentRequestCmd.paymentRequest.card.cvc)
      .cardholderName(merchantPaymentRequestCmd.paymentRequest.card.cardHolder.name + " " + merchantPaymentRequestCmd.paymentRequest.card.cardHolder.surname)
      .expirationMonth(merchantPaymentRequestCmd.paymentRequest.card.expiryMonth)
      .expirationYear(merchantPaymentRequestCmd.paymentRequest.card.expiryDate)
      .number(merchantPaymentRequestCmd.paymentRequest.card.cardNumber)
   try {
     val response = Gateway.getGateWayInstance(config).paymentMethod().create(request)
      responseGateway = Future.successful(PaymentResponseMessage(responseCode = "00",paymentReference = merchantPaymentReference,responseMessage = ""))
     return responseGateway
   }catch{
     case x: AuthenticationException =>
       responseGateway = Future.successful(PaymentResponseMessage(responseCode = "99",paymentReference = merchantPaymentReference,responseMessage = "Authentication Failed"))
       return responseGateway
     case x: AuthorizationException =>
       responseGateway = Future.successful(PaymentResponseMessage(responseCode = "99",paymentReference = merchantPaymentReference,responseMessage = "Authorization Failed"))
       return responseGateway
     case _ =>
       responseGateway = Future.successful(PaymentResponseMessage(responseCode = "99",paymentReference = merchantPaymentReference,responseMessage = "Failed"))
       return responseGateway
   }

  }


}