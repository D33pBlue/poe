/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-23
 * @Project: Proof of Evolution
 * @Filename: rsa.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-24
 * @Copyright: 2020
 */

package utils

import (
    "crypto"
    "crypto/rand"
    "crypto/rsa"
    "fmt"
    "crypto/x509"
    "encoding/pem"
    "strings"
)

type Key *rsa.PrivateKey
type Addr string

func GenerateKey()(Key,error) {
  reader := rand.Reader
	return rsa.GenerateKey(reader,2048)
}

func ExportPublicKeyAsPemStr(key Key) string {
  pubkey := &key.PublicKey
  pubkey_pem := string(pem.EncodeToMemory(&pem.Block{Type:  "RSA PUBLIC KEY",Bytes: x509.MarshalPKCS1PublicKey(pubkey)}))
  return pubkey_pem
}

func ExportPrivateKeyAsPemStr(key Key) string {
  privatekey_pem := string(pem.EncodeToMemory(&pem.Block{Type:  "RSA PRIVATE KEY",Bytes: x509.MarshalPKCS1PrivateKey(key)}))
  return privatekey_pem
}

func GetAddr(key Key)Addr{
  lines := strings.Split(ExportPublicKeyAsPemStr(key),"\n")
  return Addr(strings.Join(lines[1:len(lines)-2],""))
}

func publicKeyFromAddr(addr Addr)*rsa.PublicKey{
  return nil // TODO: implement later
}

type SignBuilder struct{
  data []byte
}

func (self *SignBuilder)Add(data interface{}){
  var binary []byte = []byte(fmt.Sprintf("%v",data))
  self.data = append(self.data,binary...)
}

func (self *SignBuilder)GetSignature(key Key)[]byte{
  var opts rsa.PSSOptions
  opts.SaltLength = rsa.PSSSaltLengthAuto // for simple example
  newhash := crypto.SHA256
  pssh := newhash.New()
  pssh.Write(self.data)
  hashed := pssh.Sum(nil)
  signature,err := rsa.SignPSS(rand.Reader,key,newhash,hashed,&opts)
  if err != nil {
  		fmt.Println(err)
      return nil
  	}
  return signature
}

func CheckSignature(sign,hashed []byte,addr Addr)bool{
  var opts rsa.PSSOptions
  opts.SaltLength = rsa.PSSSaltLengthAuto
  newhash := crypto.SHA256
  publicKey := publicKeyFromAddr(addr)
  err := rsa.VerifyPSS(publicKey,newhash,hashed,sign,&opts)
  if err != nil {
    return false
  }
  return true
}
