/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-16
 * @Project: Proof of Evolution
 * @Filename: wallet.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-30
 * @Copyright: 2020
 */



// Package wallet contains the definitions of the
// wallets, with miner's informations.
package wallet

import (
  "os"
  "fmt"
  "net"
  "io/ioutil"
  "bufio"
  "strconv"
  "github.com/D33pBlue/poe/utils"
)

type Wallet struct{
  Id utils.Addr
  Key utils.Key
  MinerIp string
}

func New(path,ip string)*Wallet{
  wallet := new(Wallet)
  wallet.MinerIp = ip
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

func (self *Wallet)GetTotal()string  {
  conn, err := net.Dial("tcp",self.MinerIp)
  if err!=nil{
    fmt.Println(err)
    return ""
  }
  fmt.Fprintf(conn,"get_total\n")
  fmt.Fprintf(conn,string(self.Id)+"\n")
  reader := bufio.NewReader(conn)
  total, _ := reader.ReadString('\n')
  return total[:len(total)-1]
}

func (self *Wallet)SendMoney(amount int,receiver utils.Addr)error{
  // TODO: implement later
  return nil
}

func (self *Wallet)SubmitJob(job string)error{
  // TODO: implement later
  return nil
}


func generateKey()(utils.Key,error){
  key,err := utils.GenerateKey()
  if err!=nil{ return nil,err }
  i := 0
  for ;fileExists("data/key"+strconv.Itoa(i)+".pem");i++{}
  name := "data/key"+strconv.Itoa(i)+".pem"
  err = ioutil.WriteFile(name,[]byte(utils.ExportPublicKeyAsPemStr(key)), 0644)
  if err!=nil{ return nil,err }
  err = ioutil.WriteFile(name+".priv",[]byte(utils.ExportPrivateKeyAsPemStr(key)), 0644)
  return key,err
}

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


func fileExists(filename string) bool {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}
