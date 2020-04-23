/**
 * @Author: Bizzaro Francesco <d33pblue>
 * @Date:   2020-Apr-23
 * @Project: Proof of Evolution
 * @Filename: rsa.go
 * @Last modified by:   d33pblue
 * @Last modified time: 2020-Apr-23
 * @Copyright: 2020
 */

package wallet

import (
    "crypto"
    "crypto/rand"
    "crypto/rsa"
    "crypto/sha256"
    "fmt"
    "crypto/x509"
    "encoding/pem"
)

func GenerateKey()(*rsa.PrivateKey,error) {
  reader := rand.Reader
	return rsa.GenerateKey(reader,2048)
}

func EncryptMessage(message []byte,yourPrivKey *rsa.PrivateKey,
        otherPublKey *rsa.PublicKey){
  label := []byte("")
  hash := sha256.New()
  ciphertext,err := rsa.EncryptOAEP(hash,rand.Reader,otherPublKey,message,label)
  fmt.Printf("OAEP encrypted [%s] to \n[%x]\n", string(message), ciphertext)
  // Message - Signature
	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto // for simple example
	PSSmessage := message
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)
	signature,err := rsa.SignPSS(rand.Reader,yourPrivKey,newhash,hashed,&opts)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("PSS Signature : %x\n",signature)
}


func DecryptMessage(ciphertext,hashed,signature []byte,yourPrivKey *rsa.PrivateKey,
        otherPublKey *rsa.PublicKey ){
  label := []byte("")
  hash := sha256.New()
  plainText, err := rsa.DecryptOAEP(hash,rand.Reader,yourPrivKey,ciphertext,label)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("OAEP decrypted [%x] to \n[%s]\n", ciphertext, plainText)
	//Verify Signature
  var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto
  newhash := crypto.SHA256
	err = rsa.VerifyPSS(otherPublKey,newhash,hashed,signature,&opts)
	if err != nil {
		fmt.Println("Who are U? Verify Signature failed")
	} else {
		fmt.Println("Verify Signature successful")
	}
}

func ExportPublicKeyAsPemStr(pubkey *rsa.PublicKey) string {
    pubkey_pem := string(pem.EncodeToMemory(&pem.Block{Type:  "RSA PUBLIC KEY",Bytes: x509.MarshalPKCS1PublicKey(pubkey)}))
    return pubkey_pem
}

func ExportPrivateKeyAsPemStr(privatekey *rsa.PrivateKey) string {
    privatekey_pem := string(pem.EncodeToMemory(&pem.Block{Type:  "RSA PRIVATE KEY",Bytes: x509.MarshalPKCS1PrivateKey(privatekey)}))
    return privatekey_pem
}
