package com.bigdataconcept.payment.load.test

import java.util.Base64
import akka.Done
import akka.actor.{Actor, ActorLogging, Props, Status}
import akka.http.scaladsl.Http
import akka.http.scaladsl.marshalling.Marshal
import akka.stream.Materializer
import com.bigdataconcept.payment.load.test
import akka.http.scaladsl.model.{ContentType, HttpEntity, HttpMethods, HttpRequest, MediaTypes, MessageEntity}
import akka.http.scaladsl.unmarshalling.Unmarshal
import io.alphash.faker._
import com.bigdataconcept.payment.load.test.Domain._
import com.github.javafaker.Faker
import com.github.javafaker.Name
import com.google.gson.Gson
import akka.pattern.pipe
import com.typesafe.config.ConfigFactory

import scala.concurrent.duration._
import scala.concurrent.{ExecutionContext, Future}
import scala.util.Random


object Simulator{
     object Start
     object Stop
     def props(ip: String,targetPorts: Seq[Int])(implicit mat: Materializer, ec: ExecutionContext) : Props =  Props(new Simulator(ip,targetPorts))
}
import scala.concurrent.ExecutionContext

class Simulator(ip: String,targetPorts: Seq[Int]) (implicit mat: Materializer, ec: ExecutionContext) extends Actor with MerchantPaymentJsonFormats with ActorLogging{

 import Simulator._

  val merchantList = Seq("Amazon:12345768728:AMZ","AliExpress:72345768728:ALI", "Walmart:82345768728:WAL","Cosco:92345768728:COSC")
  val countryList = Seq("CA","US","CH","NG","GH","IN","SG")
  val random = new Random
  private def port: Int = Random.shuffle(targetPorts).head
  private def countryCode: String = Random.shuffle(countryList).head
  private val http = Http(context.system)

  private var startTime: Long = System.currentTimeMillis()


  context.self ! Start

  override def receive: Receive = {
    case Start =>
    startTime = System.currentTimeMillis()
      log.info("Start")
    run().pipeTo(self)
    case paymentResponse: PaymentResponse =>
      log.info(s"Response ${paymentResponse} in ${System.currentTimeMillis() - startTime} ms")
      context.self ! Start
    case Status.Failure(ex) =>
      log.info(s"Failed: ${ex.getMessage}")
  }

  private def run(): Future[PaymentResponse] = {
    for {
      paymentResponse <-  sendMerchantPayment()
    } yield{
      paymentResponse
    }
  }
  private def sendMerchantPayment(): Future[PaymentResponse] = {
 for {
      requestEntity <- Marshal(buildMerchantPaymentRequest()).to[MessageEntity]
      request = HttpRequest(HttpMethods.POST, s"http://$ip:$port/sendPayment", entity = requestEntity.withContentType(ContentType(MediaTypes.`application/json`)))
      response <- http.singleRequest(request)
      responseEntity <- response.entity.toStrict(5.seconds)
       reps <- Unmarshal(responseEntity).to[PaymentResponse]
 } yield {
   reps
    }}

   private def buildMerchantPaymentRequest() : MerchantPaymentRequest={
     val faker = new Faker()
     val payment = Payment()
     val person = Person()
     val creditCardNumber = payment.creditCardNumber();
     val name = person.firstNameMale
     val surname = person.lastName
     val cardHolder = CardHolder(surname,name)
     val cvc =generateRandomNumber(3)
     val token = creditCardNumber + cvc
     val tokenEnc = new String(Base64.getEncoder.encode(token.getBytes))
     val card = Card(creditCardNumber,"05/2018",payment.creditCardType,tokenEnc,cardHolder,countryCode,cvc)
     val address = faker.address()
     val billingAddress = BillingAddress(address.streetAddress(),address.countryCode(),address.zipCode())
     val merchantCode = "MCT-" + generateRandomNumber(3)
     val merchantDetails = merchantList(random.nextInt(merchantList.length)).split(":")
     val merchant = Merchant(merchantDetails(0),merchantCode,merchantDetails(1))
     val merchantReferenceNumber = merchantDetails(2) + "/" + generateRandomNumber(6)
     val paymentRequest = new PaymentRequest(merchantReferenceNumber,card,"USD",8.0,billingAddress,"abc")
     val merchantPaymentRequest = MerchantPaymentRequest(paymentRequest,merchant)

     return merchantPaymentRequest
   }


    def generateRandomNumber(length: Integer) : String ={
      val  random = new Random()
      val number = random.nextInt(100000000)
       val num = String.valueOf(number)

      return num.substring(0,num.length)
    }



}
