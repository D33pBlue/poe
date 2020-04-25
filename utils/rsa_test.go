/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-25
 * @Project: Proof of Evolution
 * @Filename: rsa_test.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-25
 * @Copyright: 2020
 */

package utils

import(
  "testing"
)

func TestLoad(t *testing.T)  {
  key,_ := GenerateKey()
  addr := GetAddr(key)
  pubk := publicKeyFromAddr(addr)
  if *pubk==key.PublicKey{
    t.Errorf("The key from address is different:\n%v\n\t!=\n%v\n",*pubk,key.PublicKey)
  }
}
