/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-23
 * @Project: Proof of Evolution
 * @Filename: std_trans.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-01
 * @Copyright: 2020
 */

package blockchain

import(
  "fmt"
  "time"
  "encoding/json"
  "github.com/D33pBlue/poe/utils"
)

type StdTransaction struct{
  Timestamp time.Time
  Inputs []TrInput
  Outputs []TrOutput
  Creator utils.Addr
  Hash string
  Signature string
  spent bool
}

func MakeStdTransaction(creator utils.Addr,key utils.Key,
      inps []TrInput,outs []TrOutput)(*StdTransaction,error){
  tr := new(StdTransaction)
  tr.Timestamp = time.Now()
  tr.Creator = creator
  tr.Inputs = inps
  tr.Outputs = outs
  for i:=0;i<len(tr.Outputs);i++{
    tr.Outputs[i].spent = false
  }
  tr.Hash = tr.GetHash()
  tr.Signature = fmt.Sprintf("%x",utils.GetSignatureFromHash(tr.Hash,key))
 return tr,nil
}

func (self *StdTransaction)Check(chain *Blockchain)bool{
  hash2 := self.GetHash()
  if hash2!=self.Hash{return false}
  if !utils.CheckSignature(self.Signature,self.Hash,self.Creator){
    return false
  }
  var tot int = 0
  for i:=0;i<len(self.Inputs);i++{
    tot += self.Inputs[i].GetValue(self.Creator,chain)
  }
  var spent int = 0
  for i:=0;i<len(self.Outputs);i++{
    spent += self.Outputs[i].Value
  }
  return tot>=spent
}

func (self *StdTransaction)GetCreator()utils.Addr{
  return self.Creator
}

func (self *StdTransaction)IsSpent()bool{
  return false // TODO: implement later
}

func (self *StdTransaction)SetSpent(){
  self.spent = true
}

func (self *StdTransaction)GetHash()string{
  hb := new(utils.HashBuilder)
  hb.Add(self.Creator)
  hb.Add(self.Timestamp)
  for i:=0;i<len(self.Inputs);i++{
    hb.Add(self.Inputs[i].Block)
    hb.Add(self.Inputs[i].ToSpend)
    hb.Add(self.Inputs[i].Index)
  }
  for i:=0;i<len(self.Outputs);i++{
    hb.Add(self.Outputs[i].Address)
    hb.Add(self.Outputs[i].Value)
  }
  return fmt.Sprintf("%x",hb.GetHash())
}

func (self *StdTransaction)GetHashCached()string{
  return self.Hash
}

// func (self *StdTransaction)Serialize()[]byte{
//   return nil // TODO:  implement later
// }

func MarshalStdTransaction(data []byte)*StdTransaction{
  var objmap map[string]json.RawMessage
  json.Unmarshal(data, &objmap)
  tr := new(StdTransaction)
  json.Unmarshal(objmap["Timestamp"],&tr.Timestamp)
  json.Unmarshal(objmap["Inputs"],&tr.Inputs)
  json.Unmarshal(objmap["Outputs"],&tr.Outputs)
  json.Unmarshal(objmap["Creator"],&tr.Creator)
  json.Unmarshal(objmap["Hash"],&tr.Hash)
  json.Unmarshal(objmap["Signature"],&tr.Signature)
  tr.spent = false
  return tr
}

func (self *StdTransaction)GetType()string{
  return TrStd
}
