package akka_akka_streams.homework

import akka.actor.ActorSystem
import akka.stream.{ActorMaterializer, ClosedShape, UniformFanInShape}
import akka.stream.scaladsl.RunnableGraph
import akka.NotUsed
import akka.actor.ActorSystem
import akka.stream.scaladsl.{Broadcast, Flow, GraphDSL, RunnableGraph, Sink, Source, Zip, ZipWith}


object homeworktemplate {
  implicit val system = ActorSystem("fusion")
  implicit val materializer = ActorMaterializer()

  val sumOfThree = GraphDSL.create() { implicit b =>
    import GraphDSL.Implicits._

    val zip1 = b.add(ZipWith[Int, Int, Int](math.addExact _))
    val zip2 = b.add(ZipWith[Int, Int, Int](math.addExact _))
    zip1.out ~> zip2.in0

    UniformFanInShape(zip2.out, zip1.in0, zip1.in1, zip2.in1)
  }

  val graph = GraphDSL.create(){ implicit  builder: GraphDSL.Builder[NotUsed] =>
    import GraphDSL.Implicits._

    //1. source
    val input = builder.add(Source(1 to 5))

    val broadcast = builder.add(Broadcast[Int](3))

    val multiplier10 = builder.add(Flow[Int].map(x=>x*10))
    val multiplier2 = builder.add(Flow[Int].map(x=>x*2))
    val multiplier3 = builder.add(Flow[Int].map(x=>x*3))

    val output = builder.add(Sink.foreach[(Int)](println))

    val sum = builder.add(sumOfThree)

    //shape
    input ~> broadcast

    broadcast.out(0) ~> multiplier10 ~> sum.in(0)
    broadcast.out(1) ~> multiplier2 ~> sum.in(1)
    broadcast.out(2) ~> multiplier3 ~> sum.in(2)

    sum.out ~> output

    //close shape
    ClosedShape
  }

  def main(args: Array[String]) : Unit ={
    RunnableGraph.fromGraph(graph).run()

  }
}