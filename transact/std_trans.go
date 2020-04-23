/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-23
 * @Project: Proof of Evolution
 * @Filename: std_trans.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-24
 * @Copyright: 2020
 */

package transact

import(
  "github.com/D33pBlue/poe/utils"
)

type StdTransaction struct{
  Inputs []TrInput
  Outputs []TrOutput
  Creator utils.Addr
  Hash []byte
  Signature []byte
}

func MakeStdTransaction(creator []byte,key utils.Key)

func (self *StdTransaction)Check(chain *blockchain.Blockchain)bool{

  return true
}

func (self *StdTransaction)IsSpent()bool{
  return false // TODO: implement later
}

func (self *StdTransaction)GetHash()[]byte{
  hb := new(utils.HashBuilder)
  hb.Add(self.Creator)
  // hb.Add(self.Signature)
  for i:=0;i<len(self.Inputs);i++{
    hb.Add(self.Inputs[i].Block)
    hb.Add(self.Inputs[i].ToSpend)
    hb.Add(self.Inputs[i].Index)
  }
  for i:=0;i<len(self.Outputs);i++{
    hb.Add(self.Outputs[i].Address)
    hb.Add(self.Outputs[i].Value)
  }
  return hb.GetHash()
}
