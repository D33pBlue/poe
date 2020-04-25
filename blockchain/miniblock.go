/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-25
 * @Project: Proof of Evolution
 * @Filename: miniblock.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-25
 * @Copyright: 2020
 */

package blockchain

import(
  "github.com/D33pBlue/poe/utils"
)

type MiniBlock struct{
  // JobTransaction
  HashPrevBlock []byte
  Miner utils.Addr
  Hash []byte
  Nonce *Nonce
}

// The Check method checks the validity of the miniblock.
func (self *MiniBlock)Check(prevBlock *Block)bool{
  return true // TODO: implement later
}
