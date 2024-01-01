package akka_akka_streams.Akka

import akka.actor.typed.Behavior
import akka.actor.typed.scaladsl.Behaviors
import akka.actor.typed.scaladsl.ActorContext
import akka.actor.typed.scaladsl.AbstractBehavior


object intro_actors {
  // 1.functional
  object behaviour_factory_method {
    object  Echo {
      def apply(): Behavior[String] = Behaviors.setup{ ctx =>
        Behaviors.receiveMessage{
          case msg =>
            ctx.log.info(msg)
            Behaviors.same
        }
      }
    }
  }

  //2. OOP
  object abstract_behaviour {
    class Echo(ctx: ActorContext[String]) extends  AbstractBehavior[String](ctx) {
      override def onMessage(msg: String): Behavior[String] = {
        ctx.log.info(msg)
        this
      }
    }

    object Echo{
      def apply(): Behavior[String] = Behaviors.setup{ctx =>
        new Echo(ctx)
      }
    }
  }
}