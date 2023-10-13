package sign

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	cryptorand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"testing"
)

func generatePEM(randReader io.Reader, password []byte) (buf *bytes.Buffer, err error) {
	k, err := rsa.GenerateKey(randReader, 1024)
	if err != nil {
		return nil, err
	}

	derBytes := x509.MarshalPKCS1PrivateKey(k)

	var block *pem.Block
	if password != nil {
		block, err = x509.EncryptPEMBlock(randReader, "RSA PRIVATE KEY", derBytes, password, x509.PEMCipherAES128)
	} else {
		block = &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: derBytes,
		}
	}

	buf = &bytes.Buffer{}
	err = pem.Encode(buf, block)
	return buf, err
}

func TestLoadPemPrivKey(t *testing.T) {
	reader, err := generatePEM(newRandomReader(rand.New(rand.NewSource(1))), nil)
	if err != nil {
		t.Errorf("Unexpected pem generation err %s", err.Error())
	}

	privKey, err := LoadPEMPrivKey(reader)
	if err != nil {
		t.Errorf("Unexpected key load error, %s", err.Error())
	}
	if privKey == nil {
		t.Errorf("Expected valid privKey, but got nil")
	}
}

func TestLoadPemPrivKeyInvalidPEM(t *testing.T) {
	reader := strings.NewReader("invalid PEM data")
	privKey, err := LoadPEMPrivKey(reader)

	if err == nil {
		t.Errorf("Expected error invalid PEM data error")
	}
	if privKey != nil {
		t.Errorf("Expected nil privKey but got %#v", privKey)
	}
}

func TestLoadEncryptedPEMPrivKey(t *testing.T) {
	reader, err := generatePEM(newRandomReader(rand.New(rand.NewSource(1))), []byte("password"))
	if err != nil {
		t.Errorf("Unexpected pem generation err %s", err.Error())
	}

	privKey, err := LoadEncryptedPEMPrivKey(reader, []byte("password"))

	if err != nil {
		t.Errorf("Unexpected key load error, %s", err.Error())
	}
	if privKey == nil {
		t.Errorf("Expected valid privKey, but got nil")
	}
}

func TestLoadEncryptedPEMPrivKeyWrongPassword(t *testing.T) {
	reader, err := generatePEM(newRandomReader(rand.New(rand.NewSource(1))), []byte("password"))
	privKey, err := LoadEncryptedPEMPrivKey(reader, []byte("wrong password"))

	if err == nil {
		t.Errorf("Expected error invalid PEM data error")
	}
	if privKey != nil {
		t.Errorf("Expected nil privKey but got %#v", privKey)
	}
}

func TestLoadPEMPrivKeyPKCS8_RSA(t *testing.T) {
	r, err := generatePKCS8(keyTypeRSA)
	if err != nil {
		t.Errorf("generate pkcs8 key: %v", err)
	}

	var rr bytes.Buffer
	tee := io.TeeReader(r, &rr)

	key, err := LoadPEMPrivKeyPKCS8(tee)
	if err != nil {
		t.Errorf("load pkcs8 key: %v", err)
	}

	if _, ok := key.(*rsa.PrivateKey); !ok {
		t.Errorf("key should be rsa but was %T", key)
	}

	_, err = LoadPEMPrivKeyPKCS8AsSigner(&rr)
	if err != nil {
		t.Errorf("load pkcs8 key as signer: %v", err)
	}
}

func TestLoadPEMPrivKeyPKCS8_ECDSA(t *testing.T) {
	r, err := generatePKCS8(keyTypeECDSA)
	if err != nil {
		t.Errorf("generate pkcs8 key: %v", err)
	}

	var rr bytes.Buffer
	tee := io.TeeReader(r, &rr)

	key, err := LoadPEMPrivKeyPKCS8(tee)
	if err != nil {
		t.Errorf("load pkcs8 key: %v", err)
	}

	if _, ok := key.(*ecdsa.PrivateKey); !ok {
		t.Errorf("key should be ecdsa but was %T", key)
	}

	_, err = LoadPEMPrivKeyPKCS8AsSigner(&rr)
	if err != nil {
		t.Errorf("load pkcs8 key as signer: %v", err)
	}
}

func TestLoadPEMPrivKeyPKCS8_ED25519(t *testing.T) {
	r, err := generatePKCS8(keyTypeED25519)
	if err != nil {
		t.Errorf("generate pkcs8 key: %v", err)
	}

	var rr bytes.Buffer
	tee := io.TeeReader(r, &rr)

	key, err := LoadPEMPrivKeyPKCS8(tee)
	if err != nil {
		t.Errorf("load pkcs8 key: %v", err)
	}

	if _, ok := key.(ed25519.PrivateKey); !ok {
		t.Errorf("key should be ed25519 but was %T", key)
	}

	_, err = LoadPEMPrivKeyPKCS8AsSigner(&rr)
	if err != nil {
		t.Errorf("load pkcs8 key as signer: %v", err)
	}
}

type keyType int

const (
	keyTypeRSA keyType = iota
	keyTypeECDSA
	keyTypeED25519
)

func generatePKCS8(typ keyType) (io.Reader, error) {
	key, pemType, err := generatePrivateKey(typ)
	if err != nil {
		return nil, err
	}

	b, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return nil, err
	}

	block := &pem.Block{
		Type:  pemType,
		Bytes: b,
	}

	var buf bytes.Buffer
	if err := pem.Encode(&buf, block); err != nil {
		return nil, err
	}

	return &buf, err
}

func generatePrivateKey(typ keyType) (interface{}, string, error) {
	switch typ {
	case keyTypeRSA:
		key, err := rsa.GenerateKey(cryptorand.Reader, 1024)
		return key, "RSA PRIVATE KEY", err
	case keyTypeECDSA:
		key, err := ecdsa.GenerateKey(elliptic.P224(), cryptorand.Reader)
		return key, "ECDSA PRIVATE KEY", err
	case keyTypeED25519:
		_, key, err := ed25519.GenerateKey(cryptorand.Reader)
		return key, "ED25519 PRIVATE KEY", err
	default:
		return nil, "", fmt.Errorf("unsupported key type %v", typ)
	}
}
