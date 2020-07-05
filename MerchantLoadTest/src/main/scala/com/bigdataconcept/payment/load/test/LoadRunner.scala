package com.bigdataconcept.payment.load.test

import java.util.concurrent.TimeUnit
import akka.actor.ActorSystem
import com.typesafe.config.ConfigFactory
import scala.concurrent.duration._
import scala.concurrent.ExecutionContext
import collection.JavaConverters._

object LoadRunner extends App {

  val config = ConfigFactory.load()

  val testDuration = config.getDuration("load-test.duration", TimeUnit.MILLISECONDS).millis
  val parallelism = config.getInt("load-test.parallelism")
  val rampUpTime = config.getDuration("load-test.ramp-up-time", TimeUnit.MILLISECONDS).millis
  val ip = config.getString("gateway.ip")

  implicit val system: ActorSystem = ActorSystem()
  implicit val executionContext: ExecutionContext = system.dispatcher

  val ports = config.getIntList("gateway.ports").asScala.map(_.intValue()).toList

  system.log.info(s"Creating $parallelism simulations")
  (1 to parallelism).map { _ =>
    Thread.sleep((rampUpTime/parallelism).toMillis)
    val sim = system.actorOf(Simulator.props(ip,ports))
      system.scheduler.scheduleOnce(testDuration, sim, Simulator.Stop)
  }

  system.scheduler.scheduleOnce(testDuration + 15.seconds) {
    system.terminate()
  }

}
