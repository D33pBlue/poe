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

// import (
//   "crypto/rsa"
//   "strings"
//   "github.com/D33pBlue/poe/transact"
// )
//
// type Wallet struct{
//   Address string
//   Public *rsa.PublicKey
//   private *rsa.PrivateKey
//   // Balance coin.Currency
//   Unspent []Transact
// }
//
//
// func New() *Wallet{
//   wallet := new(Wallet)
//   var err error
//   wallet.private,err = GenerateKey()
//   if err!=nil{return nil}
//   
//   wallet.Balance = coin.New(0)
//   return wallet
// }
//
// func (self *Wallet)Save(path string){
//
// }
//
// func (self *Wallet)Load(path string){
//
// }
