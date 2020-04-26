/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-26
 * @Project: Proof of Evolution
 * @Filename: poe.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-26
 * @Copyright: 2020
 */

package main

import (
  "os"
  "fmt"
  "flag"
  "bufio"
  "io/ioutil"
  "strings"
  "strconv"
  "github.com/D33pBlue/poe/miner"
  "github.com/D33pBlue/poe/wallet"
  "github.com/D33pBlue/poe/utils"
)

func startMining(ip,port,keypath string){
  var publicKey utils.Addr
  if fileExists(keypath){
    data, err := ioutil.ReadFile(keypath)
    if err!=nil{
      fmt.Println(err)
      return
    }
    pub,err2 := utils.LoadPublicKeyFromPemStr(data)
    if err2!=nil {
      fmt.Println(err2)
      return
    }
    publicKey = utils.GetAddr2(pub)
  }else{
    fmt.Println("You need to link a valid public key file to start mining.")
    fmt.Println("You can generate it with mode genkey.")
    return
  }
  fmt.Println("Loaded public key:")
  fmt.Println(publicKey)
  minerNode := miner.New(port)
  go minerNode.Serve(publicKey)
  fmt.Printf("\nStarted mining node at port %v\n\n",port)
  startShell(processOnMining,minerNode)
}

func startWallet(ip,port,keypath string){
  if !fileExists(keypath){
    fmt.Println("You need to link a valid key file.")
    fmt.Println("You can generate it with mode genkey.")
    return
  }
  walletObj := wallet.New(keypath,ip+":"+port)
  if walletObj==nil{return}
  fmt.Printf("Connecting to %v:%v\n",ip,port)
  startShell(processOnWallet,walletObj)
}

func processOnMining(cmd string,args []string,obj interface{})string{
  switch cmd {
  case "node":
    if len(args)!=1 {
      return "invalid arguments"
    }
    err := obj.(*miner.Miner).AddNode(args[0])
    if err!=nil { return fmt.Sprint(err) }
    return "Added "+args[0]+" in list of miners"
  case "nodes":
    return strings.Join(obj.(*miner.Miner).GetConnected(),"\n")
  case "status":
    return "current status" // TODO:  implement later
  }
  return "invalid cmd"
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
    if len(args)!=1{
      return "invalid arguments"
    }
    err := obj.(*wallet.Wallet).SubmitJob(args[0])
    if err!=nil{ return fmt.Sprint(err) }
    return "Sent Job transaction"
  case "results":
    return "results" // TODO: implement later
  }
  return "invalid cmd"
}

func generateKey(){
  wallet := wallet.New("","")
  if wallet!=nil{
    fmt.Println("Generated public and private keys in ./data/")
  }
}

func main()  {
  fmt.Println("\n\n---------------------------------------")
  fmt.Println("---- Proof of Evolution Blockchain ----\n")
  mode := flag.String("mode", "", "Mode{mine|wallet|genkey}")
  ip := flag.String("ip", "127.0.0.1", "The IP address of the mining node")
  port := flag.String("port","4242","The port where the mining node start listening.")
  key := flag.String("key","","Path to the public key pem file")
  flag.Parse()
  if *mode=="mine"{
    startMining(*ip,*port,*key)
  }else if *mode=="wallet"{
    startWallet(*ip,*port,*key)
  }else if *mode=="genkey"{
    generateKey()
  }else{
    fmt.Println("Choose a mode in {mine|wallet|genkey}")
  }
}

func fileExists(filename string) bool {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
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
