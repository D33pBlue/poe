/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-18
 * @Project: Proof of Evolution
 * @Filename: population.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-19
 * @Copyright: 2020
 */


package ga

import(
  "github.com/D33pBlue/poe/op"
)

type Comp func(float64,float64)bool

// Function to choose for maximization problems.
func Maximize(s1,s2 float64) bool{
  return s1 > s2
}

// Function to choose for minimization problems.
func Minimize(s1,s2 float64) bool{
  return s1 < s2
}

var Optimum Comp = Minimize

// A population is a set of solution candidates, and here is
// represented as array of Sol.
type Population []Sol

func (a Population) Len() int { return len(a) }
func (a Population) Swap(i, j int){ a[i], a[j] = a[j], a[i] }
func (a Population) Less(i, j int) bool {
  // if !a[i].IsEval || !a[j].IsEval{return }
  return Optimum(a[i].Fitness,a[j].Fitness)
}

// Evals the solution candidates of the population and
// sends each evaluated one to the chNonce channel in input.
func (pop Population)eval(blockHash []byte,chNonce chan Sol) Sol  {
  st := op.MakeState(blockHash)
  best := pop[0]
  for i:=0;i<pop.Len();i++{
    pop[i].eval(st,blockHash)
    select{
    case chNonce <- pop[i]:
    default:
    }
    if Optimum(pop[i].Fitness,best.Fitness){
      best = pop[i]
    }
  }
  return best
}

// Clears the fitness of the solutions in the population.
func (self Population)reset(){
  for i:=0; i<self.Len(); i++{
    self[i].IsEval = false
  }
}

// Returns a deep copy instance of the whole population.
func (self Population)DeepCopy()(pop Population){
  for i:=0;i<len(self);i++{
    pop = append(pop,self[i].DeepCopy())
  }
  return
}
