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
