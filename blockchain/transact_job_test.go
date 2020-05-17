/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-May-10
 * @Project: Proof of Evolution
 * @Filename: transac_job_test.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-13
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

func TestJobSaveToFile(t *testing.T){
  job,_ := ioutil.ReadFile("/home/d33pblue/Documenti/uni/tesi/poe/examples/tsp/tsp.go")
  data,_ := ioutil.ReadFile("/home/d33pblue/Documenti/uni/tesi/poe/examples/tsp/data/tsp0/gr96.json")
  key,_ := utils.GenerateKey()
  addr := utils.GetAddr(key)
  var out TrOutput
  jobtr := MakeJobTransaction(addr,key,nil,out,
    fmt.Sprintf("%x",job),fmt.Sprintf("%x",data),"",100,40,45)
  err := jobtr.SaveJobInFile("/home/d33pblue/Documenti/uni/tesi/poe/examples/tsp/tsp_trcopy.go")
  if err!=nil{
    t.Errorf("%v\n",err)
  }
  if !utils.FileExists("/home/d33pblue/Documenti/uni/tesi/poe/examples/tsp/tsp_trcopy.go"){
    t.Errorf("The file has not been saved\n")
  }
  job2,_ := ioutil.ReadFile("/home/d33pblue/Documenti/uni/tesi/poe/examples/tsp/tsp_trcopy.go")
  if len(job)!=len(job2){
    t.Errorf("The saved job has a different length than the original one\n")
  }else{
    for i:=0;i<len(job);i++{
      if job2[i]!=job[i]{
        t.Errorf("The saved job is not equal to the original one\n")
        break
      }
    }
  }
  err2 := jobtr.SaveDataInFile("/home/d33pblue/Documenti/uni/tesi/poe/examples/tsp/data/tsp0/gr96_tr_copy.json")
  if err2!=nil{
    t.Errorf("%v\n",err2)
  }
  if !utils.FileExists("/home/d33pblue/Documenti/uni/tesi/poe/examples/tsp/data/tsp0/gr96_tr_copy.json"){
    t.Errorf("The file has not been saved\n")
  }
  data2,_ := ioutil.ReadFile("/home/d33pblue/Documenti/uni/tesi/poe/examples/tsp/data/tsp0/gr96_tr_copy.json")
  if len(data)!=len(data2){
    t.Errorf("The saved data has a different length than the original one\n")
  }else{
    for i:=0;i<len(data);i++{
      if data2[i]!=data[i]{
        t.Errorf("The saved data is not equal to the original one\n")
        break
      }
    }
  }
}
