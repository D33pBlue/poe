/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-16
 * @Project: Proof of Evolution
 * @Filename: wallet.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-08
 * @Copyright: 2020
 */



// Package wallet contains the definitions of the
// wallets, with miner's informations.
package wallet

import (
  // "os"
  "fmt"
  // "sort"
  "net"
  "io/ioutil"
  "bufio"
  "errors"
  "strconv"
  "github.com/D33pBlue/poe/utils"
  "github.com/D33pBlue/poe/blockchain"
)

type WalletEntry struct{
  Value int
  Spendable TrInput
}

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
  if err2!=nil{ return err2 }
  head,prev := MarshalBlock([]byte(data))
  if head==nil{
    return errors.New("Unable to marshal block")
  }
  if head.GetHash()==self.SeenBlocks[len(self.SeenBlocks)-1]{
    // The blockchain is already seen => nothing to update
    return nil
  }
  // There are new blocks => ask them, up to the first known one
  hash := prev
  block := head
  for {
    if len(hash)<=0{break}
    if self.hasSeenBlock(hash)>=0{break}
    block.Previous,hash = self.askBlock(hash,self.MinerIp)
    block = block.Previous
  }
  if !self.MinerTrusted{
    // Check the received blocks
    if block{
      seenChainLength := self.hasSeenBlock(hash)
      if block.LenSubChain!=seenChainLength{
        return errors.New("Chain length mismatch")
      }
    }
    // TODO: complete the check of the blocks
  }
  // Update self.Entries and self.SeenBlocks
  // TODO: ...
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

func (self *Wallet)SendMoney(amount int,receiver utils.Addr)error{
  // TODO: implement with new architecture

  // total,err := strconv.Atoi(self.GetTotal())
  // if err!=nil{return err}
  // if total<amount{
  //   return errors.New("You does not have enough money")
  // }
  // transactions := self.getActiveTransactions()
  // sort.SliceStable(transactions, func(i, j int) bool {
  //   return transactions[i].GetSpendingValueFor(self.Id) < transactions[j].GetSpendingValueFor(self.Id)
  // })
  // var transactToSpend []blockchain.Transaction
  // var spending int = 0
  // for i:=0;spending<amount;i++{
  //   transactToSpend = append(transactToSpend,transactions[i])
  //   spending += transactions[i].GetSpendingValueFor(self.Id)
  // }
  // // newTransact,err2 := blockchain.MakeStdTransaction()
  // // TODO: make transaction
  return nil
}

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
