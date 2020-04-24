/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-23
 * @Project: Proof of Evolution
 * @Filename: input.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-24
 * @Copyright: 2020
 */


package blockchain

import(
  "github.com/D33pBlue/poe/utils"
)

type TrOutput struct{
  Address utils.Addr // address of the receiver
  Value int // value to exchange
  // proof of ownership?
}
