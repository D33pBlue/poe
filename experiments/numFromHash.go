/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-18
 * @Project: Proof of Evolution
 * @Filename: numFromHash.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-19
 * @Copyright: 2020
 */



// To test hash generation and use

package main

import(
  "fmt"
   "encoding/binary"
   "crypto/sha256"
   // "encoding"
)

func main(){
  const (input1 = "golang test")
	first := sha256.New()
	first.Write([]byte(input1))
	var hash []byte = first.Sum(nil)
  fmt.Printf("len: %v, value: %x\n",len(hash), hash)
  // var mySlice = []byte{244, 204, 244, 244, 244, 244, 244, 244}
  data := binary.BigEndian.Uint64(hash)
  fmt.Println(data)
}
