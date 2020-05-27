/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-May-22
 * @Project: Proof of Evolution
 * @Filename: transact_sol.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-27
 * @Copyright: 2020
 */

package blockchain

import(
  "fmt"
  "time"
  "encoding/json"
  "github.com/D33pBlue/poe/utils"
  "github.com/D33pBlue/poe/conf"
  "github.com/D33pBlue/poe/ga"
)

// A SolTransaction is used to receive money discosing a solution
// that was previously declared with a ResTransaction. It is valid
// if the solution is correct and the creator deserves the prize.
type SolTransaction struct{
  Timestamp time.Time
  Creator utils.Addr
  ResBlock string // the hash of the block with the ResTransaction
  ResTrans string // the hash of the ResTransaction
  JobTrans string // the hash of the JobTransaction
  Solution []byte
  Hash string
  Signature string
  config *conf.Config
}

// Builds a new SolTransaction and signs it. This method does not check
// the data it receives in input.
func MakeSolTransaction(creator utils.Addr,key utils.Key,
      resblock,restrans,jobtrans string,
      solution []byte,config *conf.Config)*SolTransaction{
  tr := new(SolTransaction)
  tr.Timestamp = time.Now()
  tr.Creator = creator
  tr.ResBlock = resblock
  tr.ResTrans = restrans
  tr.JobTrans = jobtrans
  tr.Solution = solution
  tr.Hash = tr.GetHash()
  tr.Signature = fmt.Sprintf("%x",utils.GetSignatureFromHash(tr.Hash,key))
  tr.config = config
 return tr
}

// Check validate the transaction and update trChanges. The parameter block
// is assumed to be the block in which this transaction is stored.
// The validity is checked in relation to the chain that is linked
// to that block. The subchain must be valid/already checked.
// In order to be valid, the transaction must have:
// - the hash that matches the one declared
// - the signature verified with the public key of the creator
// - a valid link to a ResTransaction that belongs to the Creator, and that is
// stored in the previous block
// - the hash of the JobTransaction that match the declaration
// - a solution whose hash matches the declaration
// - the correct evaluation of the solution
// In addition, this transaction should not be already stored in the block.
func (self *SolTransaction)Check(block *Block,trChanges *map[string]string)bool{
  hash2 := self.GetHash()
  if hash2!=self.Hash{
    fmt.Println("The hash does not match")
    fmt.Printf("%v !=\n%v\n",hash2,self.Hash)
    return false}
  if !utils.CheckSignature(self.Signature,self.Hash,self.Creator){
    fmt.Println("The signature is not valid")
    return false}
  if self.ResBlock!=block.Previous.GetHashCached(){
    fmt.Println("The ResTransaction in this SolTransaction is not in Previous block")
    return false }
  tr := block.Previous.FindTransaction(self.ResTrans)
  if tr==nil{
    fmt.Println("The declared ResTransaction does not exists")
    return false }
  if self.Creator!=tr.GetCreator(){
    fmt.Println("The linked ResTransaction belongs to another miner")
    return false }
  resTr := tr.(*ResTransaction)
  if resTr.JobTrans!=self.JobTrans{
    fmt.Println("The jobs are different")
    return false }
  hb := new(utils.HashBuilder)
  hb.Add(self.Solution)
  hashSol := fmt.Sprintf("%x",hb.GetHash())
  if hashSol!=resTr.HashSol{
    fmt.Println("The solution's hash does not match the one declared")
    return false}
  jobBlock := block.FindPrevBlock(resTr.JobBlock)
  if jobBlock==nil{
    fmt.Println("Unable to find the block with job transaction")
    return false}
  transaction := jobBlock.FindTransaction(resTr.JobTrans)
  if transaction==nil{
    fmt.Println("Unable to find the job transaction")
    return false}
  jobTr := transaction.(*JobTransaction)
  jobPath,dataPath := self.config.GetSuitablePathForJob(self.JobTrans)
  err := jobTr.SaveJobInFile(jobPath)
  if err!=nil{
    fmt.Println(err)
    return false}
  err2 := jobTr.SaveDataInFile(dataPath)
  if err2!=nil{
    fmt.Println(err2)
    return false}
  job := ga.BuildJob(jobPath,dataPath,nil,"")
  if job==nil{
    fmt.Println("Unable to build job")
    return false}
  sol := job.EvaluateSingleSolution(self.Solution,"","")
  if sol.Fitness!=resTr.Evaluation{
    fmt.Println("The declared evaluation of the solution is wrong")
    return false}
  transactions := block.Transactions.GetTransactionArray()
  for i:=0;i<len(transactions);i++{
    if transactions[i].GetCreator()==self.Creator && transactions[i].GetType()==TrSol{
      solTr := transactions[i].(*SolTransaction)
      if self.ResTrans==solTr.ResTrans && solTr!=self{
        fmt.Println("A SolTransaction already exists for the same ResTransaction")
        return false
      }
    }
  }
  return true
}

// Recalculates and return the hash of the transaction.
func (self *SolTransaction)GetHash()string{
  hb := new(utils.HashBuilder)
  hb.Add(self.Creator)
  hb.Add(self.Timestamp.Format("2006-01-02 15:04:05"))
  hb.Add(self.ResBlock)
  hb.Add(self.ResTrans)
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

// Always returns nil
func (self *SolTransaction)GetOutputAt(i int)*TrOutput{
  return nil
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
func MarshalSolTransaction(data []byte,config *conf.Config)*SolTransaction{
  var objmap map[string]json.RawMessage
  json.Unmarshal(data, &objmap)
  tr := new(SolTransaction)
  json.Unmarshal(objmap["Timestamp"],&tr.Timestamp)
  json.Unmarshal(objmap["Creator"],&tr.Creator)
  json.Unmarshal(objmap["ResBlock"],&tr.ResBlock)
  json.Unmarshal(objmap["ResTrans"],&tr.ResTrans)
  json.Unmarshal(objmap["JobTrans"],&tr.JobTrans)
  json.Unmarshal(objmap["Solution"],&tr.Solution)
  json.Unmarshal(objmap["Hash"],&tr.Hash)
  json.Unmarshal(objmap["Signature"],&tr.Signature)
  tr.config = config
  return tr
}
