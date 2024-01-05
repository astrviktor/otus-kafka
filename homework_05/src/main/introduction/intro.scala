object scalaintro extends App{
  //1. everything is object
  val a = 1+1
  val b: Unit = println("Hi")

  val inc: Int => Int =x => x+1

  lazy val cond: Boolean = true

  val x1 = if (cond) "yes" else "no"

  val x2: Any = if (cond) println("yes") else "no"

  val x: Any = do {
    println("sdfsedf")
  } while (cond)

  val l = List(1,2,3)
  for (i<-0 until 10 by 2){
    println(i)
  }





}



object scalaintro1 {

  def main(args: Array[String]): Unit = {
    println("test1")
  }

}