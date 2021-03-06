/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-23
 * @Project: Proof of Evolution
 * @Filename: blockchain.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-28
 * @Copyright: 2020
 */


package blockchain

import (
  "fmt"
  "net"
  "sync"
  "bufio"
  "strconv"
  "strings"
  "errors"
  "io/ioutil"
  "github.com/D33pBlue/poe/utils"
  "github.com/D33pBlue/poe/ga"
  "github.com/D33pBlue/poe/conf"
)

// match any type
type any interface{}

// used to exchange transactions with miner
type MexTrans struct{
  Type string
  Data []byte
}

// used to exchange blocks with miner
type MexBlock struct{
  Data []byte
  IpSender string
}

// The blockchain has
// - a reference to the head block,
// - a current block that has to be mined before being inserted
// in the blockchain and being propagated,
// - some channels to communicate with miner and pass data
// - a publicKey (id) to get money from mining,
// - a lock for sync and some boolean to manage the mining process.
// (The miner has not to give the total money of clients: it is a problem
// of the client to imoplement this, by collecting old blocks/transactions).
type Blockchain struct{
  Head *Block
  Current *Block
  TransQueue chan MexTrans// receive transactions from miner
  BlockOut chan MexBlock// send mined block to miner
  BlockIn chan MexBlock// receive mined block from miner
  MiniBlockOut chan MexBlock// send mined miniblock to miner
  MiniBlockIn chan MexBlock// receive mined miniblock from miner
  internalBlock chan MexBlock// notify mined block to Communicate
  internalMiniBlock chan MexBlock// notify mined miniblock to Communicate (don't send the last one)
  id utils.Addr
  keepmining bool
  miningstatus chan bool
  access_data sync.Mutex
  folder string
  currentTrChanges map[string]string
  executor *ga.Executor
  config *conf.Config
  solutionToShare map[string]([]byte) // map the hash of the solution to the solution in []byte.
  chGoodSolutions chan ga.Sol // channel to receive good solutions from Executor.
}

// This method load a serialized blockchain from file.
// (Each time a block is added to the blockchain, in fact, it is
// saved in a file with its index as name).
func LoadChainFromFolder(id utils.Addr,folder string,config *conf.Config)*Blockchain{
  chain := new(Blockchain)
  chain.TransQueue = make(chan MexTrans)
  chain.BlockOut = make(chan MexBlock)
  chain.BlockIn = make(chan MexBlock)
  chain.MiniBlockOut = make(chan MexBlock)
  chain.MiniBlockIn = make(chan MexBlock)
  chain.internalBlock = make(chan MexBlock)
  chain.internalMiniBlock = make(chan MexBlock)
  chain.miningstatus = make(chan bool)
  chain.folder = folder
  chain.currentTrChanges = make(map[string]string)
  chain.Head = nil
  chain.solutionToShare = make(map[string]([]byte))
  chain.chGoodSolutions = make(chan ga.Sol,1000)
  chain.executor = ga.BuildExecutor(chain.chGoodSolutions)
  chain.config = config
  var i int = 0
  fmt.Println("Loading chain from disk")
  for{
    var filename string = fmt.Sprintf("%v/block%v.json",chain.folder,i)
    if !utils.FileExists(filename){break}
    data,err := ioutil.ReadFile(filename)
    if err!=nil{
      fmt.Println(err)
      return nil
    }
    b,hash := MarshalBlock(data,chain.config)
    if chain.Head!=nil && chain.Head.Hash!=hash{
      fmt.Println("Chain head not match")
      fmt.Println(chain.Head.Hash)
      fmt.Println(hash)
      return nil
    }
    b.Previous = chain.Head
    if !b.CheckStep1(hash){
      fmt.Println("Fail in CheckStep1")
      return nil
    }
    chain.Head = b
    fmt.Printf("Loaded block %v\n",filename)
    i += 1
  }
  fmt.Println("Checking loaded chain")
  trChanges := make(map[string]string)
  if !chain.Head.CheckStep2(&trChanges,chain.config){ return nil }
  fmt.Println("Initializing chain")
  chain.applyTransactionsChanges(&trChanges)
  chain.Current = BuildBlock(id,chain.Head)
  chain.id = id
  chain.keepmining = false
  fmt.Println("The chain is ready")
  go chain.Mine()
  return chain
}

// Create a new Blockchain (the first block) and initialize it.
// This method also start mining by calling Mine() in a goroutine.
func NewBlockchain(id utils.Addr,folder string,config *conf.Config)*Blockchain{
  chain := new(Blockchain)
  chain.TransQueue = make(chan MexTrans)
  chain.BlockOut = make(chan MexBlock)
  chain.BlockIn = make(chan MexBlock)
  chain.MiniBlockOut = make(chan MexBlock)
  chain.MiniBlockIn = make(chan MexBlock)
  chain.internalBlock = make(chan MexBlock)
  chain.internalMiniBlock = make(chan MexBlock)
  chain.miningstatus = make(chan bool)
  chain.Head = BuildFirstBlock(id)
  chain.solutionToShare = make(map[string]([]byte))
  chain.chGoodSolutions = make(chan ga.Sol,1000)
  chain.executor = ga.BuildExecutor(chain.chGoodSolutions)
  chain.config = config
  chain.folder = folder
  chain.currentTrChanges = make(map[string]string)
  var filename string = fmt.Sprintf("%v/block%v.json",chain.folder,chain.Head.LenSubChain)
  data := chain.Head.Serialize()
  err := ioutil.WriteFile(filename,data, 0644)
  if err!=nil{fmt.Println(err)}
  chain.Current = BuildBlock(id,chain.Head)
  chain.id = id
  chain.keepmining = false
  go chain.Mine()
  return chain
}

// Returns a block in the current blockchain from its hash;
// if there is no match, returns nil.
func (self *Blockchain)GetBlock(hash string)*Block{
  block := self.Head
  for{
    if block==nil{return nil}
    if block.GetHashCached()==hash{
      return block
    }
    block = block.Previous
  }
  return nil
}

// Implements the minig process. This process can end
// if the block becomes valid or if another valid block
// is received from a packet. If the process ends with a valid
// block mined, the block is sent to internalBlock channel.
func (self *Blockchain)Mine(){
  fmt.Println("Start mining")
  self.keepmining = true
  self.Current.Mine(self.id,&self.keepmining,self.internalMiniBlock,
    self.executor,self.config)
  // when mined, send blocks to internalBlock
  if self.keepmining == true{
    mex := new(MexBlock)
    mex.Data = self.Current.Serialize()
    self.internalBlock <- *mex
    self.keepmining = false
  }else{
    self.miningstatus <- false
  }
  fmt.Println("Mining ended")
}

// Serialize the current head block of the chain.
func (self *Blockchain)GetSerializedHead()[]byte{
  self.access_data.Lock()
  defer self.access_data.Unlock()
  return self.Head.Serialize()
}

// Main loop that manages the communication with
// the miner, reading the packets from channels and acting consequently.
// This method handles:
// - incoming blocks
// - mined block propagation
// - incoming transactions
func (self *Blockchain)Communicate(id utils.Addr,stop chan bool){
  for{
    select{
      case <-stop:// close message
        return
      case sol := <-self.chGoodSolutions: // a good solution for a job has been found
        self.processGoodSolutionFound(sol)
      case mex := <-self.internalBlock:// a block has been mined
        self.BlockOut <- mex
        self.startNewMiningProcess()
      case mex := <-self.BlockIn:// someone sent a block
        block,hashPrev := MarshalBlock(mex.Data,self.config)
        self.processIncomingBlock(block,hashPrev,mex.IpSender)
      case mex := <-self.TransQueue:// someone sent a transaction
        var transact Transaction
        switch mex.Type {
        case TrStd:
          transact = MarshalStdTransaction(mex.Data)
        case TrJob:
          transact = MarshalJobTransaction(mex.Data)
        case TrSol:
          transact = MarshalSolTransaction(mex.Data,self.config)
        case TrRes:
          transact = MarshalResTransaction(mex.Data)
        // TrPrize and TrCoin are inserted by miners in blocks directly,
        // thus they are not exchanged as single transaction.
        default:
          transact = nil
        }
        if transact!=nil{
          self.processIncomingTransaction(transact)
        }
      case mex := <- self.internalMiniBlock:
        // propagate self mined miniblock to other miners
        // do not propagate the last one: it is already in the block
        self.MiniBlockOut <- mex
      case mex := <- self.MiniBlockIn:
        miniblock := MarshalMiniBlock(mex.Data)
        if miniblock!=nil{
          self.access_data.Lock()
          block := self.Current
          self.access_data.Unlock()
          block.AddMiniBlock(miniblock,self.config)
        }
    }
  }
}

func (self *Blockchain)GetNextSlotForJob()(int,int){
  self.access_data.Lock()
  block := self.Head
  self.access_data.Unlock()
  return block.NextSlotForJobExectution()
}

// Called when mining process succeed to update the blockchain
// with the new current block; build a new one and restart mining
func (self *Blockchain)startNewMiningProcess(){
  self.storeCurrentBlockAndCreateNew(nil,0,nil)// with nil store self.Current
  go self.Mine()
}

// Ask another miner a block with a specific hash and returns
// it or nil (in case of error).
func askBlock(blockHash string,ipaddress string,config *conf.Config)(*Block,string){
  fmt.Println("asking for ",blockHash," to ",ipaddress)
  conn, err := net.Dial("tcp",ipaddress)
  if err!=nil{
    fmt.Println(err)
    return nil,"" }
  fmt.Fprintf(conn,"get_block\n"+blockHash+"\n")
  blockRaw,err2 := bufio.NewReader(conn).ReadString('\n')
  if err2!=nil{
    fmt.Println(err2)
    return nil,"" }
  return MarshalBlock([]byte(blockRaw),config)
}

// Checks the block and updates the blockchain.
// If it was valid and inserted as new head, restart mining
func (self *Blockchain)processIncomingBlock(block *Block,
                hashPrev string,ipSender string)  {
  var b *Block = block
  var hp string = hashPrev
  fmt.Println("checking incoming block")
  var savingPoint int = 0
  for{// try to reconstruct the blockchain if valid
    if b==nil{break}
    if !b.CheckStep1(hp){
      fmt.Println("check step 1 failed")
      fmt.Println("hhprev",hp)
      return
    }
    fmt.Println("step 1 parzial ok")
    fmt.Println("hhprev",hp)
    if hp=="" {break}
    existingBlock := self.GetBlock(hp)
    if existingBlock!=nil{
      b.Previous = existingBlock
      break
    }
    b.Previous,hp = askBlock(hp,ipSender,self.config)
    b = b.Previous
    savingPoint += 1
  }
  fmt.Println("chain has succeeded check step 1")
  transactionChanges := make(map[string]string)
  if block.CheckStep2(&transactionChanges,self.config){// the the blockchain is valid
    self.access_data.Lock()
    var replace bool = (block.LenSubChain>self.Head.LenSubChain)// || (
      // block.LenSubChain==self.Head.LenSubChain && block.NumJobs>self.Head.NumJobs))
    self.access_data.Unlock()
    fmt.Printf("Replace: %v\n",replace)
    if replace{
      // stop mining
      self.keepmining = false
      fmt.Println("keepmining:",self.keepmining)
      <-self.miningstatus // wait mining ending
      fmt.Println("Mining ended!")
      // update the blockchain with the new
      self.storeCurrentBlockAndCreateNew(block,savingPoint,&transactionChanges)
      // restart mining
      go self.Mine()
    }
  }else{
    fmt.Println("chain discarded in check step 2")
  }
}

// Set a block as new head and save to file the blocks from that one
// back to savingPoint blocks.
// Updates the map of current unspent transactions consequently.
func (self *Blockchain)storeCurrentBlockAndCreateNew(block *Block,
                        savingPoint int,trChanges *map[string]string){
  self.access_data.Lock()
  defer self.access_data.Unlock()
  if block!=nil{
    self.Head = block
    self.applyTransactionsChanges(trChanges)
  }else{
    self.Head = self.Current
    self.applyTransactionsChanges(&self.currentTrChanges)
  }
  self.Current = BuildBlock(self.id,self.Head)
  self.currentTrChanges = make(map[string]string)
  b := self.Head
  for i:=0;i<savingPoint+1;i++{
    if b==nil{break}
    var filename string = fmt.Sprintf("%v/block%v.json",self.folder,b.LenSubChain)
    data := b.Serialize()
    err := ioutil.WriteFile(filename,data, 0644)
    if err!=nil{fmt.Println(err)}
    b = b.Previous
  }
  self.checkJobsToClose()
  self.discloseDeclaredSolutionsAndIntegrateShared()
}

// Checks if in the last block (self.Head) you published a ResTransaction and,
// in case, publish now the corresponding SolTransaction.
// The solutions should be stored in self.solutionToShare.
// At the same time, checks if there are some shared solutions others
// disclosed, and adds them in the population of the job.
// This method is called inside self.access_data.Lock().
func (self *Blockchain)discloseDeclaredSolutionsAndIntegrateShared(){
  // loop through all transactions in Head and:
  transactions := self.Head.Transactions.GetTransactionArray()
  for i:=0;i<len(transactions);i++{
    switch transactions[i].GetType() {
    case TrRes:
      if transactions[i].GetCreator()==self.id{
        // disclose your solution
        hashjob := transactions[i].(*ResTransaction).JobTrans
        tr := MakeSolTransaction(self.id,self.config.GetPrivateKey(),
              self.Head.GetHashCached(),transactions[i].GetHashCached(),
              hashjob,self.solutionToShare[hashjob],self.config)
        delete(self.solutionToShare,hashjob)
        go self.propagateTransaction(tr)
      }
    case TrSol:
      if transactions[i].GetCreator()!=self.id{
        // integrate the solution from others
        tr := transactions[i].(*SolTransaction)
        hashJob := tr.JobTrans
        self.executor.InjectSharedSolution(hashJob,tr.Solution)
      }
    }
  }
}

// Decides to publish or not a solution in a ResTransaction and,
// in case, creates and propagates the transaction
func (self *Blockchain)processGoodSolutionFound(sol ga.Sol){
  fmt.Println("Process good solution")
  hashJob := sol.JobHash
  self.access_data.Lock()
  transactions := self.Current.Transactions.GetTransactionArray()
  openJobs := self.Current.getOpenJobs()
  self.access_data.Unlock()
  // check if the job is open and get block of the JobTransaction
  var open bool = false
  var jobBlock string
  for i:=0;i<len(openJobs);i++{
    if openJobs[i].GetHashCached()==hashJob{
      open = true
      jobBlock = openJobs[i].blockContainer
      break
    }
  }
  if open==false{ return }
  fmt.Println("The job is open")
  // find the current best score for the current block
  var bestScore float64
  var first bool = true
  var seen int = 0
  for i:=0;i<len(transactions);i++{
    if transactions[i].GetType()==TrSol{
      solTr := transactions[i].(*SolTransaction)
      if solTr.JobTrans==hashJob{
        seen += 1
        if first{
          bk := self.Current.FindPrevBlock(solTr.ResBlock)
          if bk!=nil{
            resTr := bk.FindTransaction(solTr.ResTrans)
            if resTr!=nil{
              bestScore = resTr.(*ResTransaction).Evaluation
              first = false
            }
          }
        }else{// bestSol != nil
          bk := self.Current.FindPrevBlock(solTr.ResBlock)
          if bk!=nil{
            resTr := bk.FindTransaction(solTr.ResTrans)
            if resTr!=nil{
              score := resTr.(*ResTransaction).Evaluation
              if (sol.IsMin && score<bestScore)||(sol.IsMin==false && score>bestScore){
                bestScore = score
              }
            }
          }
        }
      }
    }
  }
  fmt.Printf("Best score: %v your: %v\n",bestScore,sol.Fitness)
  if (seen<=0)||(sol.IsMin && sol.Fitness<bestScore)||(sol.IsMin==false && sol.Fitness>bestScore){
    // the solution is better than the current ones
    fmt.Println("The solution is good")
    // check if you possess enough money:
    var amount int = GetResTransactionCost()
    inps,tot := self.getUnspentCoin(amount)
    if tot<amount { return }
    fmt.Println("You have money")
    // make ResTransaction
    var out TrOutput
    out.Address = self.id
    out.Value = tot-amount // set back the remainder
    hb := new(utils.HashBuilder)
    serializedSol := sol.Individual.Serialize()
    hb.Add(serializedSol)
    hashSol := fmt.Sprintf("%x",hb.GetHash())
    resTr := MakeResTransaction(self.id,self.config.GetPrivateKey(),
          inps,out,jobBlock,hashJob,hashSol,sol.Fitness,sol.IsMin)
    // update self.solutionToShare
    self.solutionToShare[hashJob] = serializedSol
    // propagate the ResTransaction
    fmt.Println("Propagate")
    go self.propagateTransaction(resTr)
  }
}

// Searches backward in the chain some unspent TrOutputs untill an amount
// is reached or the chain has been completely explored. If the amount
// is covered a []TrInput is returned, with the total amount of money.
func (self *Blockchain)getUnspentCoin(amount int)([]TrInput,int){
  type TrINOUT struct{
    In TrInput
    Out TrOutput
  }
  b := self.Head
  var spent []TrInput
  var unspent []TrINOUT
  for{
    var usable []TrINOUT
    if b==nil{ return nil,0 }
    transactions := b.Transactions.GetTransactionArray()
    for i:=0;i<len(transactions);i++{
      switch transactions[i].GetType() {
      case TrStd:
        tr := transactions[i].(*StdTransaction)
        if tr.Creator==self.id{
          spent = append(spent,tr.Inputs...)
        }
        for j:=0;j<len(tr.Outputs);j++{
          if tr.Outputs[j].Address==self.id{
            var inout TrINOUT
            inout.Out = tr.Outputs[j]
            inout.In.Block = b.GetHashCached()
            inout.In.ToSpend = tr.GetHashCached()
            inout.In.Index = j
            usable = append(usable,inout)
          }
        }
      case TrCoin:
        out := transactions[i].GetOutputAt(0)
        if out.Address==self.id{
          var inout TrINOUT
          inout.Out.Address = out.Address
          inout.Out.Value = out.Value
          inout.Out.Address = out.Address
          inout.In.Block = b.GetHashCached()
          inout.In.ToSpend = transactions[i].GetHashCached()
          inout.In.Index = 0
          usable = append(usable,inout)
        }
      case TrJob:
        tr := transactions[i].(*JobTransaction)
        if tr.Creator==self.id{
          spent = append(spent,tr.Inputs...)
        }
        if tr.Output.Address==self.id{
          var inout TrINOUT
          inout.Out = tr.Output
          inout.In.Block = b.GetHashCached()
          inout.In.ToSpend = tr.GetHashCached()
          inout.In.Index = 0
          usable = append(usable,inout)
        }
      case TrRes:
        tr := transactions[i].(*ResTransaction)
        if tr.Creator==self.id{
          spent = append(spent,tr.Inputs...)
        }
        if tr.Output.Address==self.id{
          var inout TrINOUT
          inout.Out = tr.Output
          inout.In.Block = b.GetHashCached()
          inout.In.ToSpend = tr.GetHashCached()
          inout.In.Index = 0
          usable = append(usable,inout)
        }
      case TrPrize:
        out := transactions[i].GetOutputAt(0)
        if out.Address==self.id{
          var inout TrINOUT
          inout.Out.Address = out.Address
          inout.Out.Value = out.Value
          inout.Out.Address = out.Address
          inout.In.Block = b.GetHashCached()
          inout.In.ToSpend = transactions[i].GetHashCached()
          inout.In.Index = 0
          usable = append(usable,inout)
        }
      }
    }
    for i:=0;i<len(usable);i++{
      valid := true
      for j:=0;valid && j<len(spent);j++{
        if (usable[i].In.Block==spent[j].Block &&
            usable[i].In.ToSpend==spent[j].ToSpend &&
            usable[i].In.Index==spent[j].Index){
          valid = false
          break
        }
      }
      if valid{
        if usable[i].Out.Value>=amount{
          var result []TrInput
          result = append(result,usable[i].In)
          return result,usable[i].Out.Value
        }else{
          unspent = append(unspent,usable[i])
          var tot int = 0
          var result []TrInput
          for k:=0;k<len(unspent);k++{
            tot += unspent[k].Out.Value
            result = append(result,unspent[k].In)
          }
          if tot>=amount{
            return result,tot
          }
        }
      }
    }
    b = b.Previous
  }
  return nil,0
}

// Sends a transaction to the miner itself (localhost:port)
func (self *Blockchain)propagateTransaction(tr Transaction){
  conn, err := net.Dial("tcp","127.0.0.1:"+self.config.Port)
  if err!=nil{
    fmt.Println(err)
    return
  }
  fmt.Fprintf(conn,"transaction\n")
  fmt.Fprintf(conn,tr.GetType()+"\n")
  fmt.Fprintf(conn,string(tr.Serialize())+"\n")
}

// Checks if some of the jobs used in Head should stop and, in case,
// stops their execution and clears the executor.
// This method is called inside self.access_data.Lock().
func (self *Blockchain)checkJobsToClose(){
  oldJobs := self.Head.getOpenJobs() // Jobs that was open in the previous block
  jobs := self.Current.getOpenJobs() // Jobs that are open in the current block
  // find the jobs that need to be stopped:
  var toClose []*JobTransaction
  for i:=0;i<len(oldJobs);i++{
    var ended bool = true
    for j:=0;j<len(jobs);j++{
      if oldJobs[i].GetHashCached()==jobs[j].GetHashCached(){
        ended = false
        break
      }
    }
    if ended{
      toClose = append(toClose,oldJobs[i])
    }
  }
  // stop them
  for i:=0;i<len(toClose);i++{
    self.executor.StopJob(toClose[i].GetHashCached())
  }
}

// Checks the validity of a new transaction and insert it
// in the current block, if it is valid.
func (self *Blockchain)processIncomingTransaction(transaction Transaction) {
  fmt.Println("Processing transaction")
  trChanges := make(map[string]string)
  if transaction.Check(self.Current,&trChanges){
    self.access_data.Lock()
    defer self.access_data.Unlock()
    // check double spending in current unmined block
    var errDS error = nil
    for k,_ := range trChanges{
      if v2,ok := self.currentTrChanges[k]; ok{
        if v2!=""{
          errDS = errors.New("Double spending 3")
          break
        }
      }
    }
    if errDS==nil{
      err := self.Current.AddTransaction(transaction)
      if err!=nil{
        fmt.Println(err)
      }else{
        fmt.Println("Transaction inserted in current block")
        for k,v := range trChanges{
          self.currentTrChanges[k] = v
        }
      }
    }else{
      fmt.Println(errDS)
    }
  }else{
    fmt.Println("Transaction not pass check")
  }
}

// Update the transaction's spending block of the transactions
// in trChanges.
// This method is called inside the access_data lock.
func (self* Blockchain)applyTransactionsChanges(trChanges *map[string]string){
  for k,spendingBlock := range *trChanges {
    tokens := strings.Split(k,",")
    block := self.GetBlock(tokens[0])
    trans := block.FindTransaction(tokens[1])
    index,_ := strconv.Atoi(tokens[2])
    trans.GetOutputAt(index).SetSpentIn(spendingBlock)
  }
}
