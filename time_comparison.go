/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-May-19
 * @Project: Proof of Evolution
 * @Filename: time_comparison.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-21
 * @Copyright: 2020
 */

package main

import (
  "fmt"
  "time"
  "math/rand"
  "encoding/binary"
  "github.com/D33pBlue/poe/utils"
  "github.com/D33pBlue/poe/blockchain"
  "github.com/D33pBlue/poe/ga"
  "github.com/D33pBlue/poe/op"
)


func main(){
  fmt.Println("Compare execution times of hash and fitness evaluation")
  jobPath := "/home/d33pblue/Documenti/uni/tesi/poe/examples/tsp/tsp.go"
  dataPath := "/home/d33pblue/Documenti/uni/tesi/poe/examples/tsp/data/tsp0/rat99.json"
  dna,err := ga.LoadGA(jobPath,dataPath)
  if err!=nil{
    fmt.Println(err)
    return
  }
  hb := new(utils.HashBuilder)
  hb.Add(jobPath)
  hash := hb.GetHash()
  x := int64(binary.BigEndian.Uint64(hash[:8]))
  var prng *rand.Rand = rand.New(rand.NewSource(99))
  prng.Seed(x)
  t0 := time.Now().UnixNano()
  st := op.MakeState(hash)
  sol := new(ga.Sol)
  sol.Individual = dna.Generate(prng)
  sol.Fitness = 9999999999999
  sol.IsEval = false
  sol.Eval2(st,hash)
  tGA := time.Now().UnixNano()-t0
  block := blockchain.BuildFirstBlock(utils.Addr("sfvfvw"))
  t0 = time.Now().UnixNano()
  h := block.GetHash("")
  tHash := time.Now().UnixNano()-t0
  fmt.Println(h)
  fmt.Printf("Time for GA:   %v\n",tGA)
  fmt.Printf("Time for Hash: %v\n",tHash)
  fmt.Printf("Complexity for GA: %v (%v)\n",sol.Complex,float64(tGA)/sol.Complex)
  fmt.Printf("Approximation: %v\n",float64(tHash)*sol.Complex/float64(tGA))
}
