/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-16
 * @Project: Proof of Evolution
 * @Filename: address.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-19
 * @Copyright: 2020
 */



package wallet

// Addr represents an address, that is the identifier
// of a wallet.
type Addr string

func NewAddr() Addr{
  return "eijwkjcjwhblchbqliwhbchiqbhv"
}

func CheckAddr(addr Addr)bool{
  if len(addr)==64{
    return true
  }
  return false
}
