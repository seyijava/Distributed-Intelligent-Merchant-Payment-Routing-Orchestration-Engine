package com.bigdataconcept.payment.service.provider.paypal;

import com.bigdataconcept.payment.domain.Configs.PaymentServiceProviderConfig;
import  com.braintreegateway.BraintreeGateway;
import com.braintreegateway.Environment;

public class Gateway {


      private static BraintreeGateway gateway = null;

      public static BraintreeGateway getGateWayInstance(PaymentServiceProviderConfig config) {
           if(gateway == null){
             gateway = new  BraintreeGateway(Environment.SANDBOX, config.merchantId(),config.publicKey(),config.privateKey());
           }
          return gateway;
      }

}
