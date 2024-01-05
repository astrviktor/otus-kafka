package akka_akka_streams.Akka

import akka.actor.typed.scaladsl.Behaviors
import akka.actor.typed.{ActorRef, ActorSystem, Behavior, SpawnProtocol}
import akka.util.Timeout

import scala.concurrent.ExecutionContext
import akka.actor.typed.scaladsl.AskPattern._
import akka.actor.typed.Props
import akka_akka_streams.Akka.intro_actors.behaviour_factory_method

import scala.concurrent.Future
import scala.language.{existentials, postfixOps}
import scala.concurrent.duration._

object AkkaMain {
  def main(args: Array[String]): Unit = {
    val system = ActorSystem[String](behaviour_factory_method.Echo(), "Echo")
    system ! "Hello"

    Thread.sleep(3000)
    system.terminate()
  }

}