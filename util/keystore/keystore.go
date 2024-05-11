//We create a simple key store for the project to store and
//retrieve the private key.

// Package keystore implements the auth.KeyLookup interface. This implements
// an in-memory keystore for JWT support.
package keystore

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path"
	"strings"
)

// Key represent key information.
type Key struct {
	privatPEM string
	publicPEM string
}

// Keystore represent an in memory store implementation of the
// KeyLoop interface for use with the auth package
type KeyStore struct {
	store map[string]Key
}

// New Constructs an empty Keystore ready for use .
func New() *KeyStore {
	return &KeyStore{
		store: make(map[string]Key),
	}
}

// LoadRSAKeys loads a set of RSA PEM files rooted inside of a directory. The
// name of each PEM file will be used as the key id .
func (ks *KeyStore) LoadRSAKeys(fsys fs.FS) error {
	fn := func(fileName string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("Walkdir failure: %w", err)
		}

		if dirEntry.IsDir() {
			return nil
		}

		if path.Ext(fileName) != ".pem" {
			return nil
		}

		file, err := fsys.Open(fileName)
		if err != nil {
			return fmt.Errorf("Opening key file: %w", err)
		}
		defer file.Close()

		//limit PEM file size to 1 megabyte. This should be reasonable for
		//almost any PEM file and prevents shenanigans like linking the file
		//to .dev/random or something like that.
		pem, err := io.ReadAll(io.LimitReader(file, 1024*1024))
		if err != nil {
			return fmt.Errorf("reading auth private key: %w", err)
		}

		privatePEM := string(pem)
		publicPEM, err := toPublicPem(privatePEM)
		if err != nil {
			return fmt.Errorf("converting private PEM to public: %w", err)
		}

		key := Key{
			privatPEM: privatePEM,
			publicPEM: publicPEM,
		}
		ks.store[strings.TrimSuffix(dirEntry.Name(), ".pem")] = key
		return nil
	}

	if err := fs.WalkDir(fsys, ".", fn); err != nil {
		return fmt.Errorf("walking directory: %w", err)
	}

	return nil
}

// PublicKey searches the key store for a given kid and returns the public key.
func (ks *KeyStore) PublicKey(kid string) (string, error) {
	key, found := ks.store[kid]
	if !found {
		return "", errors.New("kid lookup failed")
	}
	return key.publicPEM, nil
}

// PrivateKey searches the key store for a given kid and returns the private key.
func (ks *KeyStore) PrivateKey(kid string) (string, error) {
	key, found := ks.store[kid]
	if !found {
		return "", errors.New("kid lookup failed")
	}

	return key.privatPEM, nil
}

func toPublicPem(privatePEM string) (string, error) {
	block, _ := pem.Decode([]byte(privatePEM))
	if block == nil {
		return "", errors.New("invalid key: Key must be a PEM encoded PKCS1or PKCS28 key")
	}

	var parsedKey any
	parsedKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		parsedKey, err = x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return "", err
		}
	}
	pk, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return "", errors.New("key is not a valid RSA private key")
	}

	asn1Bytes, err := x509.MarshalPKIXPublicKey(&pk.PublicKey)
	if err != nil {
		return "", fmt.Errorf("marshaling public key: %w", err)
	}

	publicBlock := pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	var buf bytes.Buffer
	if err := pem.Encode(&buf, &publicBlock); err != nil {
		return "", fmt.Errorf("encoding to public PEM: %w", err)
	}

	return buf.String(), nil

}
