package akka_akka_streams.AkkaDataStreams

import akka.NotUsed
import akka.actor.typed.{ActorRef, ActorSystem, Behavior, DeathPactException, SupervisorStrategy}
import akka.actor.typed.scaladsl.Behaviors

object Protocol{
  sealed trait Command
  case class Fail(text: String) extends Command
  case class Hello(text: String) extends Command
}

import Protocol._

object Worker {
  def apply(): Behavior[Command] =
    Behaviors.receiveMessage{
      case Fail(text) =>
        throw new RuntimeException(text)
      case Hello(text) =>
        println(text)
        Behaviors.same
    }
}

object MiddleManagment {
  def apply(): Behavior[Command] =
    Behaviors.setup[Command]{ context =>
      context.log.info("MiddleManagment starting up")
      val child = context.spawn(Worker(),"child")
      context.watch(child)
      Behaviors.receiveMessage{message =>
        child ! message
        Behaviors.same
      }

    }
}

object Boss {
  def apply(): Behavior[Command] =
    Behaviors.supervise(
      Behaviors.setup[Command]{ context =>
        context.log.info("Boss starting up")
        val middleManagment = context.spawn(MiddleManagment(), "middle-management")
        context.watch(middleManagment)
        Behaviors.receiveMessage[Command]{message =>
          middleManagment ! message
          Behaviors.same
        }

      }

    ).onFailure[DeathPactException](SupervisorStrategy.restart)
}

object  StartStop1 extends App{
  def apply(): Behavior[NotUsed] =
    Behaviors.setup{ ctx =>
      val boss = ctx.spawn(Boss(), "upper-management")
      boss ! Hello("hi 1")
      boss.tell(Fail("ping"))
      Thread.sleep(1000)
      boss ! Hello("hi 2")

      Behaviors.same

    }
  val value = StartStop1()
  implicit val system = ActorSystem(value, "skjhdgf")
  Thread.sleep(5000)
  system.terminate()
}