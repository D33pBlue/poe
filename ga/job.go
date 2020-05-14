/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-May-14
 * @Project: Proof of Evolution
 * @Filename: job.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-14
 * @Copyright: 2020
 */
package ga

type Job struct{
  ChNonce chan Sol // used to send nonce candidates
  ChUpdateIn chan Sol // used to receive shared solutions from miners
  ChUpdateOut chan Sol // used to share good solutions to miners
  Population Population
  KeepRunning bool
  BestFound []Sol
}

func BuildJob(jobpath,datapath string)*Job{
  // initialize the channels with buffers in order to made them async
  return nil // TODO: implement later
}

func (self *Job)Execute(){
  // TODO: implement later
}
