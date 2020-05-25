/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-18
 * @Project: Proof of Evolution
 * @Filename: solution.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-25
 * @Copyright: 2020
 */



package ga

import (
  "github.com/D33pBlue/poe/op"
)

// Sol represents a possible solution to the problem
// and an Individual.
type Sol struct{
  Individual DNA
  Fitness float64
  Complex float64
  IsEval bool
  Conf Config
  Gen int
  HashUsed []byte
  JobHash string // used only to share good solutions (to identify the job)
}

// Returns a deep copy instance of a Sol
func (self Sol)DeepCopy()Sol {
  sol := new(Sol)
  sol.Individual = self.Individual.DeepCopy()
  sol.Fitness = self.Fitness
  sol.Complex = self.Complex
  sol.IsEval = self.IsEval
  sol.Conf = self.Conf
  sol.Gen = self.Gen
  sol.HashUsed = self.HashUsed
  sol.JobHash = self.JobHash
  return *sol
}

// Evaluates the individual stored in a Sol, saving its
// Fitness and its Complexity.
func (self *Sol)eval(st *op.State,blockHash []byte){
  st.Reset()
  self.Fitness = self.Individual.Evaluate(st)
  self.Complex = st.NumOperations()
  self.IsEval = true
  self.HashUsed = blockHash
}

// This method is used only for testing; please use eval to evaluate a Sol.
func (self *Sol)Eval2(st *op.State,blockHash []byte){
  st.Reset()
  self.Fitness = self.Individual.Evaluate(st)
  self.Complex = st.NumOperations()
  self.IsEval = true
  self.HashUsed = blockHash
}
