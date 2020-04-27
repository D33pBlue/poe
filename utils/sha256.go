/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-23
 * @Project: Proof of Evolution
 * @Filename: utils.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-27
 * @Copyright: 2020
 */

package utils

import (
  "fmt"
  "crypto/sha256"
)

type HashBuilder struct{
  data []byte
}

func (self *HashBuilder)Add(data interface{}){
  var binary []byte = []byte(fmt.Sprintf("%v",data))
  self.data = append(self.data,binary...)
}

func (self *HashBuilder)GetHash()([]byte){
  encoder := sha256.New()
  encoder.Write(self.data)
  return encoder.Sum(nil)
}

func CompareHashes(h1,h2 []byte)bool{
  if len(h1)!=len(h2){ return false }
  for i:=0;i<len(h1);i++{
    if h1[i]!=h2[i]{ return false }
  }
  return true
}
