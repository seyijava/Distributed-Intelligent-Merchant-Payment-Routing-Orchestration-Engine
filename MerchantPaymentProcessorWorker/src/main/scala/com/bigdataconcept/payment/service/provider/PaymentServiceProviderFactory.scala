package com.bigdataconcept.payment.service.provider

import akka.actor.Props
import com.bigdataconcept.payment.domain.Configs.PaymentServiceProviderConfig
import com.bigdataconcept.payment.service.provider.adyen.AdyenPaymentServiceProviderActor
import com.bigdataconcept.payment.service.provider.paypal.PayPalBraintreePaymentServiceProviderActor
import com.bigdataconcept.payment.service.provider.stripe.StripePaymentServiceProviderActor


object PaymentServiceProvider extends Enumeration {
  type PaymentServiceProvider = Value
  val PAYPAL, STRIPE,ADYEN  = Value
}


/**
 * @author Oluwaseyi Otun
 *         PaymentServiceProviderFactory is a factory class create Payment Service Provider
 *         base on the configuration and instance that is being running on cluster.
 */

object PaymentServiceProviderFactory {

  import PaymentServiceProvider._

  def apply(psp: PaymentServiceProvider, config: PaymentServiceProviderConfig) :Props = psp match {

    case PAYPAL => PayPalBraintreePaymentServiceProviderActor.props(config)

    case STRIPE => StripePaymentServiceProviderActor.props(config)

    case ADYEN => AdyenPaymentServiceProviderActor.props(config)


  }

}
