/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-23
 * @Project: Proof of Evolution
 * @Filename: blockchain.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-28
 * @Copyright: 2020
 */


package blockchain

import (
  "github.com/D33pBlue/poe/utils"
)

type any interface{}

type MexTrans struct{
  Type string
  Data []byte
}

type MexBlock struct{
  Data []byte
}


type Blockchain struct{
  Head *Block
  Current *Block
  TransQueue chan MexTrans// receive transactions from miner
  BlockOut chan MexBlock// send mined block to miner
  BlockIn chan MexBlock// receive mined block from miner
  internalBlock chan MexBlock// notify mined block to Communicate
  id utils.Addr
  // history map[[]byte]*Block
}

func NewBlockchain(id utils.Addr)*Blockchain{
  chain := new(Blockchain)
  chain.TransQueue = make(chan MexTrans)
  chain.BlockOut = make(chan MexBlock)
  chain.BlockIn = make(chan MexBlock)
  chain.internalBlock = make(chan MexBlock)
  chain.Head = BuildFirstBlock(id)
  chain.Current = BuildBlock(id,chain.Head)
  chain.id = id
  return chain
}

func (self *Blockchain)GetBlock(hash []byte)*Block{
  return nil // TODO: implement later
}

func (self *Blockchain)Mine(stop *bool){
  // TODO: implement later
  // ẁhen mined, send blocks to internalBlock
}

func (self *Blockchain)startNewMiningProcess(){
  // TODO: implement later
  // store current block, generate new current, add trandactions
  // go self.Mine()
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
        self.processIncomingBlock(block,hashPrev)
      case mex := <-self.TransQueue:
        var transact Transaction
        switch mex.Type {
        case TrStd:
          transact = MarshalStdTransaction(mex.Data)
        case TrCoin:
          transact = MarshalCoinTransaction(mex.Data)
        case TrJob:
          transact = MarshalJobTransaction(mex.Data)
        }
        self.processIncomingTransaction(transact)
    }
  }
}

func (self *Blockchain)processIncomingBlock(block *Block,hashPrev []byte)  {
  // TODO: check the block and update the blockchain.
  // if valid restart mining
}

func (self *Blockchain)processIncomingTransaction(transaction Transaction) {
  // TODO: implement later
}
