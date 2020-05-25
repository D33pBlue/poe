/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-May-22
 * @Project: Proof of Evolution
 * @Filename: transact_sol.go
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

// A SolTransaction is used to receive money discosing a solution
// that was previously declared with a ResTransaction. It is valid
// if the solution is correct and the creator deserves the prize.
type SolTransaction struct{
  Timestamp time.Time
  Output TrOutput
  Creator utils.Addr
  ResBlock string // the hash of the block with the ResTransaction
  ResTrans string // the hash of the ResTransaction
  JobTrans string // the hash of the JobTransaction
  Solution []byte
  Hash string
  Signature string
}

// Builds a new SolTransaction and signs it. This method does not check
// the data it receives in input.
func MakeSolTransaction(creator utils.Addr,key utils.Key,
      resblock,restrans,jobtrans string,
      solution []byte,amount int)*SolTransaction{
  tr := new(SolTransaction)
  tr.Timestamp = time.Now()
  tr.Output.Address = creator
  tr.Output.Value = amount
  tr.Creator = creator
  tr.ResBlock = resblock
  tr.ResTrans = restrans
  tr.JobTrans = jobtrans
  tr.Solution = solution
  tr.Hash = tr.GetHash()
  tr.Signature = fmt.Sprintf("%x",utils.GetSignatureFromHash(tr.Hash,key))
 return tr
}

func (self *SolTransaction)Check(block *Block,trChanges *map[string]string)bool{
  return false // TODO: implement later
}

// Recalculates and return the hash of the transaction.
func (self *SolTransaction)GetHash()string{
  hb := new(utils.HashBuilder)
  hb.Add(self.Creator)
  hb.Add(self.Timestamp.Format("2006-01-02 15:04:05"))
  hb.Add(self.ResBlock)
  hb.Add(self.ResTrans)
  hb.Add(self.Output)
  hb.Add(self.JobTrans)
  hb.Add(self.Solution)
  return fmt.Sprintf("%x",hb.GetHash())
}

// Returns the cached hash of the transaction.
func (self *SolTransaction)GetHashCached()string{
  return self.Hash
}

// Returns the public key of the creator of this transaction.
func (self *SolTransaction)GetCreator()utils.Addr{
  return self.Creator
}

// Returns the timestamp stored in the transaction.
func (self *SolTransaction)GetTimestamp()time.Time{
  return self.Timestamp
}

// Returns the type of the transaction, which is "SolTransaction"
func (self *SolTransaction)GetType()string{
  return TrSol
}

// Always returns the TrOutput stored in the transaction,
// ignoring the parameter in input
func (self *SolTransaction)GetOutputAt(i int)*TrOutput{
  return &self.Output
}

// Serializes the SolTransaction and returns it as []byte.
func (self *SolTransaction)Serialize()[]byte{
  data, err := json.Marshal(self)
  if err != nil {
    fmt.Println(err)
  }
  return data
}

// Rebuilds a SolTransaction from its serialization and returns
// it, or nil.
func MarshalSolTransaction(data []byte)*SolTransaction{
  var objmap map[string]json.RawMessage
  json.Unmarshal(data, &objmap)
  tr := new(SolTransaction)
  json.Unmarshal(objmap["Timestamp"],&tr.Timestamp)
  json.Unmarshal(objmap["Output"],&tr.Output)
  json.Unmarshal(objmap["Creator"],&tr.Creator)
  json.Unmarshal(objmap["ResBlock"],&tr.ResBlock)
  json.Unmarshal(objmap["ResTrans"],&tr.ResTrans)
  json.Unmarshal(objmap["JobTrans"],&tr.JobTrans)
  json.Unmarshal(objmap["Solution"],&tr.Solution)
  json.Unmarshal(objmap["Hash"],&tr.Hash)
  json.Unmarshal(objmap["Signature"],&tr.Signature)
  return tr
}
