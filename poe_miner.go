/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-26
 * @Project: Proof of Evolution
 * @Filename: poe.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Jun-27
 * @Copyright: 2020
 */


package main

import (
  "os"
  "fmt"
  "flag"
  "bufio"
  "strings"
  "github.com/D33pBlue/poe/miner"
  "github.com/D33pBlue/poe/utils"
  "github.com/D33pBlue/poe/conf"
)

func startMining(config *conf.Config){
  var publicKey utils.Addr = config.GetPublicKey()
  fmt.Println("Loaded public key:")
  fmt.Println(publicKey)
  var port string = config.GetPort()
  minerNode := miner.New(port,publicKey,config)
  go minerNode.Serve()// start serving in a goroutine
  fmt.Printf("\nStarted mining node at port %v\n\n",port)
  // add linked miners and request updates
  linkedMinersIp := config.GetLinkedMinersIp()
  for i:=0;i<len(linkedMinersIp);i++{
    processOnMining("node",[]string{linkedMinersIp[i]},minerNode)
  }
  // start a shell for user inputs
  startShell(processOnMining,minerNode)
}

// Maps user inputs to actions. This is a resolver.
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
  configFile := flag.String("conf","conf/config0.json","path to the config file")
  flag.Parse()
  config,err := conf.LoadConfiguration(*configFile,"")
  if err!=nil{
    fmt.Println(err)
    return
  }
  startMining(config)
}

type resolver func(string,[]string,interface{})string

// Reads user inputs and calls a resolver to process them.
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
