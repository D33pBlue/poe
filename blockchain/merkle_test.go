/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:  mer 22 apr 2020 22:46:43 CEST
 * @Project: Proof of Evolution
 * @Filename: block.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-24
 * @Copyright: 2020
 */

package blockchain

 import(
 	"testing"
 //  "fmt"
 //  "github.com/D33pBlue/poe/transact"
 )

 func TestMerkleConstruction(t *testing.T){
	m := New()
  if m==nil{t.Errorf("Merkle tree not built")}
  if m.Root!=nil{t.Errorf("Bad root initialization")}
  if m.Nleaves!=0{t.Errorf("Bad Nleaves initialization")}
 }
