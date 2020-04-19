/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-16
 * @Project: Proof of Evolution
 * @Filename: wallet_test.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-19
 * @Copyright: 2020
 */



package wallet


import(
  "testing"
)

func TestAddress(t *testing.T){
  var addr Addr = "addr"
  if CheckAddr(addr)!=false{
    t.Errorf("Invalid address passed check: %v",addr)
  }
}
