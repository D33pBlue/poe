/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-18
 * @Project: Proof of Evolution
 * @Filename: transaction.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-24
 * @Copyright: 2020
 */

package blockchain


type Transaction interface{
  Check(chain *Blockchain)bool
  IsSpent()bool
  GetHash()[]byte
  GetType()string
}

const (
  TrStd = "StdTransaction"
  TrCoin = "CoinTransaction"
  TrJob = "JobTransaction"
)
