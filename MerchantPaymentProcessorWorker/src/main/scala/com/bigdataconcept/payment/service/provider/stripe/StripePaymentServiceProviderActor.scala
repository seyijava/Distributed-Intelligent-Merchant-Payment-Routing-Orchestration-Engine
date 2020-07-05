package com.bigdataconcept.payment.service.provider.stripe


import akka.actor.{Actor, ActorLogging, Props}
import akka.stream.ActorMaterializer
import com.bigdataconcept.payment.domain.Configs.PaymentServiceProviderConfig
import com.bigdataconcept.payment.domain.PaymentDomain.{MerchantPaymentRequestCommand, PaymentResponseMessage}
import com.bigdataconcept.payment.service.provider.IPaymentServiceProvider
import com.stripe.Stripe
import com.stripe.exception.StripeException
import com.stripe.model.{Charge, Token}

import scala.concurrent.Future
import scala.util.{Failure, Success}


object StripePaymentServiceProviderActor {
  def props(config: PaymentServiceProviderConfig): Props = Props(new StripePaymentServiceProviderActor(config))
}


/**
 * @author Oluwaseyi Otun
 *         Stripe Payment Service Provider Implementation
 *         Payment Service Provider Actor. Any outbound payment for Stripe Service Provider is forwarded to this
 *         actor for processing
 */

class StripePaymentServiceProviderActor(config: PaymentServiceProviderConfig) extends Actor with ActorLogging with IPaymentServiceProvider {

  implicit val materializer = ActorMaterializer()

  implicit val executionContext = context.system.dispatcher

  implicit val sys = context.system
  override def receive: Receive = {
    case cmd: MerchantPaymentRequestCommand => sendPaymentRequestToPaymentServiceProvider(cmd)
  }

  override def sendPaymentRequestToPaymentServiceProvider(merchantPaymentRequestCommand: MerchantPaymentRequestCommand): Unit = {

    log.info("Sending Payment To Stripe Payment Reference {}", merchantPaymentRequestCommand.paymentRequest.merchantPaymentReference)

     val gatewayResponse = sendAPICallToStripe(merchantPaymentRequestCommand)
     gatewayResponse.mapTo[PaymentResponseMessage].onComplete{
       case Success(response) => log.info(response.responseCode)
       case Failure(ex) => log.info(ex.getLocalizedMessage)
     }

  }


  def sendAPICallToStripe(merchantPaymentRequestCommand: MerchantPaymentRequestCommand): Future[PaymentResponseMessage] = {
    import java.util

    var responseGateway: Future[PaymentResponseMessage] = null
    try {
      Stripe.apiKey = config.publicKey
      val card = new util.HashMap[String, AnyRef]
      card.put("number", merchantPaymentRequestCommand.paymentRequest.card.cardNumber)
      card.put("exp_month", merchantPaymentRequestCommand.paymentRequest.card.expiryMonth)
      card.put("exp_year", merchantPaymentRequestCommand.paymentRequest.card.expiryDate)
      val params = new util.HashMap[String, AnyRef]
      params.put("card", card.asInstanceOf[AnyRef])
      val token = Token.create(params)
      val chargeParams = new util.HashMap[String, AnyRef]
      chargeParams.put("amount", merchantPaymentRequestCommand.paymentRequest.amount.asInstanceOf[AnyRef])
      chargeParams.put("currency", merchantPaymentRequestCommand.paymentRequest.currency.asInstanceOf[AnyRef])
      chargeParams.put("description", merchantPaymentRequestCommand.paymentRequest.paymentDescription.asInstanceOf[AnyRef])
      chargeParams.put("source", token.getId)
      val charge = Charge.create(chargeParams)
      responseGateway = Future.successful(PaymentResponseMessage(responseCode = charge.getStatus,paymentReference = merchantPaymentRequestCommand.paymentRequest.merchantPaymentReference,responseMessage = ""))
    } catch {
      case e: StripeException =>
        responseGateway = Future.successful(PaymentResponseMessage(responseCode = "99",paymentReference = merchantPaymentRequestCommand.paymentRequest.merchantPaymentReference,responseMessage = ""))
    }
    return responseGateway
  }


}
