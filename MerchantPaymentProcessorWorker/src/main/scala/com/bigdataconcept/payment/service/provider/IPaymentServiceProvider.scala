package com.bigdataconcept.payment.service.provider

import com.bigdataconcept.payment.domain.PaymentDomain.MerchantPaymentRequestCommand

trait IPaymentServiceProvider {

  def sendPaymentRequestToPaymentServiceProvider(merchantPaymentRequestCommand: MerchantPaymentRequestCommand) :Unit
}
