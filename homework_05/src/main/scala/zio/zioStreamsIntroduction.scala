package zio

import ch.qos.logback.classic.{Level, Logger}
import zio._
import zio.kafka.consumer._
import zio.kafka.producer.{Producer, ProducerSettings}
import zio.kafka.serde._
import zio.stream.{ZSink, ZStream}
import org.slf4j.LoggerFactory
import zio.Console.printLine

import scala.concurrent.TimeoutException

object zioStreamsIntroduction extends App{
  val emptyStream = ZStream.empty
  val oneIntValueStream = ZStream.succeed(4)
  val oneListValueStream = ZStream.succeed(List(1,2,3))
  val finiteIntStream = ZStream.range(1, 10)
  val infiniteIntStream = ZStream.iterate(1)(_+1)

  val stream = ZStream(1,2,3)
  val unit = ZStream.unit
  val never = ZStream.never

  val s1 = ZStream.fromChunk(Chunk(1,2,3))
  val s2 = ZStream.fromChunks(Chunk(1,2,3),Chunk(1,2,3))

  val s3 = ZStream(1,2,3) ++ ZStream.fail("error") ++ ZStream(4,5)
  val s4 = ZStream(1,2,3)

  val stream1 = s3.orElse(s4)

  val s5 = ZStream(1,2,3) ++ ZStream.fail("Error") ++ ZStream(4,5)
  val s6 = ZStream(1,2,3)

  val stream3 = s5.catchSome{
    case "Error" => s6
    case "Error1" => s4
  }

  val number = ZStream(1,2,3) ++
    ZStream
      .fromZIO(
        Console.print("Enter a number: ") *> Console.readLine
          .flatMap(x=>

          x.toIntOption match {
            case Some(value) => ZIO.succeed(value)
            case None => ZIO.fail("NaN")
          }
          )

      ).retry(Schedule.exponential(1.second))


  stream.timeoutFail(new TimeoutException)(10.seconds)

  // sink
  val sum = ZStream(1,2,3).run(ZSink.sum)

  val s7 = ZStream(1,2,3,4).runFold(0)(_+_)

  //foreach
  ZStream(1,2,3).foreach(printLine(_))


  //operations
  val stream5 = ZStream.iterate(0)(_+1)
  //0,1,2,3,

  val s8 = stream5.take(5)

  val s9 = stream5.takeWhile( _ < 5)
  val s10 = stream5.takeUntil(_ < 5)

  val s11 = s3.takeRight(3)



  val intStream = ZStream.fromIterable(0 to 100)
  val stringStream = intStream.map(_.toString)

  val s17 = ZStream.range(1, 11).filter(_%2==0)

  val s18 = for {
    i <- ZStream.range(1, 11).take(10)
    if i%2 == 0

  } yield i
  // 2 4 5 6 10

  val s19 = ZStream(1,2,3,4).filterNot(_%2==0)










}