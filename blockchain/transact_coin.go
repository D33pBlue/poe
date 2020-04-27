/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-23
 * @Project: Proof of Evolution
 * @Filename: std_trans.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-27
 * @Copyright: 2020
 */

package blockchain

import(
  "time"
  // "errors"
  "github.com/D33pBlue/poe/utils"
)

type CoinTransaction struct{
  Timestamp time.Time
  Output TrOutput
  Hash []byte
}

func MakeCoinTransaction(receiver utils.Addr,value int)(*CoinTransaction,error){
  tr := new(CoinTransaction)
  tr.Timestamp = time.Now()
  out := new(TrOutput)
  out.Address = receiver
  out.Value = value
  tr.Output = *out
  tr.Hash = tr.GetHash()
  return tr,nil
}

func (self *CoinTransaction)Check(chain *Blockchain)bool{
  return utils.CompareHashes(self.Hash,self.GetHash())
}

func (self *CoinTransaction)IsSpent()bool{
  return false // TODO: implement later
}

func (self *CoinTransaction)GetHash()[]byte{
  hb := new(utils.HashBuilder)
  hb.Add(self.Timestamp)
  hb.Add(self.Output.Address)
  hb.Add(self.Output.Value)
  return hb.GetHash()
}

func (self *CoinTransaction)GetHashCached()[]byte{
  return self.Hash
}

func (self *CoinTransaction)Serialize()[]byte{
  return nil // TODO:  implement later
}

func MarshalCoinTransaction([] byte)*CoinTransaction{
  return nil // TODO: implement later
}

func (self *CoinTransaction)GetType()string{
  return TrCoin
}
