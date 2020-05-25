/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-May-22
 * @Project: Proof of Evolution
 * @Filename: transact_res.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-23
 * @Copyright: 2020
 */

package blockchain

import(
  "fmt"
  "time"
  "encoding/json"
  "github.com/D33pBlue/poe/utils"
)

// A ResTransaction is used to declare a good result achieved. The best
// ones should receive a prize. In order to receive it the miner has to
// submit a SolTransaction in the next block after the block with this
// transaction is mined. This transaction is valid only if it is submitted
// in the proper slot (that is scheduled with the job execution).
type ResTransaction struct{
  Timestamp time.Time
  Inputs []TrInput
  Output TrOutput // the possible remainder
  Creator utils.Addr
  JobBlock string // the hash of the block with the JobTransaction
  JobTrans string // the hash of the JobTransaction
  Evaluation float64
  HashSol string // the hash of solution: [individual,evaluation,creator]
  Hash string
  Signature string
}

// Builds a new ResTransaction and signs it. This method does not check
// the data it receives in input.
func MakeResTransaction(creator utils.Addr,key utils.Key,
      inps []TrInput,out TrOutput,
      jobblock,jobtrans,hashsol string,
      evaluation float64)*ResTransaction{
  tr := new(ResTransaction)
  tr.Timestamp = time.Now()
  tr.Output = out
  tr.Inputs = inps
  tr.Creator = creator
  tr.JobBlock = jobblock
  tr.JobTrans = jobtrans
  tr.HashSol = hashsol
  tr.Evaluation = evaluation
  tr.Hash = tr.GetHash()
  tr.Signature = fmt.Sprintf("%x",utils.GetSignatureFromHash(tr.Hash,key))
 return tr
}

func (self *ResTransaction)Check(block *Block,trChanges *map[string]string)bool{
  return false // TODO: implement later
}

// Recalculates and return the hash of the transaction.
func (self *ResTransaction)GetHash()string{
  hb := new(utils.HashBuilder)
  hb.Add(self.Creator)
  hb.Add(self.Timestamp.Format("2006-01-02 15:04:05"))
  hb.Add(self.JobBlock)
  hb.Add(self.JobTrans)
  hb.Add(self.Output)
  for i:=0;i<len(self.Inputs);i++{
    hb.Add(self.Inputs[i])
  }
  hb.Add(self.HashSol)
  hb.Add(self.Evaluation)
  return fmt.Sprintf("%x",hb.GetHash())
}

// Returns the cached hash of the transaction.
func (self *ResTransaction)GetHashCached()string{
  return self.Hash
}

// Returns the public key of the creator of the transaction.
func (self *ResTransaction)GetCreator()utils.Addr{
  return self.Creator
}

// Returns the timestamp stored in the transaction.
func (self *ResTransaction)GetTimestamp()time.Time{
  return self.Timestamp
}

// Returns the type string of the transaction, which is "ResTransaction".
func (self *ResTransaction)GetType()string{
  return TrRes
}

// Always returns the only Output stored in the transaction.
func (self *ResTransaction)GetOutputAt(i int)*TrOutput{
  return &self.Output
}

// Serializes the transactions and returns it as []byte
func (self *ResTransaction)Serialize()[]byte{
  data, err := json.Marshal(self)
  if err != nil {
    fmt.Println(err)
  }
  return data
}

// Rebuilds the transaction from its serialization.
func MarshalResTransaction(data []byte)*ResTransaction{
  var objmap map[string]json.RawMessage
  json.Unmarshal(data, &objmap)
  tr := new(ResTransaction)
  json.Unmarshal(objmap["Timestamp"],&tr.Timestamp)
  json.Unmarshal(objmap["Output"],&tr.Output)
  json.Unmarshal(objmap["Inputs"],&tr.Inputs)
  json.Unmarshal(objmap["Creator"],&tr.Creator)
  json.Unmarshal(objmap["JobBlock"],&tr.JobBlock)
  json.Unmarshal(objmap["JobTrans"],&tr.JobTrans)
  json.Unmarshal(objmap["HashSol"],&tr.HashSol)
  json.Unmarshal(objmap["Evaluation"],&tr.Evaluation)
  json.Unmarshal(objmap["Hash"],&tr.Hash)
  json.Unmarshal(objmap["Signature"],&tr.Signature)
  return tr
}
