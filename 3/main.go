package main

import (
  "fmt"
  "os"
)

// структура интерфейса содержит указатель на значение и указатель на структуру, описывающую общие сведения интерфейса
// самое главное в этой структуре - тип

// interface == nil <=> когда type == nil && value == nil. 


func Foo() error {
  var err *os.PathError = nil // type = os.PathError, value = nil
  return err // ошибка - это интерфейс
}

func main() {
  err := Foo()
  fmt.Println(err) // печатает значение - nil
  fmt.Println(err == nil) // false, т.к. type != nil
}