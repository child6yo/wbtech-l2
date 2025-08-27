package main

import "fmt"

// ВАЖНО - defer в go выполняется после того, как "подготовится" return

// x проинициализирована как возвращаемая переменная функции 
// таким образом defer может изменять x после вычисления возврата
func test() (x int) {
  defer func() {
    x++
  }()
  x = 1
  return
}

// x инициализируется в теле функции, return подготавливает на выход x = 1,
// таким образом defer уже никак не повлияет на возврат
func anotherTest() int {
  var x int
  defer func() {
    x++
  }()
  x = 1
  return x
}

func main() {
  fmt.Println(test())
  fmt.Println(anotherTest())
}