/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-18
 * @Project: Proof of Evolution
 * @Filename: transaction.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-19
 * @Copyright: 2020
 */



package block

import(
  "github.com/D33pBlue/poe/coin"
)

// Transact is the interfacte for transactions.
// Transactions can be standard trans., job trans.,
// solution trans., sharing trans.
type Transact interface{
  Check()bool
}

type StdTransact struct{
  From Addr
  To Addr
  Amount coin.Currency

}

type JobTransac struct{

}

type SolTransact struct{

}
