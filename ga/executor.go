/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-May-14
 * @Project: Proof of Evolution
 * @Filename: executor.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-27
 * @Copyright: 2020
 */

package ga

type JobChannels struct{
  ChNonce chan Sol // used to send nonce candidates
  ChUpdateIn chan Sol // used to receive shared solutions from miners
}

type Executor struct{
  ActiveJobs map[string]*Job // the key is the hash of the JobTransaction
                            // in which the job is defined
  ChUpdateOut chan Sol
}

// Builds and initialize an Executor
func BuildExecutor(chUpdateOut chan Sol)*Executor{
  executor := new(Executor)
  executor.ActiveJobs = make(map[string]*Job)
  executor.ChUpdateOut = chUpdateOut
  return executor
}

// true <=> the job with a hash has been already defined
func (self *Executor)IsExecutingJob(job string)bool{
  if val, ok := self.ActiveJobs[job]; ok {
    if val==nil { return false }
    return true
  }
  return false
}

// Terminates the execution of a job by its hash.
func (self *Executor)StopJob(job string){
  if self.IsExecutingJob(job){
    self.ActiveJobs[job].KeepRunning = false
    self.ActiveJobs[job] = nil
    delete(self.ActiveJobs,job)
  }
}

// Initialize and start executing a job from the paths to the files
// with its definition and data. If the job runs correctly, this methos
// returns a JobChannels with the channels to communicate with the job;
// otherwise nil.
func (self *Executor)StartJob(hash,hashPrev,publicKey,
      jobpath,datapath string)*JobChannels{
  job := BuildJob(jobpath,datapath,self.ChUpdateOut,hash)
  if job==nil{ return nil }
  self.ActiveJobs[hash] = job
  go job.Execute(hashPrev,publicKey)
  chs := new(JobChannels)
  chs.ChNonce = job.ChNonce
  chs.ChUpdateIn = job.ChUpdateIn
  return chs
}

// Returns the channels of an already running job if it is found
// by its hash; nil otherwise.
func (self *Executor)GetChannels(job string)*JobChannels{
  if self.IsExecutingJob(job){
    executing := self.ActiveJobs[job]
    chs := new(JobChannels)
    chs.ChNonce = executing.ChNonce
    chs.ChUpdateIn = executing.ChUpdateIn
    return chs
  }
  return nil
}

// Change the hash of the block in the job's execution configuration,
// so that the coefficients for the complexity are updated.
// The solutions in ChNonce channel are resetted.
func (self *Executor)ChangeBlockHashInJob(job,hashPrev,publicKey string){
  if self.IsExecutingJob(job){
    executing := self.ActiveJobs[job]
    executing.ChangeBlockHash(hashPrev,publicKey)
  }
}

// Sends a good solution to an active job so that it can include it
// in his population.
func (self *Executor)InjectSharedSolution(job string,solSerialized []byte){
  var sol Sol
  sol.Individual.LoadFromSerialization(solSerialized)
  chs := self.GetChannels(job)
  if chs!=nil{
    chs.ChUpdateIn <- sol
  }
}
