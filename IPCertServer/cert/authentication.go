package cert

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"math/big"
	"os"
	"time"
)

type CertInformation struct {
	Country            []string
	Organization       []string
	OrganizationalUnit []string
	EmailAddress       []string
	Province           []string
	StreetAddress      []string
	SubjectKeyId       []byte
	Locality           []string
}

// 创建根证书和根私钥
func CreateRootCertAndRootPrivateKey(info CertInformation, rootCertFilePath, rootPrivateKeyFilePath string) error {

	dirPath := "./crypto"
	err0 := os.MkdirAll(dirPath, os.ModePerm)
	if err0 != nil {
		return err0
	}

	// use certInformation to new cert
	certTemp := newCertificate(info)
	// generate private key and public key
	rootPrivateKey, _ := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	// create root cert []byte
	rootCertificateByte, _ := x509.CreateCertificate(rand.Reader, certTemp, certTemp, &rootPrivateKey.PublicKey, rootPrivateKey)
	// get private key []bytes
	rootPrivateKeyByte, _ := x509.MarshalECPrivateKey(rootPrivateKey)

	//fmt.Println("rootCertificateByte:",rootCertificateByte)

	// write root certificate
	err := write(rootCertFilePath, "CERTIFICATE", rootCertificateByte)
	if err != nil {
		return err
	}
	// write root private key
	err1 := write(rootPrivateKeyFilePath, "PRIVATE KEY", rootPrivateKeyByte)
	if err1 != nil {
		return err1
	}

	return nil
}

func newCertificate(info CertInformation) *x509.Certificate {
	return &x509.Certificate{
		SerialNumber: big.NewInt(1653),
		Subject: pkix.Name{
			Country:            info.Country,
			Organization:       info.Organization,
			OrganizationalUnit: info.OrganizationalUnit,
			Province:           info.Province,
			StreetAddress:      info.StreetAddress,
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(10, 0, 0),
		//SubjectKeyId:          info.SubjectKeyId,
		BasicConstraintsValid: true,
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		//EmailAddresses: info.EmailAddress,
	}
}

func newCertificateWithCSR(req *x509.CertificateRequest) *x509.Certificate {
	return &x509.Certificate{
		SerialNumber: big.NewInt(1653),
		Subject:      req.Subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		//SubjectKeyId:          info.SubjectKeyId,
		BasicConstraintsValid: true,
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		//EmailAddresses: info.EmailAddress,
	}
}

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

func write(filename, Type string, p []byte) error {
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

func CreateCertWithCsr(csrBytes []byte) ([]byte, error) {
	var newCert []byte

	// 检测根证书和根私钥是否存在
	rootCertFilePath := "./crypto/ca.crt"
	rootPrivateKeyFilePath := "./crypto/ca.key"
	if PathIsExist(rootCertFilePath) == false {
		return newCert, errors.New("No CA RootCert!")
	}
	if PathIsExist(rootPrivateKeyFilePath) == false {
		return newCert, errors.New("No CA Private Key!")
	}

	csr, err := x509.ParseCertificateRequest(csrBytes)
	if err!=nil {
		return newCert, errors.New("parse csr error!")
	}

	// 使用csr生成临时证书
	tempCert := newCertificateWithCSR(csr)

	// 获取根证书
	certPEM, _ := ioutil.ReadFile(rootCertFilePath)
	block, _ := pem.Decode(certPEM)
	if block == nil {
		return newCert, errors.New("Failed to read root Certificate PEM.")
	}
	rootCertificate, err := x509.ParseCertificate(block.Bytes)
	if err!=nil {
		return newCert, errors.New("Failed to parse root Certificate PEM.")
	}

	// 获取根私钥
	privR, _ := ioutil.ReadFile(rootPrivateKeyFilePath)
	block2, _ := pem.Decode(privR)
	if block2 == nil {
		return newCert, errors.New("Failed to read root private key.")
	}
	rootPrivateKey, err := x509.ParseECPrivateKey(block2.Bytes)
	if err!=nil {
		return newCert, errors.New("Failed to parse root private key.")
	}

	// 签发证书
	newCert, err = x509.CreateCertificate(rand.Reader, tempCert, rootCertificate, csr.PublicKey, rootPrivateKey)
	if err!=nil {
		return newCert, errors.New("Failed to create cert: " + err.Error())
	}

	return newCert,nil
}

// VerifyCAChain 验证证书是不是由根证书签发的
func VerifyCAChain(toVerifiedCertBytes []byte) bool {
	rootCertPath := "./crypto/ca.crt"
	rootPEM, err := ioutil.ReadFile(rootCertPath)
	if err != nil {
		fmt.Println("read file")
		return false
	}
	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM([]byte(rootPEM))
	if !ok {
		fmt.Println("Failed to parse Root Certificate")
		return false
	}

	cert, err := x509.ParseCertificate(toVerifiedCertBytes)
	if err != nil {
		fmt.Println("Failed to parse Certificate PEM: " + err.Error())
		return false
	}

	opts := x509.VerifyOptions{
		Roots: roots,
	}
	if _, err := cert.Verify(opts); err != nil {
		fmt.Println("Failed to parse Certificate PEM: " + err.Error())
		return false
	}
	return true
}