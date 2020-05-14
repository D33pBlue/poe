/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-23
 * @Project: Proof of Evolution
 * @Filename: std_trans.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-14
 * @Copyright: 2020
 */

package blockchain

import(
  "fmt"
  "time"
  // "errors"
  "io/ioutil"
  "encoding/hex"
  "encoding/json"
  "github.com/D33pBlue/poe/utils"
)

type JobTransaction struct{
  Timestamp time.Time
  BlockStart int
  BlockEnd int
  Inputs []TrInput
  Output TrOutput // the possible remainder
  Creator utils.Addr
  Job string
  Data string // nil if data from url
  DataUrl string // empty if data embedded
  Prize int
  Hash string
  Signature string
  spent bool // A job transaction is "spent" when closed
  fetched string
}

// Builds a new JobTransaction and signs it. This method does not check
// the data it receives in input.
func MakeJobTransaction(creator utils.Addr,key utils.Key,
      inps []TrInput,out TrOutput,
      job,data string,dataurl string,
      prize,bkStart,bkEnd int)*JobTransaction{
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
  tr.Output = out
  tr.Hash = tr.GetHash()
  tr.Signature = fmt.Sprintf("%x",utils.GetSignatureFromHash(tr.Hash,key))
  tr.spent = false
 return tr
}

func GetJobFixedCost(job,data string,url bool)int{
  return 0 // TODO: implement later
}

func GetJobMinPrize(job,data string)int{
  return 0 // TODO: implement later
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
// - the fixed cost calculated with GetJobFixedCost paid
// - a prize >= GetJobMinPrize paid
// - the BlockStart, BlockEnd slot that matches the one obtained with
// block.NextSlotForJobExectution on the previous block
func (self *JobTransaction)Check(block *Block,trChanges *map[string]string)bool{
  hash2 := self.GetHash()
  if hash2!=self.Hash{
    fmt.Println(self.Output)
    fmt.Println("The hash does not match")
    fmt.Printf("%v !=\n%v\n",hash2,self.Hash)
    return false}
  if !utils.CheckSignature(self.Signature,self.Hash,self.Creator){
    fmt.Println("The signature is not valid")
    return false}
  var tot int = 0
  for i:=0;i<len(self.Inputs);i++{
    source := self.Inputs[i].GetSource(block)
    if source==nil{
      // the Input source does not exist in the block's chain
      fmt.Println("Inesistend input source in transaction")
      return false
    }
    if source.Address!=self.Creator{
      // The Input source does not belong to the creator of this transaction
      fmt.Println("The input does not belong to the creator of the transaction")
      return false
    }
    spentInBlock := source.GetSpentIn()
    trSourceId := self.Inputs[i].ToString()
    if spentInBlock!=""{
      if block.FindPrevBlock(spentInBlock)!=nil{
        // double spending: the TrOutput was spent in a block
        // that is reachable from the current block
        fmt.Println("Double spending case 1")
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
          fmt.Println("Double spending case 2")
          return false
        }
      }
    }
    // the TrOutput was available and it is spent now
    // => update trChanges
    (*trChanges)[trSourceId] = block.GetHashCached()
    tot += source.Value
  }
  var datafromUrl bool = false
  if len(self.Data)<=0 && len(self.DataUrl)>0{
    datafromUrl = true
    self.fetched = fmt.Sprintf("%x",utils.FetchDataFromUrl(self.DataUrl))
  }else{
    self.fetched = self.Data
  }
  fixedCost := GetJobFixedCost(self.Job,self.fetched,datafromUrl)
  if self.Prize<GetJobMinPrize(self.Job,self.fetched){
    fmt.Println("The prize is too small")
    return false
  }
  var spent int = self.Output.Value+self.Prize+fixedCost
  if tot<spent{
    fmt.Printf("Tot: %v, spent: %v\n",tot,spent)
    return false
  }
  bk1,bk2 := block.Previous.NextSlotForJobExectution()
  if bk1!=self.BlockStart || bk2!=self.BlockEnd{
    fmt.Println("Invalid execution slot")
    return false
  }
  return true
}

// Stores the code of the Job in the file given in path.
func (self *JobTransaction)SaveJobInFile(path string)error{
  data,err := hex.DecodeString(self.Job)
  if err!=nil{
    return err
  }
  return ioutil.WriteFile(path, data, 0644)
}

// Stores the data of the job in the file given in path. If the
// data is embedded, it is saved directly. Otherwise, the data
// is downloaded from the given url.
func (self *JobTransaction)SaveDataInFile(path string)error{
  var data []byte
  if len(self.DataUrl)<=0{
    var err error
    data,err = hex.DecodeString(self.Data)
    if err!=nil{ return err }
  }else{
    //
    //
    // TODO: download the data from url
    //
    //
  }
  return ioutil.WriteFile(path, data, 0644)
}

// Returns always the only TrOutput stored inside the transaction.
func (self* JobTransaction)GetOutputAt(i int)*TrOutput{
  return &self.Output
}

// Returns the timestamp stored inside the transaction.
func (self *JobTransaction)GetTimestamp()time.Time{
  return self.Timestamp
}

// Returns the public key of creator of the transaction
func (self *JobTransaction)GetCreator()utils.Addr{
  return self.Creator
}

func (self *JobTransaction)GetPeriod()(int,int){
  return self.BlockStart,self.BlockEnd
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
  hb.Add(self.Output)
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
  json.Unmarshal(objmap["Output"],&tr.Output)
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
