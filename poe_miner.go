/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-26
 * @Project: Proof of Evolution
 * @Filename: poe.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-08
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
  // "strconv"
  "github.com/D33pBlue/poe/miner"
  "github.com/D33pBlue/poe/utils"
)

func startMining(ip,port,keypath string){
  var publicKey utils.Addr
  if utils.FileExists(keypath){
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
    fmt.Println("You can generate it with a client.")
    return
  }
  fmt.Println("Loaded public key:")
  fmt.Println(publicKey)
  minerNode := miner.New(port,publicKey)
  go minerNode.Serve()
  fmt.Printf("\nStarted mining node at port %v\n\n",port)
  startShell(processOnMining,minerNode)
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
    return "<current status>" // TODO:  implement later
  }
  return "invalid cmd"
}

func main()  {
  fmt.Println("\n\n---------------------------------------")
  fmt.Println("---- Proof of Evolution Blockchain ----\n")
  ip := flag.String("ip", "127.0.0.1", "The IP address of the mining node")
  port := flag.String("port","4242","The port where the mining node start listening.")
  key := flag.String("key","","Path to the public key pem file")
  flag.Parse()
  startMining(*ip,*port,*key)
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
