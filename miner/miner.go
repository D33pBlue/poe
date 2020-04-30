/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-18
 * @Project: Proof of Evolution
 * @Filename: miner.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-30
 * @Copyright: 2020
 */

package miner

import(
  "os"
  "fmt"
  "net"
  "sync"
  "bufio"
  "regexp"
  "strings"
  "strconv"
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
  id utils.Addr
}

func New(port string,id utils.Addr)*Miner{
  miner := new(Miner)
  miner.Port = port
  miner.id = id
  var folder string = fmt.Sprintf("data/chain%v",port)
  // make dir if not exists
  _, err := os.Stat(folder)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(folder, 0755)
		if errDir != nil {fmt.Println(errDir)}
	}
  miner.Chain = blockchain.NewBlockchain(id,folder)
  miner.keepServing = false
  miner.addrch = make(chan string)
  return miner
}

func (self *Miner)Serve()  {
  stopPropagation := make(chan bool)
  stopMining := make(chan bool)
  go self.Chain.Communicate(self.id,stopMining)
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

func (self *Miner)updateAddresses(conn net.Conn, port string){
  var ipaddress = fmt.Sprint(conn.RemoteAddr())
  ipaddress = ipaddress[:strings.Index(ipaddress,":")]+":"+port
  self.connected_lock.Lock()
  found := false
  for i:=0;i<len(self.Connected);i++{
    if self.Connected[i]==ipaddress{
      found = true
      break
    }
  }
  if !found{
    self.Connected = append(self.Connected,ipaddress)
    fmt.Println("Added node ",ipaddress)
  }
  self.connected_lock.Unlock()
}

func (self *Miner)handleConnection(conn net.Conn){
  reader := bufio.NewReader(conn)
  message, _ := reader.ReadString('\n')
  fmt.Printf("Received %v",message)
  switch message[:len(message)-1] {
  case "update":
    port, _ := reader.ReadString('\n')
    self.updateAddresses(conn,port[:len(port)-1])
    conn.Write([]byte(self.Chain.GetSerializedHead()))
    conn.Write([]byte("\n"))
  case "get_block":
    hash, _ := reader.ReadString('\n')
    fmt.Printf("Requested %v\n",hash)
    block := self.Chain.GetBlock(hash[:len(hash)-1])
    fmt.Printf("Found %v\n",block)
    if block==nil{
      conn.Write([]byte("\n"))
    }else{
      conn.Write([]byte(block.Serialize()))
    }
    conn.Write([]byte("\n"))
  case "chain":
    port,_ := reader.ReadString('\n')
    block,err2 := reader.ReadString('\n')
    fmt.Println(err2)
    fmt.Printf("Content of chain: %v\n",block)
    if err2==nil{
      mexBlock := new(blockchain.MexBlock)
      mexBlock.Data = []byte(block)
      var ipaddress string = fmt.Sprint(conn.RemoteAddr())
      mexBlock.IpSender = ipaddress[:strings.Index(ipaddress,":")]+":"+port[:len(port)-1]
      self.Chain.BlockIn <- *mexBlock
    }
  case "get_total":
    publicKey,_ := reader.ReadString('\n')
    var addr utils.Addr = utils.Addr(publicKey[:len(publicKey)-1])
    total := self.Chain.GetTotal(addr)
    conn.Write([]byte(strconv.Itoa(total)+"\n"))
  }
}

// ask to another miner his current blockchain
// and send a MexBlock through self.Chain.BlockIn channel
func (self *Miner)requestUpdate(ipaddress string)error{
  conn, err := net.Dial("tcp",ipaddress)
  if err!=nil{ return err }
  // send to socket
  fmt.Fprintf(conn,"update\n")
  fmt.Fprintf(conn,self.Port+"\n")
  // listen for reply
  fmt.Println("Listening")
  block,err2 := bufio.NewReader(conn).ReadString('\n')
  fmt.Println(err2)
  fmt.Printf("Received %v\n",block)
  if err2!=nil{ return err2 }
  mexBlock := new(blockchain.MexBlock)
  mexBlock.Data = []byte(block)
  mexBlock.IpSender = fmt.Sprint(conn.RemoteAddr())
  self.Chain.BlockIn <- *mexBlock
  return nil
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
      case blockMex := <-self.Chain.BlockOut:
        self.connected_lock.Lock()
        for i:=0;i<len(self.Connected);i++{
          fmt.Println("Send update to ",self.Connected[i])
          go self.sendBlockUpdate(self.Connected[i],string(blockMex.Data))
        }
        self.connected_lock.Unlock()
    }
  }
}


func (self *Miner)sendMexAck(address,mex string)error{
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

func (self *Miner)sendBlockUpdate(address string,mex string){
  conn, err := net.Dial("tcp",address)
  if err!=nil{ return }
  fmt.Fprintf(conn,"chain\n")
  fmt.Fprintf(conn,self.Port+"\n")
  fmt.Fprintf(conn,mex+"\n")
}
