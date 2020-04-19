/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-18
 * @Project: Proof of Evolution
 * @Filename: dna.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-19
 * @Copyright: 2020
 */



package ga

import(
  "math/rand"
  "github.com/D33pBlue/poe/op"
)

// DNA is the interface the user has to implement
// to define a problem for a Job transaction.
type DNA interface {
  Generate(prng *rand.Rand) DNA
  Mutate(prng *rand.Rand) DNA
  Crossover(ind2 DNA,prng *rand.Rand) DNA
  Evaluate(st *op.State) float64
  DeepCopy() DNA
  HasToMinimize() bool
}
