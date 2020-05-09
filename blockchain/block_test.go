/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-27
 * @Project: Proof of Evolution
 * @Filename: block_test.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-09
 * @Copyright: 2020
 */

package blockchain

import(
  // "fmt"
  "testing"
  "github.com/D33pBlue/poe/utils"
)

func TestSerialization(t *testing.T)  {
  key,_ := utils.GenerateKey()
  first := BuildFirstBlock(utils.GetAddr(key))
  block := BuildBlock(utils.GetAddr(key),first)
  ser := block.Serialize()
  block2,_ := MarshalBlock(ser)
  if block.LenSubChain!=block2.LenSubChain{
    t.Errorf("%v != %v\n",block.LenSubChain,block2.LenSubChain)
  }
  if block.NumJobs!=block2.NumJobs{
    t.Errorf("%v != %v\n",block.NumJobs,block2.NumJobs)
  }
  if block.Hardness!=block2.Hardness{
    t.Errorf("%v != %v\n",block.Hardness,block2.Hardness)
  }
  if !block.Timestamp.Equal(block2.Timestamp){
    t.Errorf("%v != %v\n",block.Timestamp,block2.Timestamp)
  }
  if !utils.CompareHashes(block.Hash,block2.Hash){
    t.Errorf("%v != %v\n",block.Hash,block2.Hash)
  }
  // fmt.Println(block.Transactions.Root.Transaction)
  // fmt.Println(block2.Transactions.Root.Transaction)
}

func TestBlockHash(t *testing.T){
  key,_ := utils.GenerateKey()
  addr := utils.GetAddr(key)
  first := BuildFirstBlock(addr)
  block := BuildBlock(addr,first)
  hash1 := block.GetHash("")
  hash2 := block.GetHash("")
  for i:=0;i<10;i++{
    transact := MakeStdTransaction(addr,key,nil,nil)
    block.AddTransaction(transact)
  }
  hash3 := block.GetHash("")
  if hash1!=hash2{
    t.Errorf("Inconsistent hash\n")
  }
  if hash1==hash3{
    t.Errorf("Same hashes after inserting transaction!")
  }
  if !block.Transactions.Check(){
    t.Error("Bug in merkle tree check or hash-update")
  }
  block2,_ := MarshalBlock(block.Serialize())
  if !block2.Transactions.Check(){
    t.Error("Bug in merkle tree check or hash-update after deserialization")
  }
}
