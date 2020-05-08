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

// StdTransaction is the implementation of the standard Transaction,
// that can be used to send money.
// It collects a set of TrInput (source) and TrOutput (destination).
// All TrInput must point to TrOutput that belongs to the Creator
// of the transaction.
type StdTransaction struct{
  Timestamp time.Time
  Inputs []TrInput
  Outputs []TrOutput
  Creator utils.Addr
  Hash string
  Signature string
}

// MakeStdTransaction build and initialize a StdTransaction.
// This method does not check the validity of the parameters: they
// are assumed to be valid.
func MakeStdTransaction(creator utils.Addr,key utils.Key,
      inps []TrInput,outs []TrOutput)*StdTransaction{
  tr := new(StdTransaction)
  tr.Timestamp = time.Now()
  tr.Creator = creator
  tr.Inputs = inps
  tr.Outputs = outs
  tr.Hash = tr.GetHash()
  tr.Signature = fmt.Sprintf("%x",utils.GetSignatureFromHash(tr.Hash,key))
 return tr
}

// Check validate the transaction and update trChanges. The parameter block
// is assumed to be the block in which this transaction is stored.
// The validity is checked in relation to the chain that is linked
// to that block. The subchain must be valid/already checked.
// In order to be valid, the transaction must have:
// - the hash that matches the one declared
// - the signature verified with the public key of the creator
// - all sources that belongs to the creator
// - all sources unspent (no double spending)
// - the total amount of money in sources >= the total amount of money spent.
func (self *StdTransaction)Check(block *Block,trChanges *map[string]string)bool{
  hash2 := self.GetHash()
  if hash2!=self.Hash{return false}
  if !utils.CheckSignature(self.Signature,self.Hash,self.Creator){
    return false
  }
  var tot int = 0
  for i:=0;i<len(self.Inputs);i++{
    source := self.Inputs[i].GetSource(block)
    if source==nil{
      // the Input source does not exist in the block's chain
      return false
    }
    if source.Address!=self.Creator{
      // The Input source does not belong to the creator of this transaction
      return false
    }
    spentInBlock := source.GetSpentIn()
    trSourceId := self.Inputs[i].ToString()
    if spentInBlock!=""{
      if block.FindPrevBlock(spentInBlock)!=nil{
        // double spending: the TrOutput was spent in a block
        // that is reachable from the current block
        return false
      }else{
        // if block.FindPrevBlock(spentInBlock)==nil but spentInBlock!="",
        // the TrOutput was spent in a previous blockchain, but now the
        // blockchain considered is different (maybe due to a fork), and
        // in this new blockchain the TrOutput may be unspent. However,
        // it can be spent in the new chain, thus, trChanges must be
        // checked.
        if _, ok := (*trChanges)[trSourceId]; ok {
          // the transaction is spent in a block of the new chain
          // => double spending
          return false
        }
      }
    }
    // the TrOutput was available and it is spent now
    // => update trChanges
    (*trChanges)[trSourceId] = block.GetHashCached()
    tot += source.Value
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

func (self* StdTransaction)GetOutputAt(i int)*TrOutput{
  return &self.Outputs[i]
}

func (self *StdTransaction)GetTimestamp()time.Time{
  return self.Timestamp
}

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
  return tr
}

func (self *StdTransaction)GetType()string{
  return TrStd
}
