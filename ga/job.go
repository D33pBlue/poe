/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-May-14
 * @Project: Proof of Evolution
 * @Filename: job.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-25
 * @Copyright: 2020
 */
package ga

import(
  "fmt"
  "math/rand"
  "encoding/binary"
  "github.com/D33pBlue/poe/utils"
  "github.com/D33pBlue/poe/op"
)

type Job struct{
  ChNonce chan Sol // used to send nonce candidates
  ChUpdateIn chan Sol // used to receive shared solutions from miners
  ChUpdateOut chan Sol // used to share good solutions to miners
  KeepRunning bool
  dna DNA
  conf *Config
  jobHash string
}

// Initialize a Job, loading its dna and data from file after
// compiling a user-defined plugin.
func BuildJob(jobpath,datapath string,chUpdateOut chan Sol,jobHash string)*Job{
  // initialize the channels with buffers in order to made them async
  job := new(Job)
  job.ChNonce = make(chan Sol,1000)
  job.ChUpdateIn = make(chan Sol,100)
  job.ChUpdateOut = chUpdateOut
  // compile and load DNA
  var err error
  job.dna,err = LoadGA(jobpath,datapath)
  if err!=nil{
    fmt.Println(err)
    return nil
  }
  job.jobHash = jobHash
  return job
}

// Keeps the GA execution alive and should be called in a goroutine.
// Stops when self.KeepRunning==false.
// Add to the current population the Sol passed to self.ChUpdateIn
// and publish the best found Sol to self.ChUpdateOut every 100 epochs.
func (self *Job)Execute(hashPrev,publicKey string){
  self.KeepRunning = true
  hb := new(utils.HashBuilder)
  hb.Add(hashPrev)
  hb.Add(publicKey)
  hash := hb.GetHash()
  self.conf = BuildBlockchainGAConfig(hash,&self.KeepRunning,500)
  RunGA(self.dna,self.conf,self.ChUpdateOut,self.ChUpdateIn,self.ChNonce,self.jobHash)
}

// Change the hash used during the job execution to initialize
// the coefficients of the operations, that are stored in the state.
func (self *Job)ChangeBlockHash(hashPrev,publicKey string){
  hb := new(utils.HashBuilder)
  hb.Add(hashPrev)
  hb.Add(publicKey)
  hash := hb.GetHash()
  self.conf.ChangeHash(hash)
  // clean the channel ChNonce from Sol with old hash
  sol := <- self.ChNonce
  for ; (!utils.CompareSlices(sol.HashUsed,hash)) ;{
    sol = <-self.ChNonce
  }
}

// Builds a job and returns the evaluation of an individual,
// without start executing the genetic algorithm.
func (self *Job)EvaluateSingleSolution(individual []byte,
      hashPrev,publicKey string)*Sol{
  hb := new(utils.HashBuilder)
  hb.Add(hashPrev)
  hb.Add(publicKey)
  hash := hb.GetHash()
  x := int64(binary.BigEndian.Uint64(hash[:8]))
  var prng *rand.Rand = rand.New(rand.NewSource(99))
  prng.Seed(x)
  st := op.MakeState(hash)
  sol := new(Sol)
  sol.Individual = self.dna.Generate(prng)
  sol.Individual.LoadFromSerialization(individual)
  sol.Fitness = 9999999999999
  sol.IsEval = false
  sol.eval(st,hash)
  return sol
}
