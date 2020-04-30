/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-27
 * @Project: Proof of Evolution
 * @Filename: sha256_test.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-28
 * @Copyright: 2020
 */

package utils

import(
  // "fmt"
  "testing"
)

func TestHash(t *testing.T)  {
  nonce := 0
  var hash []byte
  for {
    hb := new(HashBuilder)
    hb.Add(nonce)
    hash = hb.GetHash()
    if hash[0]==0{
      break
    }
    nonce ++
  }
  if hash[0]!=0{
    t.Errorf("The hash does not start with zero: %v\n",hash)
  }
}
