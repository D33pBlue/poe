// Copyright 2020 D33pBlue

// Package ga provides functions and types to define
// and execute genetic algorithms.
package ga

import (
  "sort"
  "fmt"
  "math/rand"
)

// generatePopulation returns a random generated
// population of n Sol, using an existing prng
// and a dna.
func generatePopulation(n int,dna DNA,prng *rand.Rand) (pop Population){
  for i:=0; i<n; i++{
    sol := new(Sol)
    sol.Individual = dna.Generate(prng)
    sol.Fitness = 9999999999999
    sol.IsEval = false
    pop = append(pop,*sol)
  }
  return
}

func selectStd(pop Population,n int) (selected Population) {
  if len(pop)<=n{
    return pop
  }
  sort.Sort(pop)
  for i:=0; i<n-1; i++ {
    selected = append(selected,pop[i])
  }
  return
}

func offspring(pop Population,n int,pcross,pmut float64,prng *rand.Rand)(off Population){
  off = append(off,pop[0])
  for i:=0;i<n;i++{
    j := prng.Intn(len(pop))
    var ind Sol = *new(Sol)
    ind.Individual = pop[j].Individual.DeepCopy()
    ind.IsEval = false
    if prng.Float64()<pcross{
      ind2 := pop[prng.Intn(len(pop))]
      ind.Individual = ind.Individual.Crossover(ind2.Individual,prng)
      ind.IsEval = false
    }
    if prng.Float64()<pmut{
      ind.Individual = ind.Individual.Mutate(prng)
      ind.IsEval = false
    }
    off = append(off,ind)
  }
  return
}


func RunGA(dna DNA,conf *Config,chOut,chIn chan Packet){// (Population,Sol) {
  if dna.HasToMinimize(){Optimum = Minimize
  }else{Optimum = Maximize}
  var prng *rand.Rand = rand.New(rand.NewSource(99))
  prng.Seed(conf.Miner)
  var population Population = generatePopulation(conf.NPop,dna,prng)
  var bestOfAll Sol = population.eval(conf.BlockHash)
  // bestOfAll.eval()
  for epoch:=0; epoch<conf.Gen; epoch++{
    population = selectStd(population,conf.Mu)
    population = offspring(population,conf.Lambda,conf.Pcross,conf.Pmut,prng)
    population.reset()
    best := population.eval(conf.BlockHash)
    if Optimum(best.Fitness,bestOfAll.Fitness){
      bestOfAll = best
    }
    if conf.Verbose==2 {
      fmt.Printf("[%d]gen %d best fit: %f\n",conf.Miner,epoch,best.Fitness)
    }
    if epoch%conf.Step==0{
      if conf.Verbose==1{
        fmt.Printf("[%d,%d,%f]",conf.Miner,epoch,bestOfAll.Fitness)
      }
      bestOfAll.Conf = *conf
      bestOfAll.Gen = epoch
      pk := new(Packet)
      pk.Solution = bestOfAll
      pk.End = false
      chOut <- *pk
      pk2 := <- chIn
      for ii:=0;ii<len(pk2.Shared);ii++{
        ss := pk2.Shared[ii]
        ss.Conf = *conf
        population = append(population,ss)
      }
    }
  }
  // sort.Sort(population)
  bestOfAll.Conf = *conf
  bestOfAll.Gen = conf.Gen
  pk := new(Packet)
  pk.Solution = bestOfAll
  pk.End = true
  chOut <- *pk
}
