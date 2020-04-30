/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-23
 * @Project: Proof of Evolution
 * @Filename: rsa.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-30
 * @Copyright: 2020
 */

package utils

import (
    "crypto"
    "crypto/rand"
    "crypto/rsa"
    "fmt"
    "errors"
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

func ExportPublicKeyAsPemStr2(pubkey *rsa.PublicKey) string {
  pubkey_pem := string(pem.EncodeToMemory(&pem.Block{Type:  "RSA PUBLIC KEY",Bytes: x509.MarshalPKCS1PublicKey(pubkey)}))
  return pubkey_pem
}

func ExportPrivateKeyAsPemStr(key Key) string {
  privatekey_pem := string(pem.EncodeToMemory(&pem.Block{Type:  "RSA PRIVATE KEY",Bytes: x509.MarshalPKCS1PrivateKey(key)}))
  return privatekey_pem
}

func LoadPrivateKeyFromPemStr(pemStr []byte)(Key,error) {
  privPem, _ := pem.Decode(pemStr)
	if privPem.Type != "RSA PRIVATE KEY" {
    return nil,errors.New("RSA private key is of the wrong type")
  }
  key, err := x509.ParsePKCS1PrivateKey(privPem.Bytes)
	if err != nil {return nil,err}
	return key,nil
}

func LoadPublicKeyFromPemStr(pemStr []byte)(*rsa.PublicKey,error) {
  pubPem, _ := pem.Decode(pemStr)
	if pubPem.Type != "RSA PUBLIC KEY" {
    return nil,errors.New("RSA public key is of the wrong type")
  }
  key, err := x509.ParsePKCS1PublicKey(pubPem.Bytes)
	if err != nil {return nil,err}
	return key,nil
}

func GetAddr(key Key)Addr{
  lines := strings.Split(ExportPublicKeyAsPemStr(key),"\n")
  return Addr(strings.Join(lines[1:len(lines)-2],""))
}

func GetAddr2(key *rsa.PublicKey)Addr{
  lines := strings.Split(ExportPublicKeyAsPemStr2(key),"\n")
  return Addr(strings.Join(lines[1:len(lines)-2],""))
}

func publicKeyFromAddr(addr Addr)*rsa.PublicKey{
  var pemStr string = "-----BEGIN RSA PUBLIC KEY-----\n"
  pemStr += string(addr)+"\n"
  pemStr += "-----END RSA PUBLIC KEY-----"
  pub,err := LoadPublicKeyFromPemStr([]byte(pemStr))
  if err!=nil{fmt.Println(err)}
  return pub
}

type SignBuilder struct{
  data []byte
}

func (self *SignBuilder)Add(data interface{}){
  var binary []byte = []byte(fmt.Sprintf("%v",data))
  self.data = append(self.data,binary...)
}

func GetSignatureFromHash(hashed string,key Key)[]byte{
  var opts rsa.PSSOptions
  opts.SaltLength = rsa.PSSSaltLengthAuto // for simple example
  newhash := crypto.SHA256
  signature,err := rsa.SignPSS(rand.Reader,key,newhash,[]byte(hashed),&opts)
  if err != nil {
  		fmt.Println(err)
      return nil
  	}
  return signature
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

func CheckSignature(sign,hashed string,addr Addr)bool{
  var opts rsa.PSSOptions
  opts.SaltLength = rsa.PSSSaltLengthAuto
  newhash := crypto.SHA256
  publicKey := publicKeyFromAddr(addr)
  err := rsa.VerifyPSS(publicKey,newhash,[]byte(hashed),[]byte(sign),&opts)
  if err != nil {
    return false
  }
  return true
}
