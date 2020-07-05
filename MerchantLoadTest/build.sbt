import com.typesafe.sbt.SbtMultiJvm.multiJvmSettings
import com.typesafe.sbt.SbtMultiJvm.MultiJvmKeys.MultiJvm



val akka = "2.6.0"
val akkaHttpVersion ="10.1.12"

scalaVersion := "2.12.8"



assemblyMergeStrategy in assembly := {
  case "application.conf"     => MergeStrategy.concat
  case PathList(ps @ _*) if ps.last contains "module-info" => MergeStrategy.first
  case x =>
    val oldStrategy = (assemblyMergeStrategy in assembly).value
    oldStrategy(x)
}

lazy val `payment-load-test` = project
  .in(file("."))
  .settings(multiJvmSettings: _*)
  .settings(
    organization := "payment-load-test",
    scalaVersion := "2.12.8",
    scalacOptions in Compile ++= Seq("-deprecation", "-feature", "-unchecked", "-Xlog-reflective-calls", "-Xlint"),
    javacOptions in Compile ++= Seq("-Xlint:unchecked", "-Xlint:deprecation"),
    javaOptions in run ++= Seq("-Xms128m", "-Xmx1024m", "-Djava.library.path=./target/native"),
    libraryDependencies ++= Seq(
      "com.typesafe.akka" %% "akka-actor"   % akka,
      "com.typesafe.akka" %% "akka-stream" % akka,
      "com.google.code.gson" % "gson" % "2.8.5" ,
      "com.typesafe.akka" %% "akka-http" % akka % akkaHttpVersion,
      "com.typesafe.akka" %% "akka-http-spray-json" % akkaHttpVersion,
      "org.json" % "json" % "20190722",
      "commons-codec" % "commons-codec" % "1.11",
      "com.github.stevenchen3" %% "scala-faker" %  "0.1.1",
      "com.github.javafaker" % "javafaker" % "1.0.1",
    ),
    fork in run := true,
    mainClass in (Compile, run) := Some("com.bigdataconcept.payment.load.test.LoadRunner"),
    // disable parallel tests
    parallelExecution in Test := false,
    licenses := Seq(("CC0", url("http://creativecommons.org/publicdomain/zero/1.0")))
  )
  .configs (MultiJvm)





