/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-27
 * @Project: Proof of Evolution
 * @Filename: sha256_test.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-27
 * @Copyright: 2020
 */

package utils

import(
  "fmt"
  "testing"
)

func TestHash(t *testing.T)  {
  nonce := 0
  for {
    hb := new(HashBuilder)
    hb.Add(nonce)
    if hb.GetHash()[0]==0{
      fmt.Println(hb.GetHash()[0])
      break
    }
    nonce ++
  }
  hb := new(HashBuilder)
  hb.Add(nonce)
  fmt.Println("Hash:")
  fmt.Println(hb.GetHash())
  fmt.Println("nonce:",nonce)
}
