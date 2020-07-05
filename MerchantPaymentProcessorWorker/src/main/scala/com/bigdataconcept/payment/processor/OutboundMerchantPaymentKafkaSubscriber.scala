package com.bigdataconcept.payment.processor

import akka.Done
import akka.actor.{Actor, ActorLogging, ActorRef, Props}
import akka.kafka.scaladsl.Consumer
import akka.kafka.{ConsumerSettings, Subscriptions}
import akka.stream.ActorMaterializer
import akka.stream.scaladsl.Sink
import com.bigdataconcept.payment.domain.PaymentDomain.MerchantPaymentRequestCommand
import com.google.gson.Gson
import org.apache.kafka.clients.consumer.ConsumerRecord
import org.apache.kafka.common.serialization.StringDeserializer

import scala.concurrent.{ExecutionContext, Future}

object OutboundMerchantPaymentKafkaSubscriber{

  def props(kafkaTopic: String, kafkaConsumerGroupId: String, kafkaConsumerInstanceGroupId: String, paymentRouter:ActorRef)(implicit  mat: ActorMaterializer): Props = Props(new OutboundMerchantPaymentKafkaSubscriber(kafkaTopic,kafkaConsumerGroupId,kafkaConsumerInstanceGroupId,paymentRouter))
}



/**
 * @author Oluwaseyi Otun
 *         OutboundMerchantPaymentSubscriber is a kafka consumer
 *
 */

class OutboundMerchantPaymentKafkaSubscriber(kafkaTopic: String, kafkaConsumerGroupId: String, kafkaConsumerInstanceGroupId: String, paymentRouter: ActorRef)(implicit val mat: ActorMaterializer)  extends  Actor with ActorLogging{

  implicit val dispacter: ExecutionContext = context.system.dispatcher


  override def preStart(): Unit = {
    log.info("Start consuming Merchant Payment Event From kafka {}", kafkaTopic)
    startMerchantPaymentConsumerEvent()
  }

 def receive: Receive = {
      case _ => "Nothing to Do"
    }



  def startMerchantPaymentConsumerEvent(): Unit = {
    log.info("Start Consuming Merchant Payment Events from Kafka Topic {} GroupId {} ", kafkaTopic, kafkaConsumerGroupId)
    val consumerSettings = ConsumerSettings.create(context.system, new StringDeserializer, new StringDeserializer)
      .withGroupId(kafkaConsumerGroupId)
        .withGroupInstanceId(kafkaConsumerInstanceGroupId)
    Consumer.plainSource(consumerSettings, Subscriptions.topics(kafkaTopic)).mapAsync(5)(sendMerchantPaymentRequestToPSPActor)
      .runWith(Sink.ignore)
  }

   def sendMerchantPaymentRequestToPSPActor(event: ConsumerRecord[String, String]): Future[Done] = {
   val payload = event.value()
    log.info("Handle  Incoming Merchant Payment Payload Event {}  ", payload )
     val merchantPaymentCmd = new Gson().fromJson(payload,classOf[MerchantPaymentRequestCommand])

     paymentRouter ! merchantPaymentCmd
     log.info("Merchant payment Sent to Payment Service Provider Worker Actor Pool {} Payment Reference {}  ", merchantPaymentCmd,merchantPaymentCmd.paymentRequest.merchantPaymentReference )
     return Future.successful(Done)
  }
}
