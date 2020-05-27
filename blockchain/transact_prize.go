/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-May-25
 * @Project: Proof of Evolution
 * @Filename: transact_prize.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-27
 * @Copyright: 2020
 */

package blockchain

import(
  "time"
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

func (self *PrizeTransaction)Check(block *Block,trChanges *map[string]string)bool{
  return false // TODO: implement later
}

func (self *PrizeTransaction)GetHash()string{
  return "" // TODO: implement later
}

func (self *PrizeTransaction)GetHashCached()string{
  return self.Hash
}

func (self *PrizeTransaction)GetCreator()utils.Addr{
  return self.Creator
}

func (self *PrizeTransaction)GetTimestamp()time.Time{
  return self.Timestamp
}

func (self *PrizeTransaction)GetType()string{
  return TrPrize
}

func (self *PrizeTransaction)GetOutputAt(int)*TrOutput{
  return &self.Output
}

func (self *PrizeTransaction)Serialize()[]byte{
  return nil // TODO: implement later
}

func MarshalPrizeTransaction(data []byte)*PrizeTransaction{
  return nil // TODO: implement later
}
