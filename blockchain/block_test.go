/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-27
 * @Project: Proof of Evolution
 * @Filename: block_test.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-28
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
}
