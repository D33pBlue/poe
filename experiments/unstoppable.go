package main

import(
  "fmt"
  "time"
)

func content(name string){
  fmt.Println("Keep going",name)
  time.Sleep(1*time.Second)
}

func unstoppable(name string)int{
  for{
    content(name)
  }
  return 0
}

func limit(name string){
  ch := make(chan int, 1)
  go func() {
    ch <- unstoppable(name)
  }()
  select{
    case v := <-ch: fmt.Println("unstoppable returned",v)
    case <-time.After(5*time.Second): fmt.Println("time limit")
  }
  fmt.Println("limit ends")
}


func main(){
  go limit("A")
  unstoppable("B")
}
