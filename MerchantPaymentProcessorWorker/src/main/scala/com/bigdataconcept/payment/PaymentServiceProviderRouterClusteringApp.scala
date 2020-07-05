package com.bigdataconcept.payment

import akka.actor.{ActorRef, ActorSystem, Props}
import akka.routing.RandomPool
import akka.stream.ActorMaterializer
import com.bigdataconcept.payment.domain.Configs.PaymentServiceProviderConfig
import com.bigdataconcept.payment.processor.OutboundMerchantPaymentKafkaSubscriber
import com.bigdataconcept.payment.service.provider.{PaymentServiceProvider, PaymentServiceProviderFactory}
import com.typesafe.config.ConfigFactory

object PaymentServiceProviderRouterClusteringApp extends App {

  val config = ConfigFactory.load()

  private val clusterName = config.getString("paymentServiceProvider.clusterName")

  implicit  val system = ActorSystem(clusterName)

  implicit val mat = ActorMaterializer()



  private val noOfInstance = system.settings.config getInt "paymentServiceProvider.numberOfInstance"

  private val paymentServiceProviderUrl = system.settings.config getString "paymentServiceProvider.config.serviceUrl"


  private val paymentServiceProviderName = system.settings.config getString "paymentServiceProvider.name"

  private val userName = system.settings.config getString "paymentServiceProvider.config.userName"

  private val password = system.settings.config getString "paymentServiceProvider.config.password"

  private val merchantAccount = system.settings.config getString "paymentServiceProvider.config.merchantAccount"

  private val privateKey = system.settings.config getString "paymentServiceProvider.config.privateKey"

  private  val  merchantId = system.settings.config getString "paymentServiceProvider.config.merchantId"

  private val publicKey = system.settings.config getString "paymentServiceProvider.config.publicKey"



  private val kafkaTopic = system.settings.config getString "kafka.topic"

  private val kafkaGroupId = system.settings.config getString "kafka.groupId"

  private val kafkaGroupInstanceId = system.settings.config getString "kafka.groupInstanceId"

  private val paymentServiceProviderConfig =  PaymentServiceProviderConfig(serviceUrl = paymentServiceProviderUrl, userName=userName, password=
    password,merchantAccount = merchantAccount,merchantId=merchantId,privateKey=privateKey,publicKey=publicKey)

  private val paymentServiceProvider = PaymentServiceProvider.withName(paymentServiceProviderName)

  private val props: Props =   PaymentServiceProviderFactory.apply(paymentServiceProvider,paymentServiceProviderConfig)

  private val paymentServiceProviderRouter: ActorRef = system.actorOf(props.withRouter(RandomPool(noOfInstance)), "paymentServiceProviderRouter")

   system.actorOf(OutboundMerchantPaymentKafkaSubscriber.props(kafkaTopic,kafkaGroupId,kafkaGroupInstanceId,paymentServiceProviderRouter))
}


/*
val router = system.actorOf(ClusterRouterPool(
  RoundRobinPool(0),
  ClusterRouterPoolSettings(
    totalInstances = 20,
    maxInstancesPerNode = 1,
    allowLocalRoutees = false,
    useRole = None
  )
).props(Props[Worker]), name = "router")
 */