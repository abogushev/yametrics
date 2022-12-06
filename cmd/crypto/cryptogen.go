package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"yametrics/internal/crypto"
)

func main() {
	// создаём шаблон сертификата
	//cert := &x509.Certificate{
	//	// указываем уникальный номер сертификата
	//	SerialNumber: big.NewInt(1658),
	//	// заполняем базовую информацию о владельце сертификата
	//	Subject: pkix.Name{
	//		Organization: []string{"Yandex.Praktikum"},
	//		Country:      []string{"RU"},
	//	},
	//	// разрешаем использование сертификата для 127.0.0.1 и ::1
	//	IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
	//	// сертификат верен, начиная со времени создания
	//	NotBefore: time.Now(),
	//	// время жизни сертификата — 10 лет
	//	NotAfter:     time.Now().AddDate(10, 0, 0),
	//	SubjectKeyId: []byte{1, 2, 3, 4, 6},
	//	// устанавливаем использование ключа для цифровой подписи,
	//	// а также клиентской и серверной авторизации
	//	ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
	//	KeyUsage:    x509.KeyUsageDigitalSignature,
	//}

	// создаём новый приватный RSA-ключ длиной 4096 бит
	// обратите внимание, что для генерации ключа и сертификата
	// используется rand.Reader в качестве источника случайных данных

	// создаём сертификат x.509
	//certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	//if err != nil {
	//	log.Fatal(err)
	//}

	// кодируем сертификат и ключ в формате PEM, который
	// используется для хранения и обмена криптографическими ключами
	//var certPEM bytes.Buffer
	//pem.Encode(&certPEM, &pem.Block{
	//	Type:  "CERTIFICATE",
	//	Bytes: certBytes,
	//})
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatal(err)
	}
	//
	//var privateKeyPEM bytes.Buffer
	//pem.Encode(&privateKeyPEM, &pem.Block{
	//	Type:  "RSA PRIVATE KEY",
	//	Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	//})

	pemPrivateFile, err := os.Create("private_key.pem")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = pem.Encode(pemPrivateFile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	pemPrivateFile.Close()

	pemPublicFile, err := os.Create("public_key.pem")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = pem.Encode(pemPublicFile, &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	pemPublicFile.Close()

	pubkey, err := crypto.ReadPublicKey(pemPublicFile.Name())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	privkey, err := crypto.ReadPrivateKey(pemPrivateFile.Name())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	r, err := crypto.Encrypt(pubkey, []byte("hello"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	m, err := crypto.Decrypt(privkey, r)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(m))
}
