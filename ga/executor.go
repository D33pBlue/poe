/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-May-14
 * @Project: Proof of Evolution
 * @Filename: executor.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-14
 * @Copyright: 2020
 */

package ga

type JobChannels struct{
  ChNonce chan Sol // used to send nonce candidates
  ChUpdateIn chan Sol // used to receive shared solutions from miners
  ChUpdateOut chan Sol // used to share good solutions to miners
}

type Executor struct{
  ActiveJobs map[string]*Job // the key is the hash of the JobTransaction
                            // in which the job is defined
}

// true <=> the job with a hash has been already defined and is
// still active.
func (self *Executor)IsExecutingJob(job string)bool{
  return false // TODO: implement later
}

// Terminates the execution of a job by its hash.
func (self *Executor)StopJob(job string){
  // TODO: implement later
}

// Initialize and start executing a job from the paths to the files
// with its definition and data. If the job runs correctly, this methos
// returns a JobChannels with the channels to communicate with the job;
// otherwise nil.
func (self *Executor)StartJob(hash,jobpath,datapath string)*JobChannels{
  return nil // TODO: implement later
}

// Returns the channels of an already running job if it is found
// by its hash; nil otherwise.
func (self *Executor)GetChannels(job string)*JobChannels{
  return nil // TODO: implement later
}

// Sends a good solution to an active job so that it can include it
// in his population.
func (self *Executor)InjectSharedSolution(job string,sol Sol){
  chs := self.GetChannels(job)
  if chs!=nil{
    chs.ChUpdateIn <- sol
  }
}
