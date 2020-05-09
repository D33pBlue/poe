/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-16
 * @Project: Proof of Evolution
 * @Filename: wallet.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-09
 * @Copyright: 2020
 */



// Package wallet contains the definitions of the
// wallets, with miner's informations.
package wallet

import (
  "fmt"
  "sort"
  "net"
  "io/ioutil"
  "bufio"
  "errors"
  "strconv"
  "github.com/D33pBlue/poe/utils"
  . "github.com/D33pBlue/poe/blockchain"
)

// Used to keep track of the available transactions in Wallet.
type WalletEntry struct{
  Value int
  Spendable TrInput
}

// A Wallet collects the spendable transactions and can
// interoperate with a miner. It can:
// - ask the total amount of money available
// - create a standart transaction to send money
// - create a job transaction to submit a job
// - ask for the best found solutions of a job
type Wallet struct{
  Id utils.Addr
  Key utils.Key
  MinerIp string
  MinerTrusted bool
  Entries []WalletEntry // keep track of spendable transactions
  SeenBlocks []string // keep hashes of checked blocks
}

// Initializes a Wallet. If a non-empty path is given the public key
// is loaded from file, otherwise a new public key is generated.
func New(path,ip string,trust bool)*Wallet{
  wallet := new(Wallet)
  wallet.MinerIp = ip
  wallet.MinerTrusted = trust
  var err error
  if path==""{
    wallet.Key,err = generateKey()
  }else{
    wallet.Key,err = loadKey(path)
  }
  if err!=nil{
    fmt.Println(err)
    return nil
  }
  wallet.Id = utils.GetAddr(wallet.Key)
  return wallet
}

// Updates the wallet asking to a (trusted?) miner the new blocks
// of the blockchain up to the first seen block.
func (self *Wallet)Update()error {
  conn, err := net.Dial("tcp",self.MinerIp)
  if err!=nil{ return err }
  fmt.Fprintf(conn,"update_wallet\n")
  reader := bufio.NewReader(conn)
  data, err2 := reader.ReadString('\n')
  fmt.Printf("Received: %v",data)
  if err2!=nil{ return err2 }
  head,prev := MarshalBlock([]byte(data))
  if head==nil{
    return errors.New("Unable to marshal block")
  }
  if len(self.SeenBlocks)>0 && head.GetHash(prev)==self.SeenBlocks[len(self.SeenBlocks)-1]{
    // The blockchain is already seen => nothing to update
    return nil
  }
  // There are new blocks => ask them, up to the first known one
  hash := prev
  block := head
  var unseenBlocks []*Block
  unseenBlocks = append(unseenBlocks,head)
  for {
    if len(hash)<=0{break}
    if self.hasSeenBlock(hash)>=0{break}
    block.Previous,hash = self.askBlock(hash,self.MinerIp)
    block = block.Previous
    unseenBlocks = append([]*Block{block},unseenBlocks...)
  }
  if !self.MinerTrusted{
    err := self.checkReceivedBlocks(head,block,hash)
    if err!=nil{ return err }
  }
  // Update self.Entries and self.SeenBlocks
  for i:=0;i<len(unseenBlocks);i++{
    self.SeenBlocks = append(self.SeenBlocks,unseenBlocks[i].GetHashCached())
    transactions := unseenBlocks[i].Transactions.GetTransactionArray()
    for j:=0;j<len(transactions);j++{
      switch transactions[j].GetType() {
      case TrStd:
        stdTr := transactions[j].(*StdTransaction)
        if stdTr.GetCreator()==self.Id{
          for k:=0;k<len(stdTr.Inputs);k++{
            self.removeSpentTransaction(
              stdTr.Inputs[k].Block,
              stdTr.Inputs[k].ToSpend,
              stdTr.Inputs[k].Index)
          }
        }
        for k:=0;k<len(stdTr.Outputs);k++{
          if stdTr.Outputs[k].Address==self.Id{
            self.addSpendableTransaction(
              unseenBlocks[i].GetHashCached(),
              stdTr.GetHashCached(),
              k,stdTr.Outputs[k].Value)
          }
        }
      case TrCoin:
        coinTr := transactions[j].(*CoinTransaction)
        if coinTr.Output.Address==self.Id{
          self.addSpendableTransaction(
            unseenBlocks[i].GetHashCached(),
            coinTr.GetHashCached(),
            0,coinTr.Output.Value)
        }
      // case TrJob:
        // ...
      }
    }
  }
  return nil
}


// Returns the total amount of money earned, summing up all the
// incoming unspent transactions
func (self *Wallet)GetTotal()string{
  err := self.Update()
  if err!=nil{
    fmt.Println(err)
    return ""
  }
  var total int = 0
  for i:=0;i<len(self.Entries);i++{
    total += self.Entries[i].Value
  }
  return strconv.Itoa(total)
}

// SendMoney creates a StdTransaction, after checking
// if the balance of the wallet is big enough, and sends
// it to the miner.
func (self *Wallet)SendMoney(amount int,receiver utils.Addr)error{
  total,err := strconv.Atoi(self.GetTotal())
  if err!=nil{return err}
  if total<amount{
    return errors.New("You does not have enough money")
  }
  // sort the available Entries by value
  sort.SliceStable(self.Entries, func(i, j int) bool {
    return self.Entries[i].Value < self.Entries[j].Value
  })
  var inputs []TrInput
  var outputs []TrOutput
  var spending int = 0
  // collects the entries to use up to the amount
  for i:=0;spending<amount;i++{
    inputs = append(inputs,self.Entries[i].Spendable)
    spending += self.Entries[i].Value
  }
  out := new(TrOutput)
  out.Address = receiver
  out.Value = amount
  outputs = append(outputs,*out)
  // send back the remainder
  if spending-amount>0{
    remainder := new(TrOutput)
    remainder.Address = self.Id
    remainder.Value = spending-amount
    outputs = append(outputs,*remainder)
  }
  // create the transaction
  newTransact := MakeStdTransaction(
    self.Id,
    self.Key,
    inputs,
    outputs)
  if !utils.CheckSignature(newTransact.Signature,newTransact.Hash,newTransact.Creator){
    return errors.New("Invalid Signature")
  }
  if newTransact==nil{
    return errors.New("Error in StdTransaction creation")
  }
  // send the transaction to the miner
  conn, err := net.Dial("tcp",self.MinerIp)
  if err!=nil{ return err }
  fmt.Fprintf(conn,"transaction\n")
  fmt.Fprintf(conn,TrStd+"\n")
  fmt.Fprintf(conn,string(newTransact.Serialize())+"\n")
  fmt.Println("Sent transaction request to miner")
  return nil
}

// SubmitJob creates a JobTransaction and
// sends it to the miner.
func (self *Wallet)SubmitJob(job string)error{
  // TODO: implement later
  return nil
}

// Generates a new couple of public and private keys, and stores
// them in "./data/key<n>.pem" and "./data/key<n>.pem.priv".
func generateKey()(utils.Key,error){
  key,err := utils.GenerateKey()
  if err!=nil{ return nil,err }
  i := 0
  for ;utils.FileExists("data/key"+strconv.Itoa(i)+".pem");i++{}
  name := "data/key"+strconv.Itoa(i)+".pem"
  err = ioutil.WriteFile(name,[]byte(utils.ExportPublicKeyAsPemStr(key)), 0644)
  if err!=nil{ return nil,err }
  err = ioutil.WriteFile(name+".priv",[]byte(utils.ExportPrivateKeyAsPemStr(key)), 0644)
  if err!=nil{
    fmt.Println("Keys stored in "+name+" and "+name+".priv\n")
  }
  return key,err
}

// Loads from file a couple of public and private keys.
func loadKey(path string)(utils.Key,error){
  data, err := ioutil.ReadFile(path)
  if err!=nil{ return nil,err }
  pub,err2 := utils.LoadPublicKeyFromPemStr(data)
  if err2!=nil { return nil,err2 }
  data, err = ioutil.ReadFile(path+".priv")
  if err!=nil { return nil,err }
  priv,err3 := utils.LoadPrivateKeyFromPemStr(data)
  if err3!=nil { return nil,err3 }
  priv.PublicKey = *pub
  return priv,nil
}

// Returns -1 if the hash of the block is new; otherwise
// its index in the seen blockchain.
func (self *Wallet)hasSeenBlock(hash string)int{
  for i:=len(self.SeenBlocks)-1;i>=0;i--{
    if self.SeenBlocks[i]==hash{
      return i
    }
  }
  return -1
}

// Ask to a miner a block with a specific hash and returns
// it and the hash of the previous one; or nil in case of error.
func (self *Wallet)askBlock(blockHash string,ipaddress string)(*Block,string){
  fmt.Println("asking for ",blockHash," to ",ipaddress)
  conn, err := net.Dial("tcp",ipaddress)
  if err!=nil{
    fmt.Println(err)
    return nil,"" }
  fmt.Fprintf(conn,"get_block\n"+blockHash+"\n")
  blockRaw,err2 := bufio.NewReader(conn).ReadString('\n')
  if err2!=nil{
    fmt.Println(err2)
    return nil,"" }
  return MarshalBlock([]byte(blockRaw))
}

// Checks the received unseen blocks
func (self *Wallet)checkReceivedBlocks(head,lastBlock *Block,lastHash string)error{
  if lastBlock!=nil{
    seenChainLength := self.hasSeenBlock(lastHash)
    if lastBlock.LenSubChain!=seenChainLength{
      return errors.New("Chain length mismatch")
    }
  }
  //
  // TODO: complete check
  //
  //
  return nil
}


// Removes from Wallet.Entries an entry.
func (self *Wallet)removeSpentTransaction(block,transact string,index int){
  var updatedSpedable []WalletEntry
  for i:=0;i<len(self.Entries);i++{
    if !(self.Entries[i].Spendable.Block==block &&
        self.Entries[i].Spendable.ToSpend==transact &&
        self.Entries[i].Spendable.Index==index){
      updatedSpedable = append(updatedSpedable,self.Entries[i])
    }
  }
  self.Entries = updatedSpedable
}

// Add an entry to Wallet.Entries.
func (self *Wallet)addSpendableTransaction(block,transact string,index,value int){
  entry := new(WalletEntry)
  entry.Value = value
  entry.Spendable.Block = block
  entry.Spendable.ToSpend = transact
  entry.Spendable.Index = index
  self.Entries = append(self.Entries,*entry)
}
