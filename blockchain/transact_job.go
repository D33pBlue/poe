/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-23
 * @Project: Proof of Evolution
 * @Filename: std_trans.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-30
 * @Copyright: 2020
 */

package blockchain

import(
  "fmt"
  "time"
  // "errors"
  "github.com/D33pBlue/poe/utils"
)

type JobTransaction struct{
  Timestamp time.Time
  Inputs []TrInput
  Creator utils.Addr
  Job string
  Prize int
  Hash string
  Signature string
}

// func MakeJobTransaction(creator utils.Addr,key utils.Key,
//       inps []TrInput,outs []TrOutput)(*JobTransaction,error){
//   tr := new(JobTransaction)
//   tr.Timestamp = time.Now()
//   tr.Creator = creator
//   tr.Inputs = inps
//   tr.Outputs = outs
//   tr.Hash = tr.GetHash()
//   tr.Signature = utils.GetSignatureFromHash(tr.Hash,key)
//  return tr,nil
// }

func (self *JobTransaction)Check(chain *Blockchain)bool{
  return true // TODO: implement later
}

func (self *JobTransaction)IsSpent()bool{
  return false // TODO: implement later
}

func (self *JobTransaction)GetHash()string{
  hb := new(utils.HashBuilder)
  hb.Add(self.Creator)
  hb.Add(self.Timestamp)
  for i:=0;i<len(self.Inputs);i++{
    hb.Add(self.Inputs[i].Block)
    hb.Add(self.Inputs[i].ToSpend)
    hb.Add(self.Inputs[i].Index)
  }
  hb.Add(self.Job)
  hb.Add(self.Prize)
  return fmt.Sprintf("%x",hb.GetHash())
}

func (self *JobTransaction)GetHashCached()string{
  return self.Hash
}

func (self *JobTransaction)Serialize()[]byte{
  return nil // TODO:  implement later
}

func MarshalJobTransaction([] byte)*JobTransaction{
  return nil // TODO: implement later
}

func (self *JobTransaction)GetType()string{
  return TrJob
}
