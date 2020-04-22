/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-18
 * @Project: Proof of Evolution
 * @Filename: transaction.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-22
 * @Copyright: 2020
 */



package block

import(
 "github.com/D33pBlue/poe/coin"
 "crypto/sha256"
 "fmt"
)

// Transact is the interfacte for transactions.
// Transactions can be standard trans., job trans.,
// solution trans., sharing trans.
type Transact interface{
  Check()bool
  GetHash()([]byte)
}

type StdTransact struct{
 // From Addr
 // To Addr
 Amount coin.Currency
}

type JobTransac struct{

}

type SolTransact struct{

}

func (self *StdTransact)Check()bool{
  return self.Amount.Check()
}

func (self *StdTransact)GetHash()([]byte)  {
  input1 := []byte(fmt.Sprintf("%v", *self))
  // fmt.Printf("%x\n",input1)
  encoder := sha256.New()
  encoder.Write([]byte(input1))
  return encoder.Sum(nil)
}

func MakeStdTransaction(v coin.Any)Transact{
  trans := new(StdTransact)
  trans.Amount = coin.New(v)
  return trans
}
