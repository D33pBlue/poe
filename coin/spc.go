/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-18
 * @Project: Proof of Evolution
 * @Filename: spc.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-22
 * @Copyright: 2020
 */



package coin

import(
  "fmt"
)

// SpC (Simple Coin) is the default implementation
// of the Currency interface.
type SpC struct{
  value float64
}

func (s SpC)New(v Any)Currency{
  c := new(SpC)
  c.value = -1.0
  switch i := v.(type) {
    case float64: c.value = i
    case float32: c.value = float64(i)
    case int: c.value = float64(i)
    case int8: c.value = float64(i)
    case int16: c.value = float64(i)
    case int32: c.value = float64(i)
    case int64: c.value = float64(i)
  }
  return c
}

func (s *SpC)String() string{
  return fmt.Sprintf("%v§",s.value)
}

func (s *SpC)Check()bool{
  if s.value>=0 {
    return true
  }
  s.value = -1
  return false
}

func (s *SpC)Compare(b Currency) int{
  if s.value < b.(*SpC).value{
    return -1
  }
  if s.value > b.(*SpC).value{
    return 1
  }
  return 0
}

func (s *SpC)Add(b Currency)Currency{
  if s.Check()&&b.Check(){
    return s.New(s.value+b.(*SpC).value)
  }
  return s.New("")
}

func (s *SpC)Sub(b Currency)Currency{
  if s.Check()&&b.Check(){
    return s.New(s.value-b.(*SpC).value)
  }
  return s.New("")
}

func (s *SpC)Mul(t float64)Currency{
  if s.Check(){
    return s.New(s.value*t)
  }
  return s.New("")
}

func (s *SpC)Div(t float64)Currency{
  if s.Check()&&t!=0{
    return s.New(s.value/t)
  }
  return s.New("")
}
