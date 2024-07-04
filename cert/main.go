package main

import (
	"bytes"
	cryptorand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"time"

	netutils "k8s.io/utils/net"
)

func main() {
	ca := `
	`
	caKey := `
	`

	cadata, keydata, err := genApiServerCert([]byte(ca), []byte(caKey), "", nil, []string{"www.baidu.com"})
	if err != nil {
		panic(err)
	}
	fmt.Println(string(cadata), string(keydata))
}

const (
	PrivateKeyBlockType    = "PRIVATE KEY"
	PublicKeyBlockType     = "PUBLIC KEY"
	CertificateBlockType   = "CERTIFICATE"
	RSAPrivateKeyBlockType = "RSA PRIVATE KEY"
	rsaKeySize             = 2048
)

func ParsePem(ca, cakey []byte) (*x509.Certificate, interface{}, error) {
	caBlock, _ := pem.Decode(ca)
	caDer, err := x509.ParseCertificate(caBlock.Bytes)
	if err != nil {
		return nil, nil, err
	}

	caKeyBlock, _ := pem.Decode(cakey)
	var caKeyDer *rsa.PrivateKey
	caKeyDer, err = x509.ParsePKCS1PrivateKey(caKeyBlock.Bytes)
	if err != nil {
		//caKeyDer, err = x509.ParsePKCS8PrivateKey(caKeyBlock.Bytes)
		if err != nil {
			return nil, nil, err
		}
	}
	return caDer, caKeyDer, nil
}

func genApiServerCert(ca, cakey []byte, host string, alternateIPs []net.IP, alternateDNS []string) ([]byte, []byte, error) {
	caDer, caKeyDer, err := ParsePem(ca, cakey)
	if err != nil {
		panic(err)
	}

	//"names": [{
	//    "C": "CN",
	//    "ST": "BeiJing",
	//    "L": "BeiJing",
	//    "O": "k8s",
	//    "OU": "System"
	//}],
	validFrom := time.Now().Add(-time.Hour)
	maxAge := time.Hour * 24 * 365 * 99
	template := x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject: pkix.Name{
			CommonName:         "kubernetes",        // CN
			Country:            []string{"CN"},      // C
			Organization:       []string{"k8s"},     // O
			OrganizationalUnit: []string{"System"},  // Ou
			Locality:           []string{"BeiJing"}, // L
			Province:           []string{"BeiJing"}, // ST
		},
		NotBefore: validFrom,
		NotAfter:  validFrom.Add(maxAge),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	if ip := netutils.ParseIPSloppy(host); ip != nil {
		template.IPAddresses = append(template.IPAddresses, ip)
	} else {
		template.DNSNames = append(template.DNSNames, host)
	}
	template.IPAddresses = append(template.IPAddresses, alternateIPs...)
	template.DNSNames = append(template.DNSNames, alternateDNS...)

	// 私钥
	priv, err := rsa.GenerateKey(cryptorand.Reader, rsaKeySize)
	if err != nil {
		panic(err)
	}

	// 证书
	derBytes, err := x509.CreateCertificate(cryptorand.Reader, &template, caDer, &priv.PublicKey, caKeyDer)
	if err != nil {
		panic(err)
	}

	keyBuffer := bytes.Buffer{}
	if err := pem.Encode(&keyBuffer, &pem.Block{Type: RSAPrivateKeyBlockType, Bytes: x509.MarshalPKCS1PrivateKey(priv)}); err != nil {
		panic(err)
	}

	certBuffer := bytes.Buffer{}
	if err := pem.Encode(&certBuffer, &pem.Block{Type: CertificateBlockType, Bytes: derBytes}); err != nil {
		panic(err)
	}

	return certBuffer.Bytes(), keyBuffer.Bytes(), nil
}
