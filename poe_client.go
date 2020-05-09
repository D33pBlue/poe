/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-May-08
 * @Project: Proof of Evolution
 * @Filename: client.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-09
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
)


func startWallet(ip,port,keypath string,trusted bool){
  if !utils.FileExists(keypath){
    fmt.Println("You need to link a valid key file.")
    fmt.Println("You can generate it with mode genkey.")
    return
  }
  walletObj := wallet.New(keypath,ip+":"+port,trusted)
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
  wallet := wallet.New("","",true)
  if wallet!=nil{
    fmt.Println("Generated public and private keys in ./data/")
  }
}

func main()  {
  fmt.Println("\n\n---------------------------------------")
  fmt.Println("------ Proof of Evolution Client ------\n")
  mode := flag.String("mode", "wallet", "Mode{wallet|genkey}")
  ip := flag.String("ip", "127.0.0.1", "The IP address of the mining node")
  trusted := flag.Bool("trusted",true,"Set to true only if the miner is trusted")
  port := flag.String("port","4242","The port where the mining node start listening.")
  key := flag.String("key","","Path to the public key pem file")
  flag.Parse()
  if *mode=="wallet"{
    startWallet(*ip,*port,*key,*trusted)
  }else if *mode=="genkey"{
    generateKey()
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
