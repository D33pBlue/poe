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
  "strings"
  "strconv"
  "github.com/D33pBlue/poe/miner"
  "github.com/D33pBlue/poe/wallet"
  "github.com/D33pBlue/poe/utils"
)

func startMining(ip,port,keypath string){
  if fileExists(keypath){
    fmt.Println("key:",keypath)
  }else{
    fmt.Println("You need to link a valid public key file to start mining.")
    fmt.Println("You can generate it with mode genkey.")
    return
  }
  fmt.Printf("Starting mining node at %v:%v\n",ip,port)
  minerNode := miner.New(port)
  fmt.Println(minerNode)
  // start go routines for mining..
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
  fmt.Println("---- Proof of Evolution ----")
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
