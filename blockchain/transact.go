/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-18
 * @Project: Proof of Evolution
 * @Filename: transaction.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-01
 * @Copyright: 2020
 */

package blockchain

import(
  "github.com/D33pBlue/poe/utils"
)

type Transaction interface{
  Check(chain *Blockchain)bool
  IsSpent()bool
  SetSpent()
  GetHash()string
  GetHashCached()string
  GetType()string
  GetCreator()utils.Addr
  // Serialize()[]byte
  // Marshal([]byte)Transaction
}

const (
  TrStd = "StdTransaction"
  TrCoin = "CoinTransaction"
  TrJob = "JobTransaction"
)

func GetTransactionMoneyForWallet(chain *Blockchain,tr Transaction,wallet utils.Addr)int{
  switch tr.GetType() {
  case TrCoin:
    var transact *CoinTransaction = tr.(*CoinTransaction)
    if transact.Output.Address==wallet{
      return transact.Output.Value
    }
  case TrStd:
    tot := 0
    var transact *StdTransaction = tr.(*StdTransaction)
    if transact.Creator==wallet{
      for i:=0;i<len(transact.Inputs);i++{
        tot -= chain.SpendSubTransaction(&transact.Inputs[i],wallet)
      }
    }
    for i:=0;i<len(transact.Outputs);i++{
      if transact.Outputs[i].Address==wallet{
        tot += transact.Outputs[i].Value
      }
    }
    return tot
  // case TrJob:
  //   var transact *JobTransaction = tr.(*JobTransaction)
  }
  return 0
}
