package zio

import ch.qos.logback.classic.{Level, Logger}
import zio._
import zio.kafka.consumer._
import zio.kafka.producer.{Producer, ProducerSettings}
import zio.kafka.serde._
import zio.stream.ZStream
import org.slf4j.LoggerFactory

object MainApp extends ZIOAppDefault {
  LoggerFactory
    .getLogger(org.slf4j.Logger.ROOT_LOGGER_NAME)
    .asInstanceOf[Logger]
    .setLevel(Level.ERROR)

  val producer =
    ZStream
      .repeatZIO(Random.nextIntBetween(0, Int.MaxValue))
      .schedule(Schedule.fixed(2.seconds))
      .mapZIO{ random =>
        Producer.produce[Any, Long, String](
          topic = "random",
          key = random % 4,
          value = random.toString,
          keySerializer = Serde.long,
          valueSerializer = Serde.string
        )
      }
      .drain

  val consumer =
    Consumer
      .plainStream(Subscription.topics("random"), Serde.long, Serde.string)
      .tap(r=> Console.printLine(r.value))
      .map(_.offset)
      .aggregateAsync(Consumer.offsetBatches)
      .mapZIO(_.commit)
      .drain


  def producerLayer =
    ZLayer.scoped(
      Producer.make(
        settings = ProducerSettings(List("localhost:29092"))
      )
    )

  def consumerLayer =
    ZLayer.scoped(
      Consumer.make(
        ConsumerSettings(List("localhost:29092")).withGroupId("group")
      )

    )

  override def run =
    producer.merge(consumer)
      .runDrain
      .provide(producerLayer, consumerLayer)
}
