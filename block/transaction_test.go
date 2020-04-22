/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-22
 * @Project: Proof of Evolution
 * @Filename: transaction_test.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-22
 * @Copyright: 2020
 */

package block

import(
  "testing"
  "fmt"
)

func TestTransaction(t *testing.T){
  trs := MakeStdTransaction(111111111111111)
  fmt.Println(trs.Check())
  fmt.Printf("%x\n",trs.GetHash())
}
