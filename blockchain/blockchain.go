/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-23
 * @Project: Proof of Evolution
 * @Filename: blockchain.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-08
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
  "io/ioutil"
  "github.com/D33pBlue/poe/utils"
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
// - a lock for sync and some booleand to manage the mining process.
// (The miner has not to give the total money of clients: it is a problem
// of the client to imoplement this, by collecting old blocks/transactions).
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
  currentTrChanges map[string]string
}

// This method load a serialized blockchain from file.
// (Each time a block is added to the blockchain, in fact, it is
// saved in a file with its index as name).
func LoadChainFromFolder(id utils.Addr,folder string)*Blockchain{
  chain := new(Blockchain)
  chain.TransQueue = make(chan MexTrans)
  chain.BlockOut = make(chan MexBlock)
  chain.BlockIn = make(chan MexBlock)
  chain.internalBlock = make(chan MexBlock)
  chain.miningstatus = make(chan bool)
  chain.folder = folder
  chain.currentTrChanges = make(map[string]string)
  chain.Head = nil
  var i int = 0
  fmt.Println("Loading chain from disk")
  for{
    var filename string = fmt.Sprintf("%v/block%v.json",chain.folder,i)
    if !utils.FileExists(filename){break}
    data,err := ioutil.ReadFile(filename)
    if err!=nil{return nil}
    b,hash := MarshalBlock(data)
    if chain.Head!=nil && chain.Head.Hash!=hash{return nil}
    b.Previous = chain.Head
    if !b.CheckStep1(hash){ return nil }
    chain.Head = b
    fmt.Printf("Loaded block %v\n",filename)
    i += 1
  }
  fmt.Println("Checking loaded chain")
  trChanges := make(map[string]string)
  if !chain.Head.CheckStep2(&trChanges){ return nil }
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
func NewBlockchain(id utils.Addr,folder string)*Blockchain{
  chain := new(Blockchain)
  chain.TransQueue = make(chan MexTrans)
  chain.BlockOut = make(chan MexBlock)
  chain.BlockIn = make(chan MexBlock)
  chain.internalBlock = make(chan MexBlock)
  chain.miningstatus = make(chan bool)
  chain.Head = BuildFirstBlock(id)
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

// Serialize the current head block of the chain.
func (self *Blockchain)GetSerializedHead()[]byte{
  self.access_data.Lock()
  defer self.access_data.Unlock()
  return self.Head.Serialize()
}

// Main loop (called by miner), that manages the communication with
// the miner, reading the packets from channels and acting consequently.
func (self *Blockchain)Communicate(id utils.Addr,stop chan bool){
  for{
    select{
      case <-stop:// close message
        return
      case mex := <-self.internalBlock:// a block has been mined
        self.BlockOut <- mex
        self.startNewMiningProcess()
      case mex := <-self.BlockIn:// someone sent a block
        block,hashPrev := MarshalBlock(mex.Data)
        self.processIncomingBlock(block,hashPrev,mex.IpSender)
      case mex := <-self.TransQueue:// someone sent a transaction
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

// Called when mining process succeed to update the blockchain
// with the new current block; build a new one and restart mining
func (self *Blockchain)startNewMiningProcess(){
  self.storeCurrentBlockAndCreateNew(nil,0,nil)// with nil store self.Current
  go self.Mine()
}

// Ask another miner a block with a specific hash and returns
// it or nil (in case of error).
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
    b.Previous,hp = askBlock(hp,ipSender)
    b = b.Previous
    savingPoint += 1
  }
  fmt.Println("chain has succeeded check step 1")
  transactionChanges := make(map[string]string)
  if block.CheckStep2(&transactionChanges){// the the blockchain is valid
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
}

// Checks the validity of a new transaction and insert it
// in the current block, if it is valid.
func (self *Blockchain)processIncomingTransaction(transaction Transaction) {
  // TODO: implement later
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
