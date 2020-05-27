/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-18
 * @Project: Proof of Evolution
 * @Filename: block.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-27
 * @Copyright: 2020
 */

// Package block contains the definitions
// of blocks, transactions, nonce
package blockchain

import(
  "fmt"
  "time"
  "sync"
  "sort"
  "encoding/json"
  "errors"
  "github.com/D33pBlue/poe/utils"
  "github.com/D33pBlue/poe/ga"
  "github.com/D33pBlue/poe/conf"
)

// The block struct models the blocks of the blockchain. Each block,
// except the first one, is linked to a previous one.
type Block struct{
  Previous *Block
  LenSubChain int
  Transactions *Tree
  MerkleHash string
  Timestamp time.Time
  NumJobs int
  Hardness int
  NonceNoJob NonceNoJob
  MiniBlocks []MiniBlock
  Hash string
  checked bool
  mined bool
  access_data sync.Mutex
  jobs []*JobTransaction // cached list of jobs not ended
  incomingMiniblock chan *MiniBlock // to receive a MiniBlock from others
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
  block.MerkleHash = block.Transactions.GetHash()
  transact,_ := MakeCoinTransaction(id,block.calculateMiningValue())
  block.Transactions.Add(transact)
  block.Hash = block.GetHash("")
  block.checked = false
  block.mined = true
  return block
}

// Make a new block (except the first one). This block needs
// to be mined. While the block is not mined, it can accept transactions.
func BuildBlock(id utils.Addr,prev *Block)*Block{
  block := new(Block)
  block.Previous = prev
  block.LenSubChain = prev.LenSubChain+1
  block.Timestamp = time.Now()
  block.NumJobs = block.calculateNumJobs()
  block.Hardness = block.calculateHardness()
  block.NonceNoJob.Value = 0
  block.Transactions = BuildMerkleTree()
  // with no jobs a CoinTransaction is added for the miner who mines the block;
  // with jobs, CoinTransactions are added later with MiniBlocks
  if block.NumJobs<=0{
    transact,_ := MakeCoinTransaction(id,block.calculateMiningValue())
    block.Transactions.Add(transact)
  }
  block.Hash = block.GetHash("")
  block.checked = false
  block.mined = false
  block.incomingMiniblock = make(chan *MiniBlock,10)
  block.assignPrizes()
  return block
}

// The function to call in order to start mining a block.
// It should be called in a goroutine. If there are jobs,
// this function calls mineWithJobs; otherwise mineNoJob.
// Miniblocks are sent to miniblockout whenever mined.
func (self *Block)Mine(id utils.Addr,keepmining *bool,
      miniblockout chan MexBlock,executor *ga.Executor,config *conf.Config){
  self.mined = false
  if self.NumJobs==0{
    self.mineNoJob(keepmining)
  }else{
    self.mineWithJobs(id,keepmining,miniblockout,executor,config)
  }
}

// Returns a block in the chain pointed by this block,
// searching with an hash as key; if there is no match, returns nil.
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

// TODO: improve efficiency
// Returns (if exists) the transaction in this block with
// a specific hash, or nil
func (self *Block)FindTransaction(hash string)Transaction{
  transacts := self.Transactions.GetTransactionArray()
  for i:=0;i<len(transacts);i++{
    if transacts[i].GetHashCached()==hash{
      return transacts[i]
    }
  }
  return nil
}

// Serializes the block returning a []byte
func (self *Block)Serialize()[]byte{
  type Block2 struct{
    Previous string
    LenSubChain int
    Transactions *Tree
    MerkleHash string
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
  block.MerkleHash = self.MerkleHash
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

// Returns a Block parsing the serialized json, and also
// the hash of the previous block. The block is not linked to
// the previous one.
func MarshalBlock(data []byte,config *conf.Config)(*Block,string){
  var objmap map[string]json.RawMessage
  json.Unmarshal(data, &objmap)
  var block *Block = new(Block)
  block.checked = false
  block.mined = true
  json.Unmarshal(objmap["LenSubChain"],&block.LenSubChain)
  json.Unmarshal(objmap["Timestamp"],&block.Timestamp)
  json.Unmarshal(objmap["NumJobs"],&block.NumJobs)
  json.Unmarshal(objmap["Hardness"],&block.Hardness)
  json.Unmarshal(objmap["MerkleHash"],&block.MerkleHash)
  json.Unmarshal(objmap["NonceNoJob"],&block.NonceNoJob)
  json.Unmarshal(objmap["MiniBlocks"],&block.MiniBlocks)
  json.Unmarshal(objmap["Hash"],&block.Hash)
  block.Transactions = MarshalMerkleTree(objmap["Transactions"],config)
  var prev string = ""
  json.Unmarshal(objmap["Previous"],&prev)
  return block,prev
}

// The mining process with jobs. Istantiate the right number of miniblocks
// and call their mining method.
func (self *Block)mineWithJobs(id utils.Addr,keepmining *bool,
          miniblockout chan MexBlock,executor *ga.Executor,config *conf.Config){
  MBkeepaliveIndex := make(map[string]int)
  var MBkeepalive []bool
  MBminedChan := make(chan *MiniBlock)
  jobtransactions := self.getJobsForThisBlock()
  remaining := 0
  for i:=0;i<len(jobtransactions);i++{
    transact := jobtransactions[i]
    hash := transact.GetHashCached()
    // create or retrieve job channels through executor
    var jobchs *ga.JobChannels
    if executor.IsExecutingJob(hash){
      jobchs = executor.GetChannels(hash)
      executor.ChangeBlockHashInJob(hash,
        self.Previous.GetHashCached(),string(id))
    }else{
      jobPath,dataPath := config.GetSuitablePathForJob(hash)
      err := transact.SaveJobInFile(jobPath)
      err2 := transact.SaveDataInFile(dataPath)
      if err==nil && err2==nil{
        jobchs = executor.StartJob(hash,
          self.Previous.GetHashCached(),string(id),jobPath,dataPath)
      }else{
        fmt.Println(err,err2)
      }
    }
    if jobchs!=nil{
      // create a keepmining variable for each miniblock
      MBkeepaliveIndex[hash] = len(MBkeepalive)
      MBkeepalive = append(MBkeepalive,true)
      remaining += 1
      // initialize the miniblocks and start mining them
      miniblock := BuildMiniBlock(self.Previous.GetHashCached(),
        transact.blockContainer,hash,id,jobchs.ChNonce)
      go miniblock.Mine(self.Hardness,&MBkeepalive[MBkeepaliveIndex[hash]],MBminedChan)
    }
  }
  // listen for mined miniblock in return through MBminedChan
  // and propagate when received through miniblockout (except the last one),
  // listen also for miniblock from other miners through incomingMiniblock
  // and eventually stop current miniblock's mining process.
  for {
    // remaining is decreased each time a MiniBlock is mined by someone
    if remaining<=0{
      break
    }
    // fmt.Println(*keepmining)
    if !(*keepmining){
      // stop all miniblock's mining processes
      fmt.Println("Stopping all MBkeepalive")
      for i:=0;i<len(MBkeepalive);i++{
        MBkeepalive[i] = false
      }
    }
    select{
    case mb := <- MBminedChan:
      // insert miniblock and propagate through miniblockout
      if mb!=nil{
        if MBkeepalive[MBkeepaliveIndex[mb.JobTrans]]{
          MBkeepalive[MBkeepaliveIndex[mb.JobTrans]] = false
          self.storeMiniblockInBlock(mb)
          if remaining>1{
            // send the mined MiniBlocks to the others, but
            // do not propagate the last MiniBlock: it is
            // propagated with the full Block
            mex := new(MexBlock)
            mex.Data = mb.Serialize()
            mex.IpSender = string(id)
            miniblockout <- (*mex)
          }
        }
      }
      if remaining<=1{
        // last MiniBlock mined
        self.access_data.Lock()
        self.MerkleHash = self.Transactions.GetHash()
        self.Hash = self.GetHash("")
        self.mined = true
        self.access_data.Unlock()
      }
      remaining -= 1
    case mb := <- self.incomingMiniblock:
      // insert miniblock and stop its mining
      if index, ok := MBkeepaliveIndex[mb.JobTrans]; ok {
        if MBkeepalive[index]{
          MBkeepalive[index] = false
          fmt.Println("Storing incoming miniblock")
          self.storeMiniblockInBlock(mb)
        }
      }
    default:
    }
  }
  fmt.Println("End/stop mining (miniblocks)")
  //
  // where to stop jobs from executor when the slot expires??
}

// Adds all prizes for the solutions of the previous block
func (self *Block)assignPrizes(){
  if self.Previous == nil{ return }
  transactions := self.Previous.Transactions.GetTransactionArray()
  solutions := make(map[string]([]*SolTransaction))
  // collect all the solutions submitted in the previous block
  for i:=0;i<len(transactions);i++{
    if transactions[i].GetType()==TrSol{
      tr := transactions[i].(*SolTransaction)
      solutions[tr.JobTrans] = append(solutions[tr.JobTrans],tr)
    }
  }
  // sort the solutions by the fitness score for each job
  for k,_ := range solutions{
    sort.SliceStable(solutions[k], func(i, j int) bool {
      res_i := self.Previous.Previous.FindTransaction(solutions[k][i].ResTrans).(*ResTransaction)
      res_j := self.Previous.Previous.FindTransaction(solutions[k][j].ResTrans).(*ResTransaction)
      isMin := res_i.IsMin
      if isMin{
        return res_i.Evaluation < res_j.Evaluation
      }
      return  res_i.Evaluation > res_j.Evaluation })
    resTr := self.Previous.Previous.FindTransaction(solutions[k][0].ResTrans).(*ResTransaction)
    jobBlock := self.FindPrevBlock(resTr.JobBlock)
    jobTr := jobBlock.FindTransaction(resTr.JobTrans).(*JobTransaction)
    for i:=0;i<len(solutions[k]);i++{
      prize := jobTr.GetSharingBlockPrize(i+1)
      if prize<=0{
        break
      }else{
        // TODO: make PrizeTransaction
      }
    }
  }
}

// The mining process without jobs => PoW.
func (self *Block)mineNoJob(keepmining *bool){
  for{
    if !(*keepmining) {break}
    if self.mined{
      break
    }
    self.access_data.Lock()
    self.NonceNoJob.Next()
    self.MerkleHash = self.Transactions.GetHash()
    self.Hash = self.GetHash("")
    self.mined = self.checkNonceNoJob()
    self.access_data.Unlock()
  }
  fmt.Println("ckck",self.GetHash(""))
  fmt.Println("ckck2",self.GetHash(self.Previous.GetHashCached()))
}

// Inserts a transaction in the block without checking it.
// It is assumed that the transaction is valid.
// Returns an error if the block is already mined.
func (self *Block)AddTransaction(transact Transaction)error{
  self.access_data.Lock()
  if self.mined{
    self.access_data.Unlock()
    return errors.New("Tried to add transaction in block already mined")
  }
  self.Transactions.Add(transact)
  self.access_data.Unlock()
  return nil
}


func (self *Block)AddMiniBlock(miniblock *MiniBlock,config *conf.Config){
  // check MiniBlock
  if !miniblock.CheckStep1(self.Previous.Hash){
    fmt.Println("Miniblock does not pass CheckStep1")
    return
  }
  if !miniblock.CheckStep2(self,self.Hardness,config){
    fmt.Printf("Error in miniblock in block %v\n",self.GetBlockIndex())
    return
  }
  self.incomingMiniblock <- miniblock
  fmt.Println("Incoming miniblock processed.")
}

func (self *Block)storeMiniblockInBlock(miniblock *MiniBlock){
  self.access_data.Lock()
  if self.mined{
    self.access_data.Unlock()
    return
  }
  self.MiniBlocks = append(self.MiniBlocks,*miniblock)
  transact,_ := MakeCoinTransaction(miniblock.Miner,self.calculateMiningValue())
  fmt.Println("CoinTransaction inserted for MiniBlock")
  // fmt.Printf("Before: %v %v\n",self.Transactions.CountCoinTr2(),CountCoinTr(self.Transactions.Root))
  self.Transactions.Add(transact)
  // fmt.Printf("After: %v %v\n",self.Transactions.CountCoinTr2(),CountCoinTr(self.Transactions.Root))
  self.access_data.Unlock()
}

// The CheckStep1 method checks the validity of the content of
// the block without checking Previous and its connections.
func (self *Block)CheckStep1(hashPrev string)bool{
  if self.checked { return true }
  // if !utils.CompareHashes(self.Hash,self.GetHash(hashPrev)){
  if self.Hash!=self.GetHash(hashPrev){
    fmt.Printf("tree %v\n",self.Transactions.GetHash())
    fmt.Printf("cach %v\n",self.MerkleHash)
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
    if !self.checkNonceJobsStep1(hashPrev) {
      fmt.Println("error in checkNonceJobsStep1")
      return false
    }
  }else{ return false }
  return true
}

// The CheckStep2 method checks the validity of the links
// among blocks and of the depending data.
func (self *Block)CheckStep2(transactionChanges *map[string]string,config *conf.Config)bool{
  if self.checked { return true }
  if self.NumJobs>0{
    if !self.checkNonceJobsStep2(config){
      fmt.Println("fail in checkNonceJobsStep2")
      return false
    }
  }
  if !self.checkNumJobs() {
    fmt.Println("fail in num jobs")
    return false }
  if !self.checkHardness() {
    fmt.Println("fail hardness")
    fmt.Println(self.GetBlockIndex(),self.Hardness,self.calculateHardness())
    return false }
  if self.LenSubChain>0{
    if self.Previous==nil {
      fmt.Println("fail in previous")
      return false }
    if self.Previous.LenSubChain!=self.LenSubChain-1 {
      fmt.Println("fail in LenSubChain")
      return false }
    if !self.Previous.CheckStep2(transactionChanges,config){ return false } // TODO: remove recursion
  }
  // transactions are checked only if the previous blocks are valid
  if !self.checkTransactions(transactionChanges) {
    fmt.Println("fail in transactions")
    return false
  }
  //
  //
  // TODO: add check of the prizes!!
  //
  //
  self.checked = true
  return true
}

// Recalculates the hash of the block as hex string.
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

// Returns the cached hash of the block.
func (self *Block)GetHashCached()string{
  return self.Hash
}

// Returns the index of the current block in the entire Blockchain.
func (self *Block)GetBlockIndex()int{
  return self.LenSubChain
}

// Checks if the nonce has the required number of 0 bits.
func (self *Block)checkNonceNoJob()bool{
  hash := self.Hash
  for i:=0;i<self.Hardness;i++{
    if hash[i]!='0'{ return false }
  }
  return true
}

// Checks the validity of the hash of each miniblock (without
// checking their job's solution score).
func (self *Block)checkNonceJobsStep1(hashPrev string)bool{
  for i:=0;i<len(self.MiniBlocks);i++{
    if !self.MiniBlocks[i].CheckStep1(hashPrev){
      return false
    }
  }
  return true
}

// Complete the check of the validity of each MiniBlock
// evaluating the given solutions.
func (self *Block)checkNonceJobsStep2(config *conf.Config)bool{
  for i:=0;i<len(self.MiniBlocks);i++{
    if !self.MiniBlocks[i].CheckStep2(self,self.Hardness,config){
      fmt.Printf("Error in miniblock in block %v\n",self.GetBlockIndex())
      return false
    }
  }
  return true
}

// Checks if the number of miniblocks (job executed)
// is correct
func (self *Block)checkNumJobs()bool{
  return self.calculateNumJobs()==self.NumJobs
}

// Checks if the Hardness in the block matches the real one,
// recalculating it.
func (self *Block)checkHardness()bool{
  return self.calculateHardness()==self.Hardness
}

// Checks the consistency of the merkle tree and the transactions.
func (self *Block)checkTransactions(transactionChanges *map[string]string)bool{
  // check the Merkle tree hashes
  if !self.Transactions.Check(){
    fmt.Println("Invalid merkle tree")
    return false
  }
  // check all transactions in the tree
  transactions := self.Transactions.GetTransactionArray()
  minersCoin := make(map[string]int)
  for i:=0;i<len(transactions);i++{
    if !transactions[i].Check(self,transactionChanges){
      return false
    }
    // update minersCoin with the number of CoinTransaction
    if transactions[i].GetType()==TrCoin{
      var miner string = string(transactions[i].GetCreator())
      if _,ok := minersCoin[miner];ok{
        minersCoin[miner] += 1
      }else{
        minersCoin[miner] = 1
      }
    }
  }
  // check the CoinTransactions to pay who mined the block/miniblocks
  if self.NumJobs>0{
    shouldBePaid := make(map[string]int)
    for i:=0;i<self.NumJobs;i++{
      var miner string = string(self.MiniBlocks[i].Miner)
      if _,ok := shouldBePaid[miner];ok{
        shouldBePaid[miner] += 1
      }else{
        shouldBePaid[miner] = 1
      }
    }
    for k,v := range shouldBePaid{
      if minersCoin[k] != v {
        fmt.Println("A miniblock's miner has not been paid correctly in block",self.GetBlockIndex())
        fmt.Println(minersCoin[k],v)
        return false
      }
    }
  }else{
    if len(minersCoin)!=1{
      fmt.Println("Wrong number of CoinTransaction",len(minersCoin))
      return false
    }
  }
  return true
}

// Returns the list of not ended jobs (JobTransaction)
// that are stored in the chain starting from the previous block.
func (self *Block)getOpenJobs()[]*JobTransaction{
  if len(self.jobs)>0{
    return self.jobs // cached information
  }
  if self.Previous == nil{ return nil }
  // get the jobs from Previous.Previous
  previous := self.Previous.getOpenJobs()
  // add the ones of Previous
  transacts := self.Previous.Transactions.GetTransactionArray()
  hashcontainer := self.Previous.GetHashCached()
  for i:=0;i<len(transacts);i++{
    if transacts[i].GetType()==TrJob{
      jobtr := transacts[i].(*JobTransaction)
      jobtr.blockContainer = hashcontainer
      previous = append(previous,jobtr)
    }
  }
  // add the not ended ones to self.jobs
  for i:=0;i<len(previous);i++{
    _,end := previous[i].GetPeriod()
    if end>=self.GetBlockIndex(){
      self.jobs = append(self.jobs,previous[i])
    }
  }
  return self.jobs
}

func (self *Block)getJobsForThisBlock()[]*JobTransaction{
  jobtrans := self.getOpenJobs()
  var jobs []*JobTransaction
  for i:=0;i<len(jobtrans);i++{
    start,end := jobtrans[i].GetPeriod()
    index := self.GetBlockIndex()
    if index>=start && index<=end{
      jobs = append(jobs,jobtrans[i])
    }
  }
  return jobs
}

// Returns the number of Jobs that has to be executed to
// mine this block.
func (self *Block)calculateNumJobs()int{
  return len(self.getJobsForThisBlock())
}

// Returns the target hardness of this block (number of 0 chars
// at the beginning of the exadecimal hash of the block).
// This hardness is the one used in PoW, and has to be adjusted in
// MiniBlocks in according with the complexity of the job in use.
func (self *Block)calculateHardness()int{
  if self.GetBlockIndex()==0{
    return 0
  }
  // TODO: tune with mining time
  return 6
}

// Returns the slot in which the next job should execute.
// The slot is made of the indexes of the block in which the
// job should start, and the one in which it should end.
func (self *Block)NextSlotForJobExectution()(int,int){
  index := self.GetBlockIndex()
  return index+1,index+5 // TODO: tune with complexity and number of open jobs
}

// Returns the value in coin of mining this block/miniblock.
// This value depends on the hardness of the mining task.
func (self *Block)calculateMiningValue()int{
  return (1+self.Hardness)*10 // TODO: tune with mining time
}
