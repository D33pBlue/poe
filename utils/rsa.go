/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-23
 * @Project: Proof of Evolution
 * @Filename: rsa.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-May-09
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
    "encoding/hex"
    "strings"
)

type Key *rsa.PrivateKey // type for private key
type Addr string // type for public key

// Generates a couple of public and private keys.
// The public key is stored inside; in order to extract it
// you can use GetAddr(key)
func GenerateKey()(Key,error) {
  reader := rand.Reader
	return rsa.GenerateKey(reader,2048)
}

// Returns the public key as string in pem format.
func ExportPublicKeyAsPemStr(key Key) string {
  pubkey := &key.PublicKey
  pubkey_pem := string(pem.EncodeToMemory(&pem.Block{Type:  "RSA PUBLIC KEY",Bytes: x509.MarshalPKCS1PublicKey(pubkey)}))
  return pubkey_pem
}

// Returns the public key as string in pem format.
func ExportPublicKeyAsPemStr2(pubkey *rsa.PublicKey) string {
  pubkey_pem := string(pem.EncodeToMemory(&pem.Block{Type:  "RSA PUBLIC KEY",Bytes: x509.MarshalPKCS1PublicKey(pubkey)}))
  return pubkey_pem
}

// Returns the private key as string in pem format.
func ExportPrivateKeyAsPemStr(key Key) string {
  privatekey_pem := string(pem.EncodeToMemory(&pem.Block{Type:  "RSA PRIVATE KEY",Bytes: x509.MarshalPKCS1PrivateKey(key)}))
  return privatekey_pem
}

// Create a Key from the bytes of a private key in pem format.
func LoadPrivateKeyFromPemStr(pemStr []byte)(Key,error) {
  privPem, _ := pem.Decode(pemStr)
	if privPem.Type != "RSA PRIVATE KEY" {
    return nil,errors.New("RSA private key is of the wrong type")
  }
  key, err := x509.ParsePKCS1PrivateKey(privPem.Bytes)
	if err != nil {return nil,err}
	return key,nil
}

// Create a public key from the bytes of a private key in pem format.
func LoadPublicKeyFromPemStr(pemStr []byte)(*rsa.PublicKey,error) {
  pubPem, _ := pem.Decode(pemStr)
	if pubPem.Type != "RSA PUBLIC KEY" {
    return nil,errors.New("RSA public key is of the wrong type")
  }
  key, err := x509.ParsePKCS1PublicKey(pubPem.Bytes)
	if err != nil {return nil,err}
	return key,nil
}

// Returns the public key (as Addr) from a Key.
func GetAddr(key Key)Addr{
  lines := strings.Split(ExportPublicKeyAsPemStr(key),"\n")
  return Addr(strings.Join(lines[1:len(lines)-2],""))
}

// Returns the public key (as Addr) from the private key.
func GetAddr2(key *rsa.PublicKey)Addr{
  lines := strings.Split(ExportPublicKeyAsPemStr2(key),"\n")
  return Addr(strings.Join(lines[1:len(lines)-2],""))
}

// Converts an Addr to a *rsa.PublicKey
func publicKeyFromAddr(addr Addr)*rsa.PublicKey{
  var pemStr string = "-----BEGIN RSA PUBLIC KEY-----\n"
  pemStr += string(addr)+"\n"
  pemStr += "-----END RSA PUBLIC KEY-----"
  pub,err := LoadPublicKeyFromPemStr([]byte(pemStr))
  if err!=nil{fmt.Println(err)}
  return pub
}

// Can be used as builder to store data and obtain a signature.
type SignBuilder struct{
  data []byte
}

// Adds data to the SignBuilder. The data can be of any type
// and is read in []byte.
func (self *SignBuilder)Add(data interface{}){
  var binary []byte = []byte(fmt.Sprintf("%v",data))
  self.data = append(self.data,binary...)
}

// Returns the signature of the collection of the submitted
// data of a SignBuilder.
func (self *SignBuilder)GetSignature(key Key)[]byte{
  var opts rsa.PSSOptions
  opts.SaltLength = rsa.PSSSaltLengthAuto
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

// Signs data from the hash, using a private key.
func GetSignatureFromHash(hashed string,key Key)[]byte{
  var opts rsa.PSSOptions
  opts.SaltLength = rsa.PSSSaltLengthAuto
  newhash := crypto.SHA256
  decoded_hash,_ := hex.DecodeString(hashed)
  signature,err := rsa.SignPSS(rand.Reader,key,newhash,decoded_hash,&opts)
  if err != nil {
  	fmt.Println(err)
    return nil
  }
  return signature
}

// Check the signature of some data from the hash and
// a public key.
func CheckSignature(sign,hashed string,addr Addr)bool{
  var opts rsa.PSSOptions
  opts.SaltLength = rsa.PSSSaltLengthAuto
  newhash := crypto.SHA256
  publicKey := publicKeyFromAddr(addr)
  decoded_sign,_ := hex.DecodeString(sign)
  decoded_hash,_ := hex.DecodeString(hashed)
  err := rsa.VerifyPSS(publicKey,newhash,decoded_hash,decoded_sign,&opts)
  if err != nil {
    fmt.Println(err)
    return false
  }
  return true
}
