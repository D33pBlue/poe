/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-18
 * @Project: Proof of Evolution
 * @Filename: block.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-27
 * @Copyright: 2020
 */

// Package block contains the definitions
// of blocks, transactions, nonce
package blockchain

import(
  "time"
  "sync"
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
  Hash []byte
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
  block.Hash = block.GetHash()
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
  block.Hash = block.GetHash()
  block.checked = false
  block.mined = false
  return block
}

func (self *Block)Mine(){
  if self.NumJobs==0{
    self.mineNoJob()
  }else{
    self.mineWithJobs()
  }
}

func (self *Block)Serialize()[]byte{
  return nil // TODO: implement later
}

func MarshalBlock(data []byte)*Block{
  return nil // TODO: implement later
}

func (self *Block)mineWithJobs(){
  // TODO: implement later
}

func (self *Block)mineNoJob(){
  for{
    if self.mined{
      break
    }
    self.access_data.Lock()
    self.NonceNoJob.Next()
    self.Hash = self.GetHash()
    self.mined = self.checkNoJob()
    self.access_data.Unlock()
  }
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

// The Check method checks the validity of the block.
func (self *Block)Check()bool{
  if self.checked { return true }
  if !utils.CompareHashes(self.Hash,self.GetHash()){ return false }
  if self.NumJobs==0{
    if !self.checkNoJob() { return false }
  }else if self.NumJobs>0{
    if !self.checkJobs() { return false }
  }else{ return false }
  if !self.checkNumJobs() { return false }
  if !self.checkHardness() { return false }
  if !self.checkTransactions() { return false }
  if self.LenSubChain>0{
    if self.Previous==nil { return false }
    if self.Previous.LenSubChain!=self.LenSubChain-1 { return false }
    if !self.Previous.Check(){ return false } // TODO: remove recursion
  }
  self.checked = true
  return true
}

func (self *Block)GetTransaction(hash []byte)Transaction{
  return nil // TODO: implement later
}

func (self *Block)GetHash()[]byte{
  hb := new(utils.HashBuilder)
  hb.Add(self.Timestamp)
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
    hb.Add(self.Previous.GetHashCached())
  }
  return hb.GetHash()
}

func (self *Block)GetHashCached()[]byte{
  return self.Hash
}

func (self *Block)checkNoJob()bool{
  for i:=0;i<self.Hardness;i++{
    if self.Hash[i]!=0{ return false }
  }
  return true
}

func (self *Block)checkJobs()bool{
  return true // TODO: implement later
}

func (self *Block)checkNumJobs()bool{
  return true // TODO: implement later
}

func (self *Block)checkHardness()bool{
  return true // TODO: implement later
}

func (self *Block)checkTransactions()bool{
  return true // TODO: implement later
}

func (self *Block)calculateNumJobs()int{
  return 0 // TODO: implement later
}

func (self *Block)calculateHardness()int{
  return 0 // TODO: implement later
}

// Returns the value in coin of mining that block.
// This value depends on the hardness of the mining task.
func (self *Block)calculateMiningValue()int{
  return (1+self.Hardness)*10 // TODO: tune with mining time
}
