/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-23
 * @Project: Proof of Evolution
 * @Filename: std_trans.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-10
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
  BlockStart int
  BlockEnd int
  Inputs []TrInput
  Creator utils.Addr
  Job string
  Data []byte // nil if data from url
  DataUrl string // empty if data embedded
  Prize int
  Hash string
  Signature string
  spent bool // A job transaction is "spent" when closed
}

// Builds a new JobTransaction and signs it. This method does not
// the data it receives in input.
func MakeJobTransaction(creator utils.Addr,key utils.Key,
      inps []TrInput,job,dataurl string,data []byte,
      prize,bkStart,bkEnd int)(*JobTransaction,error){
  tr := new(JobTransaction)
  tr.Timestamp = time.Now()
  tr.BlockStart = bkStart
  tr.BlockEnd = bkEnd
  tr.Job = job
  tr.Data = data
  tr.DataUrl = dataurl
  tr.Prize = prize
  tr.Creator = creator
  tr.Inputs = inps
  tr.Hash = tr.GetHash()
  tr.Signature = fmt.Sprintf("%x",utils.GetSignatureFromHash(tr.Hash,key))
  tr.spent = false
 return tr,nil
}

func (self *JobTransaction)Check(block *Block,trChanges *map[string]string)bool{
  return true // TODO: implement later
}

// Stores the code of the Job in the file given in path.
func (self *JobTransaction)SaveJobInFile(path string)error{
  return nil // TODO: implement later
}

// Stores the data of the job in the file given in path. If the
// data is embedded, it is saved directly. Otherwise, the data
// is downloaded from the given url.
func (self *JobTransaction)SaveDataInFile(path string)error{
  return nil // TODO: implement later
}

// Returns always nil, because a job transaction has no output.
func (self* JobTransaction)GetOutputAt(i int)*TrOutput{
  return nil
}

// Returns the timestamp stored inside the transaction.
func (self *JobTransaction)GetTimestamp()time.Time{
  return self.Timestamp
}

// Returns the public key of creator of the transaction
func (self *JobTransaction)GetCreator()utils.Addr{
  return self.Creator
}

// Recalculates the hash of the transaction.
func (self *JobTransaction)GetHash()string{
  hb := new(utils.HashBuilder)
  hb.Add(self.Creator)
  hb.Add(self.Timestamp.Format("2006-01-02 15:04:05"))
  hb.Add(self.BlockStart)
  hb.Add(self.BlockEnd)
  for i:=0;i<len(self.Inputs);i++{
    hb.Add(self.Inputs[i].Block)
    hb.Add(self.Inputs[i].ToSpend)
    hb.Add(self.Inputs[i].Index)
  }
  hb.Add(self.Job)
  hb.Add(self.Data)
  hb.Add(self.DataUrl)
  hb.Add(self.Prize)
  return fmt.Sprintf("%x",hb.GetHash())
}

// Returns the cached hash of the transaction.
func (self *JobTransaction)GetHashCached()string{
  return self.Hash
}

// Serializes the transaction and returns it as []byte.
func (self *JobTransaction)Serialize()[]byte{
  data, err := json.Marshal(self)
  if err != nil {
    fmt.Println(err)
  }
  return data
}

// Rebuilds the transaction from its serialized data.
func MarshalJobTransaction(data []byte)*JobTransaction{
  var objmap map[string]json.RawMessage
  json.Unmarshal(data, &objmap)
  tr := new(JobTransaction)
  json.Unmarshal(objmap["Timestamp"],&tr.Timestamp)
  json.Unmarshal(objmap["BlockStart"],&tr.BlockStart)
  json.Unmarshal(objmap["BlockEnd"],&tr.BlockEnd)
  json.Unmarshal(objmap["Inputs"],&tr.Inputs)
  json.Unmarshal(objmap["Job"],&tr.Job)
  json.Unmarshal(objmap["Data"],&tr.Data)
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
