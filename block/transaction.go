package block

import(
  "github.com/D33pBlue/poe/coin"
)

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
