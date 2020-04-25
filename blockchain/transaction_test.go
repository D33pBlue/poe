/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-24
 * @Project: Proof of Evolution
 * @Filename: transaction_test.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-25
 * @Copyright: 2020
 */

package blockchain

import(
  "testing"
)

func TestStdTransactionBuild(t *testing.T)  {
  // tr,err := MakeStdTransaction("addr",nil,nil,nil)
  // if err!=nil{t.Errorf("Error: %v\n",err)}
  // if tr.GetType()!=TrStd{
  //   t.Errorf("Create transaction of invalid type %v!=%v\n",tr.GetType(),TrStd)
  // }
}

func TestCoinTransactionBuild(t *testing.T)  {
  // tr,err := MakeCoinTransaction(nil,nil)
  // if err!=nil{t.Errorf("Error: %v\n",err)}
  // if tr.GetType()!=TrCoin{
  //   t.Errorf("Create transaction of invalid type %v!=%v\n",tr.GetType(),TrCoin)
  // }
}
