/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-23
 * @Project: Proof of Evolution
 * @Filename: utils.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-18
 * @Copyright: 2020
 */

package utils

import (
  "fmt"
  "crypto/sha256"
)

// HashBuilder can be used to collect data and obtain the
// hash of the total bytes.
type HashBuilder struct{
  data []byte
}

// Adds some data to the HashBuilder. The data can be of any type
// and is read as []byte.
func (self *HashBuilder)Add(data interface{}){
  var binary []byte = []byte(fmt.Sprintf("%v",data))
  self.data = append(self.data,binary...)
}

// Returns the hash of all the collected data of an HashBuilder.
func (self *HashBuilder)GetHash()([]byte){
  encoder := sha256.New()
  encoder.Write(self.data)
  return encoder.Sum(nil)
}


func CompareHashes(h1,h2 string)bool{
  return h1==h2
}

func CompareSlices(h1,h2 []byte)bool{
  if len(h1)!=len(h2){ return false }
  for i:=0;i<len(h1);i++{
    if h1[i]!=h2[i]{ return false }
  }
  return true
}
