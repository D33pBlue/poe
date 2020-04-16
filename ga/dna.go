package ga

import(
  "math/rand"
  "github.com/D33pBlue/poe/op"
)

type DNA interface {
  Generate(prng *rand.Rand) DNA
  Mutate(prng *rand.Rand) DNA
  Crossover(ind2 DNA,prng *rand.Rand) DNA
  Evaluate(st *op.State) float64
  DeepCopy() DNA
  HasToMinimize() bool
}
