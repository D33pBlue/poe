/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-16
 * @Project: Proof of Evolution
 * @Filename: knapsack.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-19
 * @Copyright: 2020
 */



package main

import (
  "math/rand"
  "io/ioutil"
  "encoding/json"
  "github.com/D33pBlue/poe/op"
  . "github.com/D33pBlue/poe/ga"
)


type Element struct{
  Value float64
  Weight float64
}

type Probl struct{
  Items []Element
  Limit float64
}

var problem Probl

type Knap struct{
  data map[int]bool
}

type definition string

func (s definition)Initialize(path string){
  dat, err := ioutil.ReadFile(path)
  if err != nil { panic(err) }
  err = json.Unmarshal([]byte(dat), &problem)
  if err != nil { panic(err) }
}

func build()*Knap{
  var k *Knap = new(Knap)
  k.data = make(map[int]bool)
  for i:=0; i<len(problem.Items); i++{
    k.data[i] = false
  }
  return k
}

func (s definition)New()DNA{
  return build()
}

func (self* Knap)HasToMinimize() bool {
  return false
}

func (self *Knap)Generate(prng *rand.Rand) DNA{
  k := build()
  nelem := 1+prng.Intn(2)
  for i:=0; i<nelem; i++{
    k.data[prng.Intn(len(problem.Items))] = true
  }
  return k
}

func (self* Knap)DeepCopy()DNA{
  k := build()
  for el := range self.data{
    k.data[el] = self.data[el]
  }
  return k
}

func (self *Knap)Mutate(prng *rand.Rand) DNA{
  j := prng.Intn(len(problem.Items))
  self.data[j] = !self.data[j]
  return self
}

func (self *Knap)Crossover(ind2 DNA,prng *rand.Rand) DNA{
  // intersection of the sets
  var second *Knap = ind2.(*Knap)
  for i:=0;i<len(problem.Items)/10;i++{
    j := prng.Intn(len(problem.Items))
    self.data[j] = second.data[j]
  }
  return self
}


// ----------------------------------
// Evaluate in clear:
//
// func (self *Knap)Evaluate(st *op.State) float64{
//   var val float64 = 0.0
//   var wei float64 = 0.0
//   for i:=0; i<len(problem.Items);i++{
//     if self.data[i]{
//       val += problem.Items[i].Value
//       wei += problem.Items[i].Weight
//     }
//   }
//   if wei<=problem.Limit{
//     return val
//   }
//   return -1.0
// }
// ---------------------------------

func (self *Knap)Evaluate(st *op.State) float64{
  var val float64
  st.SetFloat64(&val,0.0)
  var wei float64
  st.SetFloat64(&wei,0.0)
  var i int
  for st.SetInt(&i,0); st.Lt(i,len(problem.Items));st.SetInt(&i,st.Succ(i)){
    if self.data[i]{
      st.SetFloat64(&val,st.Add(val,problem.Items[i].Value))
      st.SetFloat64(&wei,st.Add(wei,problem.Items[i].Weight))
    }
  }
  if st.Le(wei,problem.Limit){
    return val
  }
  return -1.0
}

var Definition definition
