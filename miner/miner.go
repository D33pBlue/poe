/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-18
 * @Project: Proof of Evolution
 * @Filename: miner.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-27
 * @Copyright: 2020
 */

package miner

import(
  "fmt"
  "net"
  "sync"
  "bufio"
  "regexp"
  "errors"
  "github.com/D33pBlue/poe/utils"
  "github.com/D33pBlue/poe/blockchain"
)

type Miner struct{
  Chain *blockchain.Blockchain
  Port string
  Connected []string // sinc read/write
  connected_lock sync.Mutex
  keepServing bool
  addrch chan string
}

func New(port string)*Miner{
  miner := new(Miner)
  miner.Port = port
  miner.Chain = blockchain.NewBlockchain()
  miner.keepServing = false
  miner.addrch = make(chan string)
  return miner
}

func (self *Miner)Serve(id utils.Addr)  {
  stopPropagation := make(chan bool)
  stopMining := make(chan bool)
  go self.Chain.Mine(id,stopMining)
  go self.propagateMinedBlocks(stopPropagation)
  l,err := net.Listen("tcp",":"+self.Port)
  if err!=nil{
    fmt.Println(err)
    return
  }
  defer l.Close()
  self.keepServing = true
  for {
    if !self.keepServing{
      stopMining <- true
      stopPropagation <- true
      break
    }
    conn,err := l.Accept()
    if err!=nil{
      fmt.Println(err)
      continue
    }
    go self.handleConnection(conn)
  }
}

func (self *Miner)GetConnected()[]string{
  self.connected_lock.Lock()
  var conn []string
  for i:=0;i<len(self.Connected);i++{
    conn = append(conn,self.Connected[i])
  }
  self.connected_lock.Unlock()
  return conn
}

func (self *Miner)AddNode(ipaddress string)error{
  match, _ := regexp.MatchString("[0-9]+.[0-9]+.[0-9]+.[0-9]+:[0-9]+",ipaddress)
  if !match{ return errors.New(ipaddress+" is not a valid address")}
  err := self.requestUpdate(ipaddress)
  if err!=nil{ return err }
  self.addrch <- ipaddress
  return nil
}

func (self *Miner)handleConnection(conn net.Conn){
  message, _ := bufio.NewReader(conn).ReadString('\n')
  fmt.Printf("Received %v\n",message)
  // TODO: process message
  conn.Write([]byte("ack\n"))
}

// ask to another miner his current blockchain
// and send a MexBlock through the rith channel
func (self *Miner)requestUpdate(ipaddress string)error{
  fmt.Println(ipaddress)
  return sendMexAck(ipaddress,"update")
}

func (self *Miner)propagateMinedBlocks(close chan bool){
  for{
    select{
      case <- close:
        return
      case ipaddress := <-self.addrch:
        self.connected_lock.Lock()
        self.Connected = append(self.Connected,ipaddress)
        self.connected_lock.Unlock()
      case block := <-self.Chain.BlockOut:
        self.connected_lock.Lock()
        for i:=0;i<len(self.Connected);i++{
          mex := fmt.Sprint(block) // TODO: implement later
          go sendMex(self.Connected[i],mex)
        }
        self.connected_lock.Unlock()
    }
  }
}


func sendMexAck(address,mex string)error{
  conn, err := net.Dial("tcp",address)
  if err!=nil{ return err }
  // send to socket
  fmt.Fprintf(conn,mex+"\n")
  // listen for reply
  ack, err2 := bufio.NewReader(conn).ReadString('\n')
  if err2!=nil{ return err2 }
  if ack!="ack\n"{
    return errors.New("Error in receiving ack: "+ack)
  }
  return nil
}

func sendMex(address,mex string){
  conn, err := net.Dial("tcp",address)
  if err!=nil{ return }
  // send to socket
  fmt.Fprintf(conn,mex+"\n")
}
