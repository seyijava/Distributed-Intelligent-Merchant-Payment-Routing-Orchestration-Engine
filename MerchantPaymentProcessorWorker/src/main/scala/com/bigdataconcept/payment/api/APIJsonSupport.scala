package com.bigdataconcept.payment.api

import akka.http.scaladsl.marshallers.sprayjson.SprayJsonSupport
import com.bigdataconcept.payment.domain.Adyen.{Amount, PaymentResponse}

class APIJsonSupport extends  SprayJsonSupport{
  import spray.json.DefaultJsonProtocol._

  implicit val amount = jsonFormat2(Amount)
  implicit val response = jsonFormat4(PaymentResponse)

}
