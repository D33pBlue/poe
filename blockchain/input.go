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


type TrInput struct{
  Block []byte // hash of the old block with the transaction to spend
  ToSpend []byte // hash of the old transaction to spend
  Index int // index inside ToSpend of the record to spend
  // proof of ownership?
}
