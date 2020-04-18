// Copyright 2020 D33pBlue

// Package op defines the operators the user should
// use to define his Evaluator. When used, these operators
// update a State. At the end of the execution
// the miner can get an approximation of Evaluator's complexity.
package op

import (
  "math/rand"
  "encoding/binary"
)

type State struct{
  operations float64
  coeff map[string]float64
  off map[string]float64
}

func (s *State)Reset(){
  s.operations = 0
}

func (s *State)SetBlockHash(hash []byte){
  var prng *rand.Rand = rand.New(rand.NewSource(99))
  prng.Seed(int64(binary.BigEndian.Uint64(hash)))
  s.off["+"] = prng.Float64()
  s.off["*"] = prng.Float64()
  s.off["-"] = prng.Float64()
  s.off["/"] = prng.Float64()
  s.off["%"] = prng.Float64()
  s.off["++"] = prng.Float64()
  s.off["--"] = prng.Float64()
  s.off["=="] = prng.Float64()
  s.off["<"] = prng.Float64()
  s.off["<="] = prng.Float64()
  s.off[">"] = prng.Float64()
  s.off[">="] = prng.Float64()
  s.off["neg"] = prng.Float64()
  s.off["and"] = prng.Float64()
  s.off["or"] = prng.Float64()
  s.off["pow"] = prng.Float64()
  s.off["sqrt"] = prng.Float64()
  s.off["abs"] = prng.Float64()
  s.off["sign"] = prng.Float64()
  s.off["ceil"] = prng.Float64()
  s.off["floor"] = prng.Float64()
  s.off["round"] = prng.Float64()
  s.off["min"] = prng.Float64()
  s.off["max"] = prng.Float64()
  s.off["sin"] = prng.Float64()
  s.off["cos"] = prng.Float64()
  s.off["asin"] = prng.Float64()
  s.off["acos"] = prng.Float64()
  s.off["tan"] = prng.Float64()
  s.off["atan"] = prng.Float64()
  s.off["sinh"] = prng.Float64()
  s.off["cosh"] = prng.Float64()
  s.off["asinh"] = prng.Float64()
  s.off["acosh"] = prng.Float64()
  s.off["tanh"] = prng.Float64()
  s.off["atanh"] = prng.Float64()
  s.off["log"] = prng.Float64()
  s.off["log2"] = prng.Float64()
  s.off["log10"] = prng.Float64()
  s.off["exp"] = prng.Float64()
  s.off["exp2"] = prng.Float64()
  s.off["assign"] = prng.Float64()
  s.off["append"] = prng.Float64()
  s.off["delete"] = prng.Float64()
}

func MakeState(hash []byte) *State {
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
  s.Reset()
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
