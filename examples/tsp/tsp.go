/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-19
 * @Project: Proof of Evolution
 * @Filename: tsp.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-17
 * @Copyright: 2020
 */



package main

import(
  "math/rand"
  "io/ioutil"
  "encoding/json"
  "github.com/D33pBlue/poe/op"
  . "github.com/D33pBlue/poe/ga"
)


type Problem struct{
  TourSize int
  DistanceMatrix [][]float64
  OptDistance float64
}

var problem Problem

type TSP struct{
  Data []int
}

type definition string

func (s definition)Initialize(path string){
  dat, err := ioutil.ReadFile(path)
  if err != nil { panic(err) }
  err = json.Unmarshal([]byte(dat), &problem)
  if err != nil { panic(err) }
}

func build()*TSP{
  var tsp *TSP = new(TSP)
  for i:=0;i<problem.TourSize;i++{
    tsp.Data = append(tsp.Data,i)
  }
  return tsp
}

func (s definition)New()DNA{
  return build()
}

func (self *TSP)Generate(prng *rand.Rand) DNA{
  tsp := build()
  prng.Shuffle(len(tsp.Data),func(i,j int){tsp.Data[i],tsp.Data[j]=tsp.Data[j],tsp.Data[i]})
  return tsp
}

func (self *TSP)Mutate(prng *rand.Rand) DNA{
  i := prng.Intn(len(self.Data))
  j := prng.Intn(len(self.Data))
  for ;i==j && len(self.Data)>1;{
    j = prng.Intn(len(self.Data))
  }
  self.Data[i],self.Data[j] = self.Data[j],self.Data[i]
  return self
}

func (self *TSP)Crossover(ind2 DNA,prng *rand.Rand) DNA{
  var second *TSP = ind2.(*TSP)
  j := prng.Intn(len(self.Data))
  var rem []int
  for i:=0;i<len(self.Data);i++{
    n := self.Data[i]
    var used bool = false
    for k:=0;k<j && !used;k++{
      if second.Data[k]==n{
        used = true
      }
    }
    if !used{
      rem = append(rem,n)
    }
  }
  for i:=0;i<j;i++{
    self.Data[i] = second.Data[i]
  }
  for i:=0;i<len(rem);i++{
    self.Data[j+i] = rem[i]
  }
  return self
}

// ---------------------------------
// Evaluate function in clear:
//
// func (self *TSP)Evaluate(st *op.State) float64{
//   var dist float64 = 0.0
//   for i:=1;i<len(self.Data);i++{
//     n1,n2 := self.Data[i-1],self.Data[i]
//     dist += problem.DistanceMatrix[n1][n2]
//   }
//   dist += problem.DistanceMatrix[self.Data[len(self.Data)-1]][self.Data[0]]
//   return dist
// }
// --------------------------------

func (self *TSP)Evaluate(st *op.State) float64{
  var dist float64
  st.SetFloat64(&dist,0.0)
  var i int
  for st.SetInt(&i,1);st.Lt(i,len(self.Data));st.SetInt(&i,st.Succ(i)){
    var n1 int
    var n2 int
    st.SetInt(&n1,self.Data[st.Sub(i,1).(int)])
    st.SetInt(&n2,self.Data[i])
    st.SetFloat64(&dist,st.Add(dist,problem.DistanceMatrix[n1][n2]))
  }
  var last int
  var first int
  st.SetInt(&last,self.Data[st.Sub(len(self.Data),1).(int)])
  st.SetInt(&first,self.Data[0])
  st.SetFloat64(&dist,st.Add(dist,problem.DistanceMatrix[last][first]))
  return dist
}

func (self *TSP)DeepCopy() DNA{
  k := build()
  for el := range self.Data{
    k.Data[el] = self.Data[el]
  }
  return k
}

func (self *TSP)HasToMinimize() bool{
  return true
}

func (self *TSP)Serialize()[]byte{
  data,_ := json.Marshal(self)
  return data
}

func (self *TSP)LoadFromSerialization(data []byte){
  var objmap map[string]json.RawMessage
  json.Unmarshal(data, &objmap)
  json.Unmarshal(objmap["Data"],&self.Data)
}


var Definition definition
