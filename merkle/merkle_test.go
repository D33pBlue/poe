/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:  mer 22 apr 2020 22:46:43 CEST
 * @Project: Proof of Evolution
 * @Filename: block.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-23
 * @Copyright: 2020
 */

 package merkle

 import(
 	"testing"
  "fmt"
  "github.com/D33pBlue/poe/block"
 )

 func TestMerkle(t *testing.T){
	m := New()
  fmt.Println(m.Root)
  for i := 0; i < 3; i++ {
    trs := block.MakeStdTransaction(i)
    fmt.Printf("%v: %x\n",i,trs.GetHash())
    m.Add(trs)
  }
  fmt.Printf("root: %x\n",m.Root.Hash)
  fmt.Printf("root: %x\n",m.Root.L.Hash)
  fmt.Printf("root: %x\n",m.Root.R.Hash)
  if !m.Check(){ t.Errorf("Build invalid tree")}
 }
