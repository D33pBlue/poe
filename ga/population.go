/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-18
 * @Project: Proof of Evolution
 * @Filename: population.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-18
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

type Population []Sol

func (a Population) Len() int { return len(a) }
func (a Population) Swap(i, j int){ a[i], a[j] = a[j], a[i] }
func (a Population) Less(i, j int) bool {
  // if !a[i].IsEval || !a[j].IsEval{return }
  return Optimum(a[i].Fitness,a[j].Fitness)
}

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

func (self Population)reset(){
  for i:=0; i<self.Len(); i++{
    self[i].IsEval = false
  }
}

func (self Population)DeepCopy()(pop Population){
  for i:=0;i<len(self);i++{
    pop = append(pop,self[i].DeepCopy())
  }
  return
}
