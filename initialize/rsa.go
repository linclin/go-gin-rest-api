package initialize

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/gob"
	"encoding/pem"
	"fmt"
	"os"

	"go-gin-rest-api/pkg/global"
)

func InitRSA() {
	_, err := os.Stat("conf/rsa/rsa-private.pem")
	if err != nil {
		if !os.IsExist(err) {
			GenerateRSA()
			global.Log.Info("初始化RSA证书供JWT生成Toekn使用完成")
		}
	}
}

// openssl genrsa -out conf/rsa/rsa-private.pem 2048
// openssl rsa -in conf/rsa/rsa-private.pem -pubout > conf/rsa/rsa-public.pem
func GenerateRSA() {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	checkError(err)
	publicKey := key.PublicKey
	saveGobKey("conf/rsa/rsa-private.key", key)
	savePEMKey("conf/rsa/rsa-private.pem", key)
	saveGobKey("conf/rsa/rsa-public.key", publicKey)
	savePublicPEMKey("conf/rsa/rsa-public.pem", &publicKey)
}

func saveGobKey(fileName string, key interface{}) {
	outFile, err := os.Create(fileName)
	checkError(err)
	defer outFile.Close()
	encoder := gob.NewEncoder(outFile)
	err = encoder.Encode(key)
	checkError(err)
}

func savePEMKey(fileName string, key *rsa.PrivateKey) {
	outFile, err := os.Create(fileName)
	checkError(err)
	defer outFile.Close()
	info := struct {
		Version             int
		PrivateKeyAlgorithm []asn1.ObjectIdentifier
		PrivateKey          []byte
	}{}

	info.Version = 0
	info.PrivateKeyAlgorithm = make([]asn1.ObjectIdentifier, 1)
	info.PrivateKeyAlgorithm[0] = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 1}
	info.PrivateKey = x509.MarshalPKCS1PrivateKey(key)
	prvKey, _ := asn1.Marshal(info)

	var privateKey = &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: prvKey,
	}
	err = pem.Encode(outFile, privateKey)
	checkError(err)
}

func savePublicPEMKey(fileName string, pubkey *rsa.PublicKey) {
	pemfile, err := os.Create(fileName)
	checkError(err)
	defer pemfile.Close()
	asn1Bytes, err := x509.MarshalPKIXPublicKey(pubkey)
	checkError(err)
	var pemkey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}
	err = pem.Encode(pemfile, pemkey)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		global.Log.Debug(fmt.Sprint("Fatal error:GenerateRSA", err.Error()))
		os.Exit(1)
	}
}
