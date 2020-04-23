/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-16
 * @Project: Proof of Evolution
 * @Filename: wallet.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-23
 * @Copyright: 2020
 */



// Package wallet contains the definitions of the
// wallets, with miner's informations.
package wallet

import (
  "crypto/rsa"
  "strings"
  "github.com/D33pBlue/poe/coin"
)

type Wallet struct{
  Address string
  Public *rsa.PublicKey
  private *rsa.PrivateKey
  Balance coin.Currency
}

func FromAddress(addr string,balance coin.Currency)*Wallet{
  wallet := new(Wallet)
  wallet.Address = addr
  wallet.Balance = balance
  wallet.private = nil
  // wallet.Public = ..
  return wallet
}

func New() *Wallet{
  wallet := new(Wallet)
  var err error
  wallet.private,err = GenerateKey()
  if err!=nil{return nil}
  wallet.Public = &wallet.private.PublicKey
  lines := strings.Split(ExportPublicKeyAsPemStr(wallet.Public),"\n")
  wallet.Address = strings.Join(lines[1:len(lines)-2],"")
  wallet.Balance = coin.New(0)
  return wallet
}

func (self *Wallet)Save(path string){

}

func (self *Wallet)Load(path string){

}
