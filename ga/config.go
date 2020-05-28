/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-18
 * @Project: Proof of Evolution
 * @Filename: config.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-25
 * @Copyright: 2020
 */



package ga

import(
  "fmt"
  "sync"
  "math/rand"
  "encoding/binary"
)

// Config collect the execution parameters
// the miner chooses to use while executing a GA.
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
  keepmining *bool
  access_hash sync.Mutex
}


func BuildBlockchainGAConfig(hash []byte,keepmining *bool,step int)*Config{
  fmt.Println("hash len:",len(hash))
  x := int64(binary.BigEndian.Uint64(hash[:8]))
  conf := RandConf(x,0,step)
  conf.keepmining = keepmining
  conf.BlockHash = hash
  conf.Verbose = 0
  return conf
}

func (self *Config)GetHash()[]byte{
  self.access_hash.Lock()
  var hash []byte
  for i:=0;i<len(self.BlockHash);i++{
    hash = append(hash,self.BlockHash[i])
  }
  self.access_hash.Unlock()
  return hash
}

func (self *Config)ChangeHash(hash []byte){
  self.access_hash.Lock()
  self.BlockHash = hash
  self.access_hash.Unlock()
}

// Change the current block's hash (useful to
// update the complexity accordingly).
func (c *Config)SetBlockHash(hash []byte){
  c.BlockHash = hash
}

// Generate a Config with default parameters.
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

// Generate a Config with random parameters.
func RandConf(x int64,gen,step int)*Config{
  var prng *rand.Rand = rand.New(rand.NewSource(99))
  prng.Seed(x)
  conf := new(Config)
  conf.Miner = x
  conf.Gen = gen
  conf.Step = step
  conf.NPop = prng.Intn(600)+200
  conf.Pcross = prng.Float64()+0.0001
  conf.Pmut = prng.Float64()+0.0001
  conf.Mu = prng.Intn(350)+100
  conf.Lambda = prng.Intn(350)+100
  conf.Verbose = 1
  return conf
}
