package main

type customError struct {
  msg string
}

func (e *customError) Error() string {
  return e.msg
}

func test() *customError {
  // ... do something
  return nil
}

func main() {
  // заведомо объявили err как интерфейс error
  var err error
  err = test() // возвращает nil

  // interface == nil <=> когда type == nil && value == nil
  // так как мы уже объявили тип err, оно замедомо не nil
  if err != nil {
    println("error")
    return
  }
  println("ok")
}