package wallet

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
