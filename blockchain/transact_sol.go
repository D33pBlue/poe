/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-May-22
 * @Project: Proof of Evolution
 * @Filename: transact_sol.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-22
 * @Copyright: 2020
 */

package blockchain

import(
  "time"
  "github.com/D33pBlue/poe/utils"
)

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

func (self *SolTransaction)Check(block *Block,trChanges *map[string]string)bool{
  return false // TODO: implement later
}

func (self *SolTransaction)GetHash()string{
  return "" // TODO: implement later
}

func (self *SolTransaction)GetHashCached()string{
  return self.Hash
}

func (self *SolTransaction)GetCreator()utils.Addr{
  return self.Creator
}

func (self *SolTransaction)GetTimestamp()time.Time{
  return self.Timestamp
}

func (self *SolTransaction)GetType()string{
  return TrSol
}

func (self *SolTransaction)GetOutputAt(int)*TrOutput{
  return &self.Output
}

func (self *SolTransaction)Serialize()[]byte{
  return nil // TODO: implement later
}

func MarshalSolTransaction([]byte)*SolTransaction{
  return nil // TODO: implement later
}
