package ga

import(
  "github.com/D33pBlue/poe/op"
)

type Comp func(float64,float64)bool
func Maximize(s1,s2 float64) bool{
  return s1 > s2
}
func Minimize(s1,s2 float64) bool{
  return s1 < s2
}

var Optimum Comp = Minimize

type Population []Sol

func (a Population) Len() int { return len(a) }
func (a Population) Swap(i, j int){ a[i], a[j] = a[j], a[i] }
func (a Population) Less(i, j int) bool {
  if !a[i].IsEval{
    // a[i].eval()
    panic("Solution not evaluated!")
  }
  if !a[j].IsEval{
    // a[j].eval()
    panic("Solution not evaluated!")
  }
  return Optimum(a[i].Fitness,a[j].Fitness)
}

func (pop Population)eval(blockHash []byte) Sol  {
  st := op.MakeState(blockHash)
  best := pop[0]
  for i:=0;i<pop.Len();i++{
    pop[i].eval(st)
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
