/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-23
 * @Project: Proof of Evolution
 * @Filename: std_trans.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-08
 * @Copyright: 2020
 */

package blockchain

import(
  "fmt"
  "time"
  "encoding/json"
  "github.com/D33pBlue/poe/utils"
)

type CoinTransaction struct{
  Timestamp time.Time
  Output TrOutput
  Creator utils.Addr
  Hash string
  spent bool
}

func MakeCoinTransaction(receiver utils.Addr,value int)(*CoinTransaction,error){
  tr := new(CoinTransaction)
  tr.Timestamp = time.Now()
  out := new(TrOutput)
  out.Address = receiver
  out.Value = value
  tr.Output = *out
  tr.Creator = receiver
  tr.Hash = tr.GetHash()
  tr.spent = false
  return tr,nil
}

func (self *CoinTransaction)Check(block *Block,trChanges *map[string]string)bool{
  return true // TODO: implement later
}

func (self* CoinTransaction)GetOutputAt(i int)*TrOutput{
  return &self.Output
}

func (self *CoinTransaction)GetTimestamp()time.Time{
  return self.Timestamp
}

func (self *CoinTransaction)GetCreator()utils.Addr{
  return self.Creator
}

func (self *CoinTransaction)GetHash()string{
  hb := new(utils.HashBuilder)
  hb.Add(self.Timestamp)
  hb.Add(self.Creator)
  hb.Add(self.Output.Address)
  hb.Add(self.Output.Value)
  return fmt.Sprintf("%x",hb.GetHash())
}

func (self *CoinTransaction)GetHashCached()string{
  return self.Hash
}

func MarshalCoinTransaction(data []byte)*CoinTransaction{
  var objmap map[string]json.RawMessage
  json.Unmarshal(data, &objmap)
  tr := new(CoinTransaction)
  json.Unmarshal(objmap["Timestamp"],&tr.Timestamp)
  json.Unmarshal(objmap["Output"],&tr.Output)
  json.Unmarshal(objmap["Hash"],&tr.Hash)
  tr.spent = false
  return tr
}

func (self *CoinTransaction)GetType()string{
  return TrCoin
}
