// Copyright 2020 D33pBlue

package coin

import(
  "fmt"
  "testing"
)

func TestCurrency(t *testing.T){
  c1 := New(1234.06)
  c2 := New("c")
  c3 := New(550)
  fmt.Println(c1,c1.Check(),c2,c2.Check())
  fmt.Println(c1.Compare(c2))
  fmt.Println(c2.Compare(c1))
  fmt.Println(c1.Compare(c1))
  fmt.Println(c1.Add(c2))
  fmt.Println(c1.Add(c3))
}
