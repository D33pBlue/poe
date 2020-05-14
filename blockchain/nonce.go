/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-18
 * @Project: Proof of Evolution
 * @Filename: nonce.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-14
 * @Copyright: 2020
 */

package blockchain

import(
  "github.com/D33pBlue/poe/ga"
)

type Nonce struct{
  Solution ga.DNA
  Evaluation float64
  Complexity float64
  candidates chan ga.Sol // should be a buffered chan
}

// Reads a new Nonce from candidates and change the stored values.
func (self *Nonce)Next(){
  result := <-self.candidates
  self.Solution = result.Individual
  self.Evaluation = result.Fitness
  self.Complexity = result.Complex
}

type NonceNoJob struct{
  Value int
}

func (self *NonceNoJob)Next(){
  self.Value++
}
