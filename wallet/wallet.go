/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-16
 * @Project: Proof of Evolution
 * @Filename: wallet.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-25
 * @Copyright: 2020
 */



// Package wallet contains the definitions of the
// wallets, with miner's informations.
package wallet

import (
  "github.com/D33pBlue/poe/utils"
)

type Wallet struct{
  Id utils.Addr
  Key utils.Key
  MinerIp string
}

func (self *Wallet)GetTotal()int  {
  // TODO: implement later
  return 0
}

func (self *Wallet)SendMoney(amount int,receiver utils.Addr){
  // TODO: implement later
}

func (self *Wallet)SubmitJob(job string){
  // TODO: implement later
}
