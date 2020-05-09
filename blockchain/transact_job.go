/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-23
 * @Project: Proof of Evolution
 * @Filename: std_trans.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-09
 * @Copyright: 2020
 */

package blockchain

import(
  "fmt"
  "time"
  // "errors"
  "encoding/json"
  "github.com/D33pBlue/poe/utils"
)

type JobTransaction struct{
  Timestamp time.Time
  Inputs []TrInput
  Creator utils.Addr
  Job string
  Prize int
  Hash string
  Signature string
  spent bool
}

// func MakeJobTransaction(creator utils.Addr,key utils.Key,
//       inps []TrInput,outs []TrOutput)(*JobTransaction,error){
//   tr := new(JobTransaction)
//   tr.Timestamp = time.Now()
//   tr.Creator = creator
//   tr.Inputs = inps
//   tr.Outputs = outs
//   tr.Hash = tr.GetHash()
//   tr.Signature = utils.GetSignatureFromHash(tr.Hash,key)
//  return tr,nil
// }

func (self *JobTransaction)Check(block *Block,trChanges *map[string]string)bool{
  return true // TODO: implement later
}

func (self* JobTransaction)GetOutputAt(i int)*TrOutput{
  return nil
}

func (self *JobTransaction)GetTimestamp()time.Time{
  return self.Timestamp
}

func (self *JobTransaction)GetCreator()utils.Addr{
  return self.Creator
}

func (self *JobTransaction)GetHash()string{
  hb := new(utils.HashBuilder)
  hb.Add(self.Creator)
  hb.Add(self.Timestamp)
  for i:=0;i<len(self.Inputs);i++{
    hb.Add(self.Inputs[i].Block)
    hb.Add(self.Inputs[i].ToSpend)
    hb.Add(self.Inputs[i].Index)
  }
  hb.Add(self.Job)
  hb.Add(self.Prize)
  return fmt.Sprintf("%x",hb.GetHash())
}

func (self *JobTransaction)GetHashCached()string{
  return self.Hash
}

func (self *JobTransaction)Serialize()[]byte{
  data, err := json.Marshal(self)
  if err != nil {
    fmt.Println(err)
  }
  return data
}

func MarshalJobTransaction(data []byte)*JobTransaction{
  var objmap map[string]json.RawMessage
  json.Unmarshal(data, &objmap)
  tr := new(JobTransaction)
  json.Unmarshal(objmap["Timestamp"],&tr.Timestamp)
  json.Unmarshal(objmap["Inputs"],&tr.Inputs)
  json.Unmarshal(objmap["Job"],&tr.Job)
  json.Unmarshal(objmap["Prize"],&tr.Prize)
  json.Unmarshal(objmap["Creator"],&tr.Creator)
  json.Unmarshal(objmap["Hash"],&tr.Hash)
  json.Unmarshal(objmap["Signature"],&tr.Signature)
  tr.spent = false
  return tr
}

func (self *JobTransaction)GetType()string{
  return TrJob
}
