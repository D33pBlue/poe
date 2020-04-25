/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-23
 * @Project: Proof of Evolution
 * @Filename: blockchain.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-25
 * @Copyright: 2020
 */


package blockchain

import (
  "github.com/D33pBlue/poe/utils"
)

type any interface{}

type MexTrans struct{
  Type string
  Data any
}

type MexBlock struct{
  Type string
  Data any
}


type Blockchain struct{
  Head *Block
  Current *Block
  TransQueue chan MexTrans
  BlockOut chan MexBlock
  BlockIn chan MexBlock
  // History map[[]byte]*Block
}

func NewBlockchain()*Blockchain{
  chain := new(Blockchain)
  chain.TransQueue = make(chan MexTrans)
  chain.BlockOut = make(chan MexBlock)
  chain.BlockIn = make(chan MexBlock)
  return chain
}

func (self *Blockchain)GetBlock(hash []byte)*Block{
  return nil // TODO: implement later
}

func (self *Blockchain)Mine(id utils.Addr,stop chan bool){
  for{
    select{
      case <-stop:
        return
        // TODO: other cases using self.<chan>
    }
  }
}
