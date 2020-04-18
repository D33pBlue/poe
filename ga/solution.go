package ga

import (
  "github.com/D33pBlue/poe/op"
)

type Sol struct{
  Individual DNA
  Fitness float64
  Complex float64
  IsEval bool
  Conf Config
  Gen int
}

func (self Sol)DeepCopy()Sol {
  sol := new(Sol)
  sol.Individual = self.Individual.DeepCopy()
  sol.Fitness = self.Fitness
  sol.Complex = self.Complex
  sol.IsEval = self.IsEval
  sol.Conf = self.Conf
  sol.Gen = self.Gen
  return *sol
}

func (self *Sol)eval(st *op.State){
  st.Reset()
  self.Fitness = self.Individual.Evaluate(st)
  self.Complex = st.NumOperations()
  self.IsEval = true
}
