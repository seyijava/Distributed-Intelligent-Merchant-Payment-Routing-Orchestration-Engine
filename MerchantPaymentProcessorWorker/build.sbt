name := "payment-processor-worker"
maintainer := "seyijava@gmail.com"
scalaVersion := "2.13.1"

// These options will be used for *all* versions.
scalacOptions ++= Seq(
  "-deprecation",
  "-unchecked",
  "-encoding", "UTF-8",
  "-Xlint",
)


val akka = "2.6.0"
val akkaHttpVersion ="10.1.12"


libraryDependencies ++= Seq (
  // -- Logging --
  "ch.qos.logback" % "logback-classic" % "1.2.3",
  // -- Akka --
  "com.typesafe.akka" %% "akka-actor"   % akka,
  "com.typesafe.akka" %% "akka-cluster" % akka,
  "com.typesafe.akka" %% "akka-stream-kafka" % "2.0.3",
  "com.typesafe.akka" %% "akka-stream" % akka,
  "com.google.code.gson" % "gson" % "2.8.5",
  "com.typesafe.akka" %% "akka-http" %akkaHttpVersion,
  "com.typesafe.akka" %% "akka-http-spray-json" % akkaHttpVersion,
  "org.json" % "json" % "20190722",
  "commons-codec" % "commons-codec" % "1.11",
  "com.braintreepayments.gateway" % "braintree-java" % "3.0.0",
  "com.stripe" % "stripe-java" % "19.23.0",
  "com.adyen" % "adyen-java-api-library" % "7.0.0"

)

version in Docker := "latest"
dockerExposedPorts in Docker := Seq(1600)
dockerRepository := Some("intelligentpaymentrouting")
dockerBaseImage := "java"
enablePlugins(JavaAppPackaging)