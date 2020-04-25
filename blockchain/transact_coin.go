/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-23
 * @Project: Proof of Evolution
 * @Filename: std_trans.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-25
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
  Outputs []TrOutput
  Creator utils.Addr
  Hash []byte
  Signature []byte
}

func MakeCoinTransaction(creator utils.Addr,key utils.Key,
          outs []TrOutput)(*CoinTransaction,error){
  tr := new(CoinTransaction)
  tr.Timestamp = time.Now()
  tr.Creator = creator
  tr.Outputs = outs
  tr.Hash = tr.GetHash()
  tr.Signature = utils.GetSignatureFromHash(tr.Hash,key)
  return tr,nil
}

func (self *CoinTransaction)Check(chain *Blockchain)bool{
  // TODO: implement later
  return true
}

func (self *CoinTransaction)IsSpent()bool{
  return false // TODO: implement later
}

func (self *CoinTransaction)GetHash()[]byte{
  hb := new(utils.HashBuilder)
  hb.Add(self.Creator)
  hb.Add(self.Timestamp)
  for i:=0;i<len(self.Outputs);i++{
    hb.Add(self.Outputs[i].Address)
    hb.Add(self.Outputs[i].Value)
  }
  return hb.GetHash()
}

func (self *CoinTransaction)GetType()string{
  return TrCoin
}
