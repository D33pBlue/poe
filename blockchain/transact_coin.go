/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-23
 * @Project: Proof of Evolution
 * @Filename: std_trans.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-19
 * @Copyright: 2020
 */

package blockchain

import(
  "fmt"
  "time"
  "encoding/json"
  "github.com/D33pBlue/poe/utils"
)

// A CoinTransaction is the transaction miners can use to earn some
// money when they successfully complete the minig process.
type CoinTransaction struct{
  Timestamp time.Time
  Output TrOutput
  Creator utils.Addr
  Hash string
  spent bool
}

// Instantiates a new CoinTransaction.
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

// Returns always true because a CoinTransaction never spends money.
func (self *CoinTransaction)Check(block *Block,trChanges *map[string]string)bool{
  return true
}

//Returns always the only TrOutput it stores.
func (self* CoinTransaction)GetOutputAt(i int)*TrOutput{
  return &self.Output
}

// Returns the timestamp stored in the transaction.
func (self *CoinTransaction)GetTimestamp()time.Time{
  return self.Timestamp
}

// Returns the public key of the creator of the transaction.
func (self *CoinTransaction)GetCreator()utils.Addr{
  return self.Creator
}

// Recalculates the hash of the transaction.
func (self *CoinTransaction)GetHash()string{
  hb := new(utils.HashBuilder)
  hb.Add(self.Timestamp.Format("2006-01-02 15:04:05"))
  hb.Add(self.Creator)
  hb.Add(self.Output.Address)
  hb.Add(self.Output.Value)
  return fmt.Sprintf("%x",hb.GetHash())
}

// Returns the cached hash of the transaction.
func (self *CoinTransaction)GetHashCached()string{
  return self.Hash
}

// Serializes the transaction and returns it as []byte
func (self *CoinTransaction)Serialize()[]byte{
  data, err := json.Marshal(self)
  if err != nil {
    fmt.Println(err)
  }
  return data
}

// Rebuilds the CoinTransaction from its serialized data.
func MarshalCoinTransaction(data []byte)*CoinTransaction{
  var objmap map[string]json.RawMessage
  json.Unmarshal(data, &objmap)
  tr := new(CoinTransaction)
  json.Unmarshal(objmap["Timestamp"],&tr.Timestamp)
  json.Unmarshal(objmap["Output"],&tr.Output)
  json.Unmarshal(objmap["Hash"],&tr.Hash)
  json.Unmarshal(objmap["Creator"],&tr.Creator)
  tr.spent = false
  return tr
}

func (self *CoinTransaction)GetType()string{
  return TrCoin
}
