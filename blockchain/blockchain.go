/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-23
 * @Project: Proof of Evolution
 * @Filename: blockchain.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-01
 * @Copyright: 2020
 */


package blockchain

import (
  "fmt"
  "net"
  "sync"
  "bufio"
  "io/ioutil"
  "github.com/D33pBlue/poe/utils"
)

type any interface{}

type MexTrans struct{
  Type string
  Data []byte
}

type MexBlock struct{
  Data []byte
  IpSender string
}


type Blockchain struct{
  Head *Block
  Current *Block
  TransQueue chan MexTrans// receive transactions from miner
  BlockOut chan MexBlock// send mined block to miner
  BlockIn chan MexBlock// receive mined block from miner
  internalBlock chan MexBlock// notify mined block to Communicate
  id utils.Addr
  keepmining bool
  miningstatus chan bool
  access_data sync.Mutex
  folder string
}

func NewBlockchain(id utils.Addr,folder string)*Blockchain{
  chain := new(Blockchain)
  chain.TransQueue = make(chan MexTrans)
  chain.BlockOut = make(chan MexBlock)
  chain.BlockIn = make(chan MexBlock)
  chain.internalBlock = make(chan MexBlock)
  chain.miningstatus = make(chan bool)
  chain.Head = BuildFirstBlock(id)
  chain.folder = folder
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

func (self *Blockchain)GetBlock(hash string)*Block{
  block := self.Head
  for{
    if block==nil{return nil}
    // if utils.CompareHashes(block.GetHashCached(),hash){
    if block.GetHashCached()==hash{
      return block
    }
    block = block.Previous
  }
  return nil
}

func (self *Blockchain)Mine(){
  fmt.Println("Start mining")
  self.keepmining = true
  self.Current.Mine(&self.keepmining)
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

func (self *Blockchain)GetSerializedHead()[]byte{
  self.access_data.Lock()
  defer self.access_data.Unlock()
  return self.Head.Serialize()
}

func (self *Blockchain)GetTotal(addr utils.Addr)int{
  total := 0
  self.access_data.Lock()
  b := self.Head
  for{
    if b==nil{break}
    var transactions []Transaction = b.Transactions.transactions
    for i:=0;i<len(transactions);i++{
      total += GetTransactionMoneyForWallet(self,transactions[i],addr)
    }
    b = b.Previous
  }
  self.access_data.Unlock()
  return total
}

// Spend an output of a transaction and returns the value.
// To be called inside access_data lock
func (self *Blockchain)SpendSubTransaction(inp *TrInput,wallet utils.Addr)int{
  block := self.GetBlock(inp.Block)
  if block==nil{return 0}
  transact := block.GetTransaction(inp.ToSpend)
  if transact==nil{return 0}
  // if transact.IsSpent(){return 0}
  switch transact.GetType() {
  case TrCoin:
    var tr *CoinTransaction = transact.(*CoinTransaction)
    if tr.Output.Address==wallet{
      tr.Output.spent = true
      tr.spent = true
      return tr.Output.Value
    }
  case TrStd:
    var tr *StdTransaction = transact.(*StdTransaction)
    var val int = 0
    if tr.Outputs[inp.Index].Address==wallet{
      tr.Outputs[inp.Index].spent = true
      val = tr.Outputs[inp.Index].Value
    }
    allSpent := true
    for i:=0;i<len(tr.Outputs);i++{
      if !tr.Outputs[i].spent{
        allSpent = false
        break
      }
    }
    if allSpent{
      tr.spent = true
    }
    return val
  }
  // TODO: JobTransaction case
  return 0
}

// Called when mining process succeed to update the blockchain
// with the new current block, build a new one and restart mining
func (self *Blockchain)startNewMiningProcess(){
  self.storeCurrentBlockAndCreateNew(nil)// with nil store self.Current
  go self.Mine()
}

func (self *Blockchain)Communicate(id utils.Addr,stop chan bool){
  for{
    select{
      case <-stop:
        return
      case mex := <-self.internalBlock:
        self.BlockOut <- mex
        self.startNewMiningProcess()
      case mex := <-self.BlockIn:
        block,hashPrev := MarshalBlock(mex.Data)
        self.processIncomingBlock(block,hashPrev,mex.IpSender)
      case mex := <-self.TransQueue:
        var transact Transaction
        switch mex.Type {
        case TrStd:
          transact = MarshalStdTransaction(mex.Data)
        case TrJob:
          transact = MarshalJobTransaction(mex.Data)
        default:
          transact = nil
        }
        if transact!=nil{
          self.processIncomingTransaction(transact)
        }
    }
  }
}

func askBlock(blockHash string,ipaddress string)(*Block,string){
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
  return MarshalBlock([]byte(blockRaw))
}

// Checks the block and updates the blockchain.
// If valid restart mining
func (self *Blockchain)processIncomingBlock(block *Block,
                hashPrev string,ipSender string)  {
  var b *Block = block
  var hp string = hashPrev
  fmt.Println("checking incoming block")
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
    b.Previous,hp = askBlock(hp,ipSender)
    b = b.Previous
  }
  fmt.Println("chain has succeeded check step 1")
  if block.CheckStep2(){// the the blockchain is valid
    self.access_data.Lock()
    var replace bool = (block.LenSubChain>self.Head.LenSubChain)// || (
      // block.LenSubChain==self.Head.LenSubChain && block.NumJobs>self.Head.NumJobs))
    self.access_data.Unlock()
    fmt.Printf("Replace: %v\n",replace)
    if replace{
      // stop mining
      self.keepmining = false
      <-self.miningstatus // wait mining ending
      // update the blockchain with the new
      self.storeCurrentBlockAndCreateNew(block)
      // restart mining
      go self.Mine()
    }
  }else{
    fmt.Println("chaind discarded in check step 2")
  }
}

func (self *Blockchain)storeCurrentBlockAndCreateNew(block *Block){
  self.access_data.Lock()
  if block!=nil{
    self.Head = block
  }else{
    self.Head = self.Current
  }
  self.Current = BuildBlock(self.id,self.Head)
  var filename string = fmt.Sprintf("%v/block%v.json",self.folder,self.Head.LenSubChain)
  data := self.Head.Serialize()
  err := ioutil.WriteFile(filename,data, 0644)
  if err!=nil{fmt.Println(err)}
  self.access_data.Unlock()
}

func (self *Blockchain)processIncomingTransaction(transaction Transaction) {
  // TODO: implement later
}
