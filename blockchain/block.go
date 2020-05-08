/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-18
 * @Project: Proof of Evolution
 * @Filename: block.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-08
 * @Copyright: 2020
 */

// Package block contains the definitions
// of blocks, transactions, nonce
package blockchain

import(
  "fmt"
  "time"
  "sync"
  "encoding/json"
  "errors"
  "github.com/D33pBlue/poe/utils"
)

type Block struct{
  Previous *Block
  LenSubChain int
  Transactions *Tree
  Timestamp time.Time
  NumJobs int
  Hardness int
  NonceNoJob NonceNoJob
  MiniBlocks []MiniBlock
  Hash string
  checked bool
  mined bool
  access_data sync.Mutex
}

// Build the first block of the chain.
// Miners should call this function no more than one time.
func BuildFirstBlock(id utils.Addr)*Block{
  block := new(Block)
  block.LenSubChain = 0
  block.Timestamp = time.Now()
  block.NumJobs = 0
  block.Hardness = 0
  block.NonceNoJob.Value = 0
  block.Transactions = BuildMerkleTree()
  transact,_ := MakeCoinTransaction(id,block.calculateMiningValue())
  block.Transactions.Add(transact)
  block.Hash = block.GetHash("")
  block.checked = false
  block.mined = true
  return block
}

// Make a block (except the first one).
func BuildBlock(id utils.Addr,prev *Block)*Block{
  block := new(Block)
  block.Previous = prev
  block.LenSubChain = prev.LenSubChain+1
  block.Timestamp = time.Now()
  block.NumJobs = block.calculateNumJobs()
  block.Hardness = block.calculateHardness()
  block.NonceNoJob.Value = 0
  block.Transactions = BuildMerkleTree()
  transact,_ := MakeCoinTransaction(id,block.calculateMiningValue())
  block.Transactions.Add(transact)
  block.Hash = block.GetHash("")
  block.checked = false
  block.mined = false
  return block
}


func (self *Block)Mine(keepmining *bool){
  self.mined = false
  if self.NumJobs==0{
    self.mineNoJob(keepmining)
  }else{
    self.mineWithJobs(keepmining)
  }
}

// Returns a block in the chain pointed by this block;
// if there is no match, returns nil.
func (self *Block)FindPrevBlock(hash string)*Block{
  block := self
  for{
    if block==nil{return nil}
    if block.GetHashCached()==hash{
      return block
    }
    block = block.Previous
  }
  return nil
}

func (self *Block)FindTransaction(hash string)Transaction{
  return nil // TODO: implement later.
}

func (self *Block)Serialize()[]byte{
  type Block2 struct{
    Previous string
    LenSubChain int
    Transactions *Tree
    Timestamp time.Time
    NumJobs int
    Hardness int
    NonceNoJob NonceNoJob
    MiniBlocks []MiniBlock
    Hash string
  }
  block := new(Block2)
  if self.Previous!=nil{
    block.Previous = self.Previous.Hash
  }
  block.LenSubChain = self.LenSubChain
  block.Transactions = self.Transactions
  block.Timestamp = self.Timestamp
  block.NumJobs = self.NumJobs
  block.Hardness = self.Hardness
  block.NonceNoJob = self.NonceNoJob
  block.MiniBlocks = self.MiniBlocks
  block.Hash = self.Hash
  data, err := json.Marshal(block)
  if err != nil {
    fmt.Println(err)
  }
  return data
}

// Returns a Block from json and the hash of the previous block
func MarshalBlock(data []byte)(*Block,string){
  var objmap map[string]json.RawMessage
  json.Unmarshal(data, &objmap)
  var block *Block = new(Block)
  block.checked = false
  block.mined = true
  json.Unmarshal(objmap["LenSubChain"],&block.LenSubChain)
  json.Unmarshal(objmap["Timestamp"],&block.Timestamp)
  json.Unmarshal(objmap["NumJobs"],&block.NumJobs)
  json.Unmarshal(objmap["Hardness"],&block.Hardness)
  json.Unmarshal(objmap["NonceNoJob"],&block.NonceNoJob)
  json.Unmarshal(objmap["MiniBlocks"],&block.MiniBlocks)
  json.Unmarshal(objmap["Hash"],&block.Hash)
  block.Transactions = MarshalMerkleTree(objmap["Transactions"])
  var prev string = ""
  json.Unmarshal(objmap["Previous"],&prev)
  return block,prev
}

func (self *Block)mineWithJobs(keepmining *bool){
  // TODO: implement later
}

func (self *Block)mineNoJob(keepmining *bool){
  for{
    if !(*keepmining) {break}
    if self.mined{
      break
    }
    self.access_data.Lock()
    self.NonceNoJob.Next()
    self.Hash = self.GetHash("")
    self.mined = self.checkNonceNoJob()
    self.access_data.Unlock()
  }
  fmt.Println("ckck",self.GetHash(""))
  fmt.Println("ckck2",self.GetHash(self.Previous.GetHashCached()))
}

func (self *Block)AddTransaction(transact *Transaction)error{
  self.access_data.Lock()
  if self.mined{
    self.access_data.Unlock()
    return errors.New("Tried to add transaction in block already mined")
  }
  // TODO: implement later
  self.access_data.Unlock()
  return nil
}

func (self *Block)AddMiniBlock(miniblock *MiniBlock)error{
  self.access_data.Lock()
  if self.mined{
    self.access_data.Unlock()
    return errors.New("Tried to add miniblock in block already mined")
  }
  // TODO: implement later
  self.access_data.Unlock()
  return nil
}

// The CheckStep1 method checks the validity of the content of
// the block without checking Previous and its connections.
func (self *Block)CheckStep1(hashPrev string)bool{
  if self.checked { return true }
  // if !utils.CompareHashes(self.Hash,self.GetHash(hashPrev)){
  if self.Hash!=self.GetHash(hashPrev){
    fmt.Println("error in hash")
    fmt.Println(self.Hash)
    fmt.Println(self.GetHash(hashPrev))
    return false
  }
  if self.NumJobs==0{
    if !self.checkNonceNoJob() {
      fmt.Println("error in checkNonceNoJob")
      return false }
  }else if self.NumJobs>0{
    if !self.checkNonceJobs() { return false }
  }else{ return false }
  return true
}

// The CheckStep2 method checks the validity of the links
// among blocks and of the depending data.
func (self *Block)CheckStep2(transactionChanges *map[string]string)bool{
  if self.checked { return true }
  if !self.checkNumJobs() {
    fmt.Println("fail in num jobs")
    return false }
  if !self.checkHardness() {
    fmt.Println("fail hardness")
    return false }
  if self.LenSubChain>0{
    if self.Previous==nil {
      fmt.Println("fail in previous")
      return false }
    if self.Previous.LenSubChain!=self.LenSubChain-1 {
      fmt.Println("fail in LenSubChain")
      return false }
    if !self.Previous.CheckStep2(transactionChanges){ return false } // TODO: remove recursion
  }
  // transactions are checked only if the previous blocks are valid
  if !self.checkTransactions(transactionChanges) {
    fmt.Println("fail in transactions")
    return false
  }
  self.checked = true
  return true
}

func (self *Block)GetHash(hashPrev string)string{
  hb := new(utils.HashBuilder)
  hb.Add(self.Timestamp.Format("2006-01-02 15:04:05"))
  hb.Add(self.LenSubChain)
  hb.Add(self.Transactions.GetHash())
  hb.Add(self.NumJobs)
  hb.Add(self.Hardness)
  if self.NumJobs==0{
    hb.Add(self.NonceNoJob.Value)
  }else{
    for i:=0;i<len(self.MiniBlocks);i++{
      hb.Add(self.MiniBlocks[i].GetHashCached())
    }
  }
  if self.LenSubChain>0{
    if len(hashPrev)>0{
      hb.Add(hashPrev)
    }else{
      hb.Add(self.Previous.GetHashCached())
    }
  }
  hash := hb.GetHash()
  return fmt.Sprintf("%x",hash)
}

func (self *Block)GetHashCached()string{
  return self.Hash
}

func (self *Block)checkNonceNoJob()bool{
  hash := self.Hash
  // fmt.Println(hash)
  for i:=0;i<self.Hardness;i++{
    if hash[i]!='0'{ return false }
  }
  return true
}

func (self *Block)checkNonceJobs()bool{
  return true // TODO: implement later
}

func (self *Block)checkNumJobs()bool{
  return true // TODO: implement later
}

func (self *Block)checkHardness()bool{
  return true // TODO: implement later
}

func (self *Block)checkTransactions(transactionChanges *map[string]string)bool{
  // check the Merkle tree hashes
  if !self.Transactions.Check(){ return false }
  // check all transactions in the tree
  transactions := self.Transactions.GetTransactionArray()
  for i:=0;i<len(transactions);i++{
    if !transactions[i].Check(self,transactionChanges){
      return false
    }
  }
  return true
}

func (self *Block)calculateNumJobs()int{
  return 0 // TODO: implement later
}

func (self *Block)calculateHardness()int{
  return 6 // TODO: implement later
}

// Returns the value in coin of mining that block.
// This value depends on the hardness of the mining task.
func (self *Block)calculateMiningValue()int{
  return (1+self.Hardness)*10 // TODO: tune with mining time
}
