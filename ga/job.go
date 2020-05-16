/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-May-14
 * @Project: Proof of Evolution
 * @Filename: job.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-16
 * @Copyright: 2020
 */
package ga

import(
  "fmt"
  "github.com/D33pBlue/poe/utils"
)

type Job struct{
  ChNonce chan Sol // used to send nonce candidates
  ChUpdateIn chan Sol // used to receive shared solutions from miners
  ChUpdateOut chan Sol // used to share good solutions to miners
  KeepRunning bool
  dna DNA
}

func BuildJob(jobpath,datapath string)*Job{
  // initialize the channels with buffers in order to made them async
  job := new(Job)
  job.ChNonce = make(chan Sol,1000)
  job.ChUpdateIn = make(chan Sol,100)
  job.ChUpdateOut = make(chan Sol, 10)
  // compile and load DNA
  var err error
  job.dna,err = LoadGA(jobpath,datapath)
  if err!=nil{
    fmt.Println(err)
    return nil
  }
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
  conf := BuildBlockchainGAConfig(hash,&self.KeepRunning,100)
  RunGA(self.dna,conf,self.ChUpdateOut,self.ChUpdateIn,self.ChNonce)
}
