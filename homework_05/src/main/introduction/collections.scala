object DataCollection {
  def main(args: Array[String]): Unit = {
    val l: List[String] = List("1", "2", "3")
    val collection1 = "line 1" :: "line 2" :: "line 3" :: "line 3" :: Nil
    val collection2 = collection1.toSet

    collection2.foreach(x => println(x))


    val collection3 = collection1.groupBy(x=>x).map(x=>x._1)
    collection3.foreach(x => println(x))

    val iter = collection1.iterator
    while (iter.hasNext)
      println(iter.next)

    // fold fold left|right reduce
    val demoCollection = 1::2::3::Nil
    //fold
    println(s"fold result : ${demoCollection.fold(0)((z,i) => z+i)}")
    // fold left"right
    println(s"fold left: ${demoCollection.foldLeft(1)((z,i)=> z+i)}")
    // reduce
    println(s"reduce: ${demoCollection.reduce((z,i)=>z+i)}")

    val test = List(1,2,3,4,5) :: List(1,50,3) :: List(1,2) :: Nil

    test.filter(x=>x.reduce((y,z)=>y+z) > 10).foreach(x=>println(x.mkString(",")))



    val RGB = Seq("R", "G", "B")
    val range = Range(1,4)
    val map = Map("R"-> "Red", "G"->"Green", "B"->"Blue")
/*
    for (el <- RGB)
      println(el)

    for (el <- RGB; el1 <- range)
      println(s"$el - $el1")

    for ((key,value) <- map)
      println(s"$key - $value")
*/
    for (
      el1 <- RGB;
      el2 <- RGB;
      el3 <- RGB;
      if el1 != el2;
      if el3 != el2 && el3 != el1
    ){
      println(s"$el1 $el2 $el3" )
    }

    case class Car(marke: String, model: String, year: Int)
    val cars = Car("VW", "Passat", 2005) :: Car("Lexus", "Ux", 2019) ::Car("BMW", "i3", 2021) :: Nil
    case class Garage(name: String, index: Int)
    val garages = Garage("BMW", 1) :: Garage("Lexus", 2):: Nil

    garages.flatMap{
      garage =>
        cars.filter(car => car.marke == garage.name).map(car => (car.marke, car.model, garage.index))
    }.foreach(x=> println(s"${x._1} ${x._2}  ${x._3}"))

    println("for comp.")
    val car1 = for {
      car <- cars
      garage <- garages
      if car.marke == garage.name
    }
    yield {
      (car.marke, car.model, garage.index)
    }

    car1.foreach{
      case (marke,model,index) => println(s"$marke $model $index" )

    }









  }


}