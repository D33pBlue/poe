// Copyright 2020 D33pBlue

package ga

import(
  "math/rand"
)

type Config struct{
  Miner int64
  Gen int
  Step int
  NPop int
  Pcross float64
  Pmut float64
  Mu int
  Lambda int
  Verbose int
  BlockHash []byte
}

func (c *Config)SetBlockHash(hash []byte){
  c.BlockHash = hash
}

func DefConf(x int64,gen,step int)*Config{
  conf := new(Config)
  conf.Miner = x
  conf.Gen = gen
  conf.Step = step
  conf.NPop = 600
  conf.Pcross = 0.7
  conf.Pmut = 0.3
  conf.Mu = 200
  conf.Lambda = 300
  conf.Verbose = 1
  return conf
}

func RandConf(x int64,gen,step int)*Config{
  var prng *rand.Rand = rand.New(rand.NewSource(99))
  prng.Seed(x)
  conf := new(Config)
  conf.Miner = x
  conf.Gen = gen
  conf.Step = step
  conf.NPop = prng.Intn(600)+100
  conf.Pcross = prng.Float64()+0.0001
  conf.Pmut = prng.Float64()+0.0001
  conf.Mu = prng.Intn(350)+50
  conf.Lambda = prng.Intn(350)+50
  conf.Verbose = 1
  return conf
}
