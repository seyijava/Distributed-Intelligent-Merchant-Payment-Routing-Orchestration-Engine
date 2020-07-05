package com.bigdataconcept.payment.service.provider.adyen

import java.io.IOException

import akka.actor.{Actor, ActorLogging, Props}
import com.adyen.Client
import com.adyen.enums.Environment
import com.adyen.model.{Amount, Card, PaymentRequest}
import com.adyen.service.Payment
import com.adyen.service.exception.ApiException
import com.bigdataconcept.payment.domain.Configs.PaymentServiceProviderConfig
import com.bigdataconcept.payment.domain.PaymentDomain.{MerchantPaymentRequestCommand, PaymentResponseMessage}
import com.bigdataconcept.payment.service.provider.IPaymentServiceProvider

import scala.concurrent.Future
import scala.util.{Failure, Success}

/**
 * @author Oluwaseyi Otun
 *         Adyen Payment Service Provider Implementation
 *         Payment Service Provider Actor. Any outbound payment for Adyen Service Provider is forwarded to this
 *         actor for processing
 */
object AdyenPaymentServiceProviderActor {

  def props(config: PaymentServiceProviderConfig): Props = Props(new AdyenPaymentServiceProviderActor(config))


}

class AdyenPaymentServiceProviderActor(config: PaymentServiceProviderConfig) extends Actor with ActorLogging with IPaymentServiceProvider {

  val IDEMPOTENCY_KEY = "Idempotency-Key"

  implicit val executionContext = context.system.dispatcher
  var payment: Payment = null

  implicit val sys = context.system

  override def preStart() : Unit={
    val client = new Client(config.userName, config.password, Environment.TEST, "Test")
     payment = new Payment(client)
  }

  override def receive: Receive = {
    case merchantPaymentRequestCommand: MerchantPaymentRequestCommand => sendPaymentRequestToPaymentServiceProvider(merchantPaymentRequestCommand)
  }


  override def sendPaymentRequestToPaymentServiceProvider(merchantPaymentRequestCommand: MerchantPaymentRequestCommand): Unit = {
    log.info("Sending Payment To Adyen Payment Reference {}", merchantPaymentRequestCommand.paymentRequest.merchantPaymentReference)

    val paymentResponse = makeAPIAdyenCall(merchantPaymentRequestCommand)
    paymentResponse.mapTo[PaymentResponseMessage].onComplete {
      case Success(response) => log.info(response.responseCode)
      case Failure(ex) => log.info(ex.getMessage)
    }
  }


  def makeAPIAdyenCall(merchantPaymentRequestCommand: MerchantPaymentRequestCommand): Future[PaymentResponseMessage] = {
    var responseGateway: Future[PaymentResponseMessage] = null
    val paymentRequest = new PaymentRequest()
    val requestCard = merchantPaymentRequestCommand.paymentRequest.card
    val card = new Card()
    card.setCvc(requestCard.cvc)
    card.setExpiryMonth(requestCard.expiryMonth)
    card.setExpiryYear(requestCard.expiryDate)
    card.setNumber(requestCard.cardNumber)
    val cardHolder = requestCard.cardHolder.name + " " + requestCard.cardHolder.surname
    card.setHolderName(cardHolder)
    paymentRequest.setCard(card)
    val amt = new Amount()
    amt.setCurrency(merchantPaymentRequestCommand.paymentRequest.currency)
    amt.setValue(merchantPaymentRequestCommand.paymentRequest.amount.toLong)
    paymentRequest.setAmount(amt)
    try {
      val paymentResult = payment.authorise(paymentRequest)
      responseGateway = Future.successful(new PaymentResponseMessage("00", "Successful" ,paymentResult.getAuthCode))
    }
    catch {
      case ex: ApiException =>  responseGateway = Future.successful(new PaymentResponseMessage("99", ex.getMessage, String.valueOf(ex.getStatusCode)))
      case io: IOException =>   responseGateway = Future.successful(new PaymentResponseMessage("99", io.getMessage, ""))
    }
    return responseGateway
  }

}
