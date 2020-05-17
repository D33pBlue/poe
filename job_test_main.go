/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-May-15
 * @Project: Proof of Evolution
 * @Filename: job_test.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-16
 * @Copyright: 2020
 */

package main

import(
  "fmt"
  "time"
  // "testing"
  "github.com/D33pBlue/poe/ga"
)

// func TestJob(t *testing.T){
func main(){
  job := ga.BuildJob("/home/d33pblue/Documenti/uni/tesi/poe/examples/tsp/tsp.go",
    "/home/d33pblue/Documenti/uni/tesi/poe/examples/tsp/data/tsp0/gr96.json")
  if job==nil{
    // t.Errorf("Unable to build job\n")
    fmt.Println("Unable to build job")
  }
  go job.Execute("hash","publicKey")
  for i:=0;i<10000;i++{
    // sol := <-job.ChNonce
    // fmt.Println("\n",i,sol.Fitness)
    time.Sleep(10*time.Second)
    done := false
    for ;!done;{
      select{
      case best := <-job.ChUpdateOut:
        fmt.Println("best",best.Fitness)
      default:
        done = true
      }
    }
  }
  job.KeepRunning = false
}
