/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-May-10
 * @Project: Proof of Evolution
 * @Filename: transac_job_test.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-10
 * @Copyright: 2020
 */

package blockchain

import(
  "testing"
  "io/ioutil"
  "fmt"
  "github.com/D33pBlue/poe/utils"
)


func TestJobTransactionHash(t *testing.T){
  job,_ := ioutil.ReadFile("/home/d33pblue/Documenti/uni/tesi/poe/examples/tsp/tsp.go")
  data,_ := ioutil.ReadFile("/home/d33pblue/Documenti/uni/tesi/poe/examples/tsp/data/tsp0/gr96.json")
  key,_ := utils.GenerateKey()
  addr := utils.GetAddr(key)
  var out TrOutput
  jobtr := MakeJobTransaction(addr,key,nil,out,
    fmt.Sprintf("%x",job),fmt.Sprintf("%x",data),"",100,40,45)
  if jobtr.Hash!=jobtr.GetHash(){
    t.Errorf("Hash mismatch")
  }
  jobtr2 := MarshalJobTransaction(jobtr.Serialize())
  if jobtr2.Hash!=jobtr.Hash{
    t.Errorf("Different hashes after serialization")
  }
  if jobtr2.Hash!=jobtr2.GetHash(){
    t.Errorf("Hash mismatch after serialization")
  }
}
