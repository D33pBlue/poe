/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-18
 * @Project: Proof of Evolution
 * @Filename: miner.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-16
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
  "errors"
  "github.com/D33pBlue/poe/utils"
  "github.com/D33pBlue/poe/conf"
  "github.com/D33pBlue/poe/blockchain"
)

// The Miner struct is used to keep alive a mining node.
// It stores the blockchain and the list of all connected miners.
type Miner struct{
  Chain *blockchain.Blockchain
  Port string
  Connected []string // sinc read/write
  connected_lock sync.Mutex
  keepServing bool
  addrch chan string
  id utils.Addr
  config *conf.Config
}

// Create a new Miner. In order to start mining you have to
// call miner.Serve() method (possibly in a goroutine).
func New(port string,id utils.Addr,config *conf.Config)*Miner{
  miner := new(Miner)
  miner.Port = port
  miner.id = id
  miner.config = config
  var folder string = config.GetChainFolder()//fmt.Sprintf("data/chain%v",port)
  // make dir if not exists
  _, err := os.Stat(folder)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(folder, 0755)
		if errDir != nil {fmt.Println(errDir)}
    miner.Chain = blockchain.NewBlockchain(id,folder,config)
	}else{
    miner.Chain = blockchain.LoadChainFromFolder(id,folder,config)
  }
  miner.keepServing = false
  miner.addrch = make(chan string)
  return miner
}

// This is the main loop that keeps alive a mining node
// and manages its connections calling handleConnection.
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

// Returns the list of the addresses of linked miners
func (self *Miner)GetConnected()[]string{
  self.connected_lock.Lock()
  var conn []string
  for i:=0;i<len(self.Connected);i++{
    conn = append(conn,self.Connected[i])
  }
  self.connected_lock.Unlock()
  return conn
}

// This functions can be used to connect the miner to
// another one and update the blockchain with him.
func (self *Miner)AddNode(ipaddress string)error{
  match, _ := regexp.MatchString("[0-9]+.[0-9]+.[0-9]+.[0-9]+:[0-9]+",ipaddress)
  if !match{ return errors.New(ipaddress+" is not a valid address")}
  err := self.requestUpdate(ipaddress)
  if err!=nil{ return err }
  self.addrch <- ipaddress
  return nil
}

// Keeps track of the address of another miner if it is new.
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


// handleConnection is the main function for connections:
// it reads a command from the socket and acts accordingly,
// managing an open connection untill it is closed.
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
  case "update_wallet":
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
  case "transaction":
    transacType,err := reader.ReadString('\n')
    if err!=nil{
      fmt.Println(err)
    }else{
      data,err2 := reader.ReadString('\n')
      if err2!=nil{
        fmt.Println(err2)
      }else{
        mexTransaction := new(blockchain.MexTrans)
        mexTransaction.Type = transacType[:len(transacType)-1]
        mexTransaction.Data = []byte(data)
        self.Chain.TransQueue <- *mexTransaction
      }
    }
  case "miniblock":
    // send miniblock to self.Chain.MiniBlockIn
    miniblock,err := reader.ReadString('\n')
    if err==nil{
      mexBlock := new(blockchain.MexBlock)
      mexBlock.Data = []byte(miniblock)
      self.Chain.MiniBlockIn <- *mexBlock
    }
  }
}

// Asks to another miner his current blockchain
// and send a MexBlock through self.Chain.BlockIn channel.
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

// Whenever it receives from self.Chain.BlockOut a mined block,
// this method sends the block to all the linked miners.
// Also update the list of linked miners reading new IPs
// from self.addrch channel.
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
      case mex := <-self.Chain.MiniBlockOut:
        self.connected_lock.Lock()
        for i:=0;i<len(self.Connected);i++{
          fmt.Println("Send miniblock update to ",self.Connected[i])
          go self.sendMiniBlockUpdate(self.Connected[i],string(mex.Data))
        }
        self.connected_lock.Unlock()
    }
  }
}


// Sends a message and wait for ack.
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

// Sends a block to a miner.
func (self *Miner)sendBlockUpdate(address string,mex string){
  conn, err := net.Dial("tcp",address)
  if err!=nil{ return }
  fmt.Fprintf(conn,"chain\n")
  fmt.Fprintf(conn,self.Port+"\n")
  fmt.Fprintf(conn,mex+"\n")
}

// Sends a miniblock to a miner.
func (self *Miner)sendMiniBlockUpdate(address string,mex string){
  conn, err := net.Dial("tcp",address)
  if err!=nil{ return }
  fmt.Fprintf(conn,"miniblock\n")
  fmt.Fprintf(conn,mex+"\n")
}
