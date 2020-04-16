package op

import (
  "math/rand"
)

type State struct{
  operations float64
  coeff map[string]float64
  off map[string]float64
}

func (s *State)SetBlockHash(hash string){
  s.operations = 0
  rand.Seed(43)
  s.off["+"] = rand.Float64()
  s.off["*"] = rand.Float64()
  s.off["-"] = rand.Float64()
  s.off["/"] = rand.Float64()
  s.off["%"] = rand.Float64()
  s.off["++"] = rand.Float64()
  s.off["--"] = rand.Float64()
  s.off["=="] = rand.Float64()
  s.off["<"] = rand.Float64()
  s.off["<="] = rand.Float64()
  s.off[">"] = rand.Float64()
  s.off[">="] = rand.Float64()
  s.off["neg"] = rand.Float64()
  s.off["and"] = rand.Float64()
  s.off["or"] = rand.Float64()
  s.off["pow"] = rand.Float64()
  s.off["sqrt"] = rand.Float64()
  s.off["abs"] = rand.Float64()
  s.off["sign"] = rand.Float64()
  s.off["ceil"] = rand.Float64()
  s.off["floor"] = rand.Float64()
  s.off["round"] = rand.Float64()
  s.off["min"] = rand.Float64()
  s.off["max"] = rand.Float64()
  s.off["sin"] = rand.Float64()
  s.off["cos"] = rand.Float64()
  s.off["asin"] = rand.Float64()
  s.off["acos"] = rand.Float64()
  s.off["tan"] = rand.Float64()
  s.off["atan"] = rand.Float64()
  s.off["sinh"] = rand.Float64()
  s.off["cosh"] = rand.Float64()
  s.off["asinh"] = rand.Float64()
  s.off["acosh"] = rand.Float64()
  s.off["tanh"] = rand.Float64()
  s.off["atanh"] = rand.Float64()
  s.off["log"] = rand.Float64()
  s.off["log2"] = rand.Float64()
  s.off["log10"] = rand.Float64()
  s.off["exp"] = rand.Float64()
  s.off["exp2"] = rand.Float64()
  s.off["assign"] = rand.Float64()
  s.off["append"] = rand.Float64()
  s.off["delete"] = rand.Float64()
}

func MakeState(hash string) *State {
  var s *State = new(State)
  s.coeff = make(map[string]float64)
  s.off = make(map[string]float64)
  s.coeff["+"] = 1
  s.coeff["*"] = 1
  s.coeff["-"] = 1
  s.coeff["/"] = 1
  s.coeff["%"] = 1
  s.coeff["++"] = 1
  s.coeff["--"] = 1
  s.coeff["=="] = 1
  s.coeff["<"] = 1
  s.coeff["<="] = 1
  s.coeff[">"] = 1
  s.coeff[">="] = 1
  s.coeff["neg"] = 1
  s.coeff["and"] = 1
  s.coeff["or"] = 1
  s.coeff["pow"] = 1
  s.coeff["sqrt"] = 1
  s.coeff["abs"] = 1
  s.coeff["sign"] = 1
  s.coeff["ceil"] = 1
  s.coeff["floor"] = 1
  s.coeff["round"] = 1
  s.coeff["min"] = 1
  s.coeff["max"] = 1
  s.coeff["sin"] = 1
  s.coeff["cos"] = 1
  s.coeff["asin"] = 1
  s.coeff["acos"] = 1
  s.coeff["tan"] = 1
  s.coeff["atan"] = 1
  s.coeff["sinh"] = 1
  s.coeff["cosh"] = 1
  s.coeff["asinh"] = 1
  s.coeff["acosh"] = 1
  s.coeff["tanh"] = 1
  s.coeff["atanh"] = 1
  s.coeff["log"] = 1
  s.coeff["log2"] = 1
  s.coeff["log10"] = 1
  s.coeff["exp"] = 1
  s.coeff["exp2"] = 1
  s.coeff["assign"] = 1
  s.coeff["append"] = 1
  s.coeff["delete"] = 1
  s.SetBlockHash(hash)
  return s
}

func (self *State) NumOperations()float64 {
  return self.operations
}

func (self *State) IncOperations(n float64) {
  if n>0{
      self.operations += n
  }
}

type any interface{}
