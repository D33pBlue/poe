package block

type Block struct{
  Transactions []Transact
  Nonce Nonce
}

func (self *Block)AddTransaction(t Transact){
  if t.Check() {
    self.Transactions = append(self.Transactions,t)
  }
}

func (self *Block)GetHash(withoutNonce bool)[]byte{
  return nil
}
