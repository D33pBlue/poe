/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:  mer 22 apr 2020 22:46:43 CEST
 * @Project: Proof of Evolution
 * @Filename: block.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-27
 * @Copyright: 2020
 */

package blockchain

 import(
 	"testing"
 //  "fmt"
  "github.com/D33pBlue/poe/utils"
 )


 func CountLeaves(n *Node)int{
 	if n==nil{ return 0}
 	if n.L==nil && n.R==nil{
 		return 1
 	}
 	return CountLeaves(n.L)+CountLeaves(n.R)
 }

 func CountCoinTr(n *Node)int{
 	if n==nil{return 0}
 	if n.Type==TrCoin{
 		return 1+CountCoinTr(n.L)+CountCoinTr(n.R)
 	}
 	return CountCoinTr(n.L)+CountCoinTr(n.R)
 }

 func (self *Tree)CountCoinTr2()int{
 	tot := 0
 	for i:=0;i<len(self.transactions);i++{
 		if self.transactions[i].GetType()==TrCoin{
 			tot += 1
 		}
 	}
 	return tot
 }

 func TestMerkleConstruction(t *testing.T){
	m := BuildMerkleTree()
  if m==nil{t.Errorf("Merkle tree not built")}
  if m.Root!=nil{t.Errorf("Bad root initialization")}
  if m.Nleaves!=0{t.Errorf("Bad Nleaves initialization")}
 }

 func TestAddTransaction(t *testing.T){
   m := BuildMerkleTree()
   for i:=0;i<1000;i++{
     if m.CountCoinTr2()!=CountCoinTr(m.Root){
       t.Errorf("Error before adding transaction %v\n",i)
       break
     }
     kk := m.CountCoinTr2()
     tr,_ := MakeCoinTransaction(utils.Addr("adccdcda"),10)
     m.Add(tr)
     if m.CountCoinTr2()!=CountCoinTr(m.Root){
       t.Errorf("Error after adding transaction %v\n",i)
       break
     }
     if m.CountCoinTr2()!=kk+1{
       t.Errorf("Error in number of CoinTransaction\n")
       break
     }
     if CountLeaves(m.Root)!=len(m.transactions){
       t.Errorf("Invalid num of leaves %v %v\n",CountLeaves(m.Root),len(m.transactions))
     }
   }
 }
