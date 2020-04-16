package main

import (
  "fmt"
  "github.com/D33pBlue/poe/ga"
)

func executeGA(dna ga.DNA,chIn,chOut chan ga.Packet){
  var seed int64 = 42
  var generations int = 100
  var interrupt int = 10
  ga.RunGA(dna,ga.RandConf(seed,generations,interrupt),chOut,chIn)
}

func main(){
  fmt.Println("Proof of Evolution")
  dna := ga.LoadGA("examples/tsp/","tsp",
    "examples/tsp/data/tsp0/burma14.json")
  chOut := make(chan ga.Packet)
  chIn := make(chan ga.Packet)
  go executeGA(dna,chIn,chOut)
  var best ga.Population
  for stop:=false;!stop;{
    pk := <-chOut
    best = append(best,pk.Solution)
    if pk.End {
      fmt.Printf("Final fit: %f\n",pk.Solution.Fitness)
      stop = true
    }else{
      toSend := new(ga.Packet)
      chIn <- *toSend
    }
  }
  best_of_all := best[0]
  for j:=0;j<len(best);j++{
    if ga.Optimum(best[j].Fitness,best_of_all.Fitness){
      best_of_all = best[j]
    }
  }
  fmt.Println("\n config: ")
  fmt.Println(best_of_all.Conf)
  fmt.Printf("Best fitness: %f (complex: %f)\n",
    best_of_all.Fitness,best_of_all.Complex)
}
