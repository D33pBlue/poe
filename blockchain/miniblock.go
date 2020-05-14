/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-25
 * @Project: Proof of Evolution
 * @Filename: miniblock.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-14
 * @Copyright: 2020
 */

package blockchain

import(
  "fmt"
  "github.com/D33pBlue/poe/utils"
  "github.com/D33pBlue/poe/ga"
)

type MiniBlock struct{
  // JobTransaction
  HashPrevBlock string
  Miner utils.Addr
  JobBlock string // the hash of the block with the JobTransaction
  JobTrans string // the hash of the JobTransaction
  Hash string
  Nonce Nonce
}

// Creates a new MiniBlock and initialize its channels.
func BuildMiniBlock(chNonce chan ga.Sol)*MiniBlock{
  return nil // TODO: implement later
}

// When this miniblock is mined it is sent to chOut.
// If stopMining==true the mining should stop and send nil to ch.
func Mine(hardness int,stopMining bool,chOut chan *MiniBlock){
  // TODO: implement later
}

// Checks the validity of the hash of the nonce in relation
// to the hardness, but not the validity of the evaluation of the solution.
// Also check if the hash of the previous block matches.
func (self *MiniBlock)CheckStep1(hashPrev string,hardness int)bool{
  // check the hash of the previous block
  if self.HashPrevBlock!=hashPrev{
    return false
  }
  // check the hash of this block
  if self.GetHash()!=self.GetHashCached(){
    return false
  }
  return self.checkHashPuzzle(hardness)
}

// Given an hardness, checks the validity of the hash of this MiniBlock.
func (self *MiniBlock)checkHashPuzzle(hardness int)bool {
  for i:=0;i<hardness;i++{
    if self.Hash[i]!='0'{
      return false
    }
  }
  return true
}

// Complete the check evaluating the stored solution.
// The block given must be the head of the complete chain in
// order to load the JobTransaction.
func (self *MiniBlock)CheckStep2(block *Block)bool{
  return true // TODO: implement later
}

// Recalculates and returns the hash of the MiniBlock.
func (self *MiniBlock)GetHash()string{
  hb := new(utils.HashBuilder)
  hb.Add(self.HashPrevBlock)
  hb.Add(self.Miner)
  hb.Add(self.JobBlock)
  hb.Add(self.JobTrans)
  hb.Add(self.Nonce.Solution)
  hb.Add(self.Nonce.Evaluation)
  hb.Add(self.Nonce.Complexity)
  hash := hb.GetHash()
  return fmt.Sprintf("%x",hash)
}

// Returns the cached hash of the MiniBlock.
func (self *MiniBlock)GetHashCached()string{
  return self.Hash
}
