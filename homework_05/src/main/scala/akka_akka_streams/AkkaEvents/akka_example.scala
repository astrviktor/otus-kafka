package akka_akka_streams.AkkaEvents

import akka.NotUsed
import akka.actor.typed.scaladsl.{ActorContext, Behaviors}
import akka.actor.typed.{ActorSystem, Behavior, Props}
import akka.persistence.typed.PersistenceId
import akka.persistence.typed.scaladsl.{Effect, EventSourcedBehavior}


// 100
// -5
// +3
// -4

// 1. initial state
//2. book
//3. refund

object akka_example{
  object BookingSystem {
    def apply(): Behavior[Command] = Behaviors.setup{ctx =>
      val persId = "BookingSystem"
      EventSourcedBehavior[Command, Event, State](
        PersistenceId.ofUniqueId(persId),
        State.empty,
        (state, command) => handleCommand(persId, state, command, ctx),
        (state, event) => handleEvents(state, event, ctx)
      )
    }
  }

  //main command
  sealed trait Command
  case class Book(number: Int) extends Command
  case class Refund(number: Int) extends Command

  //event
  sealed  trait Event
  case class Booked(id: String, number: Int) extends Event
  case class Refunded(id: String, number: Int) extends Event

  final case class State(value: Int){
    def book(number: Int): State = copy(value = value - number)
    def refund(number: Int): State = copy(value = value + number)
  }

  object State{
    val empty = State(100)
  }

  //handle command
  def handleCommand(persistenceId: String, state: State,
                    command: Command,
                    ctx: ActorContext[Command]): Effect[Event, State] = command match {
    case Book(number) =>
      ctx.log.info(s"event book $number, the total ticket count ${state.value}")
      val booked = Booked(persistenceId, number)
      Effect
      .persist(booked)
      .thenRun{x=>
        ctx.log.info(s"thre rest tickets number is ${x.value}")
      }
    case Refund(number) =>
      ctx.log.info(s"event refund $number, the total ticket count ${state.value}")
      Effect
      .persist(Refunded(persistenceId, number))
      .thenRun{x=>
        ctx.log.info(s"thre rest tickets number is ${x.value}")
      }
  }

  // events
  def handleEvents(state: State, event: Event,
                   ctx: ActorContext[Command]): State = event match{
    case Booked(_, number) =>
      ctx.log.info(s"receive event for tickets $number. The total tickets count is ${state.value}")
      state.book(number)
    case Refunded(_, number) =>
      ctx.log.info(s"receive event for tickets $number. The total tickets count is ${state.value}")
      state.refund(number)
  }

  def apply(): Behavior[NotUsed] = Behaviors.setup{ctx =>
    val bookingSystem = ctx.spawn(BookingSystem(), "BookingSystem", Props.empty)

    bookingSystem ! Book(5)
    bookingSystem ! Refund(2)
    bookingSystem ! Book(4)

    Behaviors.same

  }

  def main(args: Array[String]): Unit = {
    ActorSystem(akka_example(), "akka_typed")
  }


}
