/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-May-25
 * @Project: Proof of Evolution
 * @Filename: transact_prize.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-28
 * @Copyright: 2020
 */

package blockchain

import(
  "fmt"
  "time"
  "encoding/json"
  "github.com/D33pBlue/poe/utils"
)

type PrizeTransaction struct{
  Timestamp time.Time
  Output TrOutput
  Creator utils.Addr
  SolBlock string // the hash of the block with the SolTransaction
  SolTrans string // the hash of the SolTransaction
  JobTrans string // the hash of the JobTransaction
  Hash string
}

// Builds a new PrizeTransaction and returns it.
func MakePrizeTransaction(receiver utils.Addr,amount int,
      solBlock,solTrans,jobTrans string)*PrizeTransaction{
  tr := new(PrizeTransaction)
  tr.Timestamp = time.Now()
  tr.Creator = receiver
  tr.SolBlock = solBlock
  tr.SolTrans = solTrans
  tr.JobTrans = jobTrans
  tr.Output.Address = receiver
  tr.Output.Value = amount
  tr.Hash = tr.GetHash()
  return tr
}

// Check validate the transaction and update trChanges. The parameter block
// is assumed to be the block in which this transaction is stored.
// The validity is checked in relation to the chain that is linked
// to that block. The subchain must be valid/already checked.
// In order to be valid, the transaction must have:
// - the hash that matches the one declared
// - a valid link to a SolTransaction that belongs to the Creator, and that is
// stored in the previous block
// - the hash of the JobTransaction that match the declaration
// This method does not check the amount of the prize.
func (self *PrizeTransaction)Check(block *Block,trChanges *map[string]string)bool{
  hash2 := self.GetHash()
  if hash2!=self.Hash{
    fmt.Println("The hash does not match")
    fmt.Printf("%v !=\n%v\n",hash2,self.Hash)
    return false}
  if self.SolBlock!=block.Previous.GetHashCached(){
    fmt.Println("The ResTransaction in this SolTransaction is not in Previous block")
    return false }
  tr := block.Previous.FindTransaction(self.SolTrans)
  if tr==nil{
    fmt.Println("The declared ResTransaction does not exists")
    return false }
  if self.Creator!=tr.GetCreator() {
    fmt.Println("The linked ResTransaction belongs to another miner")
    return false }
  if self.Creator!=self.Output.Address {
    fmt.Println("The money are sent to the wrong miner")
    return false }
  solTr := tr.(*SolTransaction)
  if solTr.JobTrans!=self.JobTrans{
    fmt.Println("The jobs are different")
    return false }
  return true
}

// Recalculates and returns the hash of the transaction.
func (self *PrizeTransaction)GetHash()string{
  hb := new(utils.HashBuilder)
  hb.Add(self.Creator)
  hb.Add(self.Timestamp.Format("2006-01-02 15:04:05"))
  hb.Add(self.SolBlock)
  hb.Add(self.SolTrans)
  hb.Add(self.JobTrans)
  hb.Add(self.Output)
  return fmt.Sprintf("%x",hb.GetHash())
}

// Returns the cached hash of the transaction.
func (self *PrizeTransaction)GetHashCached()string{
  return self.Hash
}

// Returns the public key of the receiver of the prize.
func (self *PrizeTransaction)GetCreator()utils.Addr{
  return self.Creator
}

// Returns the timestamp stored inside the transaction.
func (self *PrizeTransaction)GetTimestamp()time.Time{
  return self.Timestamp
}

// Returns the type of the transaction (TrPrize).
func (self *PrizeTransaction)GetType()string{
  return TrPrize
}

// Always returns a reference to self.Output.
func (self *PrizeTransaction)GetOutputAt(int)*TrOutput{
  return &self.Output
}

// Returns the serialization of the transaction (as []byte).
func (self *PrizeTransaction)Serialize()[]byte{
  data, err := json.Marshal(self)
  if err != nil {
    fmt.Println(err)
  }
  return data
}

// Rebuilds a PrizeTransaction from its serialization.
func MarshalPrizeTransaction(data []byte)*PrizeTransaction{
  var objmap map[string]json.RawMessage
  json.Unmarshal(data, &objmap)
  tr := new(PrizeTransaction)
  json.Unmarshal(objmap["Timestamp"],&tr.Timestamp)
  json.Unmarshal(objmap["Creator"],&tr.Creator)
  json.Unmarshal(objmap["SolBlock"],&tr.SolBlock)
  json.Unmarshal(objmap["SolTrans"],&tr.SolTrans)
  json.Unmarshal(objmap["JobTrans"],&tr.JobTrans)
  json.Unmarshal(objmap["Output"],&tr.Output)
  json.Unmarshal(objmap["Hash"],&tr.Hash)
  return tr
}
