/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-18
 * @Project: Proof of Evolution
 * @Filename: transaction.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-23
 * @Copyright: 2020
 */

package transact

import(
  "fmt"
  "github.com/D33pBlue/poe/blockchain"
)

type Transaction interface{
  Check(chain *blockchain.Blockchain)bool
  IsSpent()bool
  GetHash()[]byte
}

func (self Transaction)GetType()string{
  return fmt.Sprintf("%T",self)
}
