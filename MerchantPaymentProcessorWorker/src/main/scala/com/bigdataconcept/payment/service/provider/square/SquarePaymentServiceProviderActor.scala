package com.bigdataconcept.payment.service.provider.square

import akka.actor.{Actor, ActorLogging, Props}
import com.bigdataconcept.payment.domain.Configs.PaymentServiceProviderConfig
import com.bigdataconcept.payment.domain.PaymentDomain
import com.bigdataconcept.payment.domain.PaymentDomain.MerchantPaymentRequestCommand
import com.bigdataconcept.payment.service.provider.IPaymentServiceProvider



/**
 * @author Oluwaseyi Otun
 *         Square Payment Service Provider Implementation
 *         Payment Service Provider Actor. Any outbound payment for Square Service Provider is forwarded to this
 *         actor for processing
 */
object SquarePaymentServiceProviderActor{

   def props(config: PaymentServiceProviderConfig) : Props = Props(new SquarePaymentServiceProviderActor(config))
}


class SquarePaymentServiceProviderActor(config: PaymentServiceProviderConfig) extends Actor with ActorLogging with IPaymentServiceProvider{


  override def receive: Receive = {
    case cmd: MerchantPaymentRequestCommand => sendPaymentRequestToPaymentServiceProvider(cmd)
  }

  override def sendPaymentRequestToPaymentServiceProvider(merchantPaymentRequestCommand: PaymentDomain.MerchantPaymentRequestCommand): Unit={

  }
}
