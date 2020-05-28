/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-May-08
 * @Project: Proof of Evolution
 * @Filename: client.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-28
 * @Copyright: 2020
 */

package main

import (
  "os"
  "fmt"
  "flag"
  "bufio"
  // "io/ioutil"
  "strings"
  "strconv"
  "github.com/D33pBlue/poe/wallet"
  "github.com/D33pBlue/poe/utils"
  "github.com/D33pBlue/poe/conf"
)


func startWallet(ip string,config *conf.Config,trusted bool){
  var port string = config.GetPort()
  walletObj := wallet.New(config,ip+":"+port,trusted)
  if walletObj==nil{return}
  fmt.Printf("Connecting to %v:%v\n",ip,port)
  startShell(processOnWallet,walletObj)
}


func processOnWallet(cmd string,args []string,obj interface{})string{
  switch cmd {
  case "public":
    return fmt.Sprint(obj.(*wallet.Wallet).Id)
  case "total":
    return fmt.Sprint(obj.(*wallet.Wallet).GetTotal())
  case "money":
    if len(args)!=2{
      return "invalid arguments"
    }
    amount,err := strconv.Atoi(args[0])
    if err!=nil{ return fmt.Sprint(err) }
    var receiver utils.Addr = utils.Addr(args[1])
    err = obj.(*wallet.Wallet).SendMoney(amount,receiver)
    if err!=nil{ return fmt.Sprint(err) }
    return fmt.Sprintf("Sent transaction of %v to %v",amount,receiver)
  case "job":
    if len(args)!=3{
      return "invalid arguments <path to job> <data> <prize>"
    }
    prize,errc := strconv.Atoi(args[2])
    if errc!=nil{ return fmt.Sprint(errc) }
    err := obj.(*wallet.Wallet).SubmitJob(args[0],args[1],prize)
    if err!=nil{ return fmt.Sprint(err) }
    return "Sent Job transaction"
  case "estimate":
    if len(args)!=2{
      return "invalid arguments <path to job> <data>"
    }
    fixedCost,advPrize := obj.(*wallet.Wallet).EstimateJobCost(args[0],args[1])
    return fmt.Sprintf("Fixed cost: %v\nMinimal prize: %v\nMin Total: %v",fixedCost,advPrize,fixedCost+advPrize)
  case "jobs":
    return obj.(*wallet.Wallet).GetSubmittedJobs()
  case "results":
    if len(args)!=1{
      return "invalid arguments <job id>"
    }
    return obj.(*wallet.Wallet).FetchAndStoreResults(args[0])
  }
  return "invalid cmd"
}

func generateKey(config *conf.Config){
  wallet := wallet.New(config,"",true)
  if wallet!=nil{
    fmt.Println("Generated public and private keys")
  }
}

func main()  {
  fmt.Println("\n\n---------------------------------------")
  fmt.Println("------ Proof of Evolution Client ------\n")
  mode := flag.String("mode", "wallet", "Mode{wallet|genkey}")
  ip := flag.String("ip", "127.0.0.1", "The IP address of the mining node")
  trusted := flag.Bool("trusted",true,"Set to true only if the miner is trusted")
  configFile := flag.String("conf","conf/config0.json","path to the config file")
  flag.Parse()
  config,err := conf.LoadConfiguration(*configFile)
  if err!=nil{
    fmt.Println(err)
    return
  }
  if *mode=="wallet"{
    startWallet(*ip,config,*trusted)
  }else if *mode=="genkey"{
    config.Key = ""
    generateKey(config)
  }else{
    fmt.Println("Choose a mode in {wallet|genkey}")
  }
}

type resolver func(string,[]string,interface{})string

func startShell(f resolver,obj interface{}){
  reader := bufio.NewReader(os.Stdin)
  fmt.Println("Accepting commands")
  fmt.Println("---------------------")
  var stop bool = false
  for ;!stop;{
    fmt.Print("-> ")
    text, _ := reader.ReadString('\n')
    text = strings.Replace(text, "\n", "", -1)
    args := strings.Split(text," ")
    if len(args)>1{
      fmt.Println(f(args[0],args[1:],obj))
    }else if len(args)==1{
      if args[0]=="close"{
        stop = true
      }else{
        fmt.Println(f(args[0],nil,obj))
      }
    }
  }
}
