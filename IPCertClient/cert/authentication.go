package cert

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"os"
)

func PathIsExist(path string) bool {
	_, err := os.Stat(path)
	var exist = false
	if err == nil {
		exist = true
	}
	if os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func Write(filename, Type string, p []byte) error {
	File, err := os.Create(filename)
	defer File.Close()
	if err != nil {
		return err
	}
	var b = &pem.Block{Bytes: p, Type: Type}
	err = pem.Encode(File, b)
	if err != nil {
		return err
	}
	return nil
}

func CreateKeyAndCSR(filePath string,name string,phone string) error{
	dirPath := filePath // "./crypto"
	err0 := os.MkdirAll(dirPath, os.ModePerm)
	if err0 != nil {
		return err0
	}

	privateKeyFilePath := filePath + "/client.key"
	csrFilePath := filePath + "/client.csr"

	// 生成公私钥对
	privateKey, _ := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	//publicKey := privateKey.PublicKey
	// 保存私钥到本地
	privB, _ := x509.MarshalECPrivateKey(privateKey)
	err := Write(privateKeyFilePath, "PRIVATE KEY", privB)
	if err != nil {
		return err
	}

	// 生成csr
	req := &x509.CertificateRequest{
		Subject: pkix.Name{
			Country:            []string{"China"},
			Organization:       []string{name},
			OrganizationalUnit: []string{phone},
			Locality:           []string{"haidian"},
			Province:           []string{"Beijing"},
			StreetAddress:      []string{"37"},
			//SubjectKeyId:       info.SubjectKeyId,
		},
		PublicKey: privateKey.PublicKey,
	}
	csrByte, _ := x509.CreateCertificateRequest(rand.Reader, req, privateKey)
	err1 := Write(csrFilePath, "CERTIFICATE REQUEST", csrByte)
	if err1 != nil {
		return err1
	}

	return nil
}
