package keystore_test

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/zacksfF/Build-A-Go-Apps-That-Scale-On-Google-Cloud/util/keystore"
)

const testDataDir = "testdata" // Assuming you have a testdata directory with PEM files

func TestKeyStore_LoadRSAKeys(t *testing.T) {
	ks := keystore.New()

	err := ks.LoadRSAKeys(os.DirFS(testDataDir))
	if err != nil {
		t.Errorf("Unexpected error loading keys: %v", err)
	}

	// Check if keys are loaded correctly (number of keys, existence of specific key)
	if len(ks.store) != 2 { // Adjust the expected number of keys based on your testdata
		t.Errorf("Expected to load 2 keys, got %d", len(ks.store))
	}

	_, ok := ks.store["key1.pem"] // Adjust key name based on your testdata
	if !ok {
		t.Errorf("Expected key 'key1.pem' not found in keystore")
	}
}

func TestKeyStore_PublicKey(t *testing.T) {
	ks := keystore.New()
	err := ks.LoadRSAKeys(os.DirFS(testDataDir))
	if err != nil {
		t.Fatal(err)
	}

	// Test successful retrieval of public key
	pubKey, err := ks.PublicKey("key1.pem") // Adjust key name based on your testdata
	if err != nil {
		t.Errorf("Unexpected error getting public key: %v", err)
	}

	// Perform additional checks on the retrieved public key (e.g., length, format)

	// Test retrieval of non-existent key
	_, err = ks.PublicKey("non-existent-key.pem")
	if err == nil || !errors.Is(err, keystore.ErrKeyNotFound) {
		t.Errorf("Expected ErrKeyNotFound for non-existent key, got: %v", err)
	}
}

func TestKeyStore_PrivateKey(t *testing.T) {
	ks := keystore.New()
	err := ks.LoadRSAKeys(os.DirFS(testDataDir))
	if err != nil {
		t.Fatal(err)
	}

	// Test successful retrieval of private key
	privKey, err := ks.PrivateKey("key1.pem") // Adjust key name based on your testdata
	if err != nil {
		t.Errorf("Unexpected error getting private key: %v", err)
	}

	// Perform additional checks on the retrieved private key (consider security implications)

	// Test retrieval of non-existent key
	_, err = ks.PrivateKey("non-existent-key.pem")
	if err == nil || !errors.Is(err, keystore.ErrKeyNotFound) {
		t.Errorf("Expected ErrKeyNotFound for non-existent key, got: %v", err)
	}
}

func TestToPublicPem(t *testing.T) {
	testData, err := ioutil.ReadFile(filepath.Join(testDataDir, "private_key.pem")) // Adjust path based on your testdata
	if err != nil {
		t.Fatal(err)
	}

	// Test successful conversion to public PEM
	pubPEM, err := keystore.ToPublicPem(string(testData))
	if err != nil {
		t.Errorf("Unexpected error converting to public PEM: %v", err)
	}

	// Perform additional checks on the public PEM (e.g., format)

	// Test invalid PEM format
	invalidPEM := "invalid_pem_data"
	_, err = keystore.ToPublicPem(invalidPEM)
	if err == nil || !strings.Contains(err.Error(), "invalid key") {
		t.Errorf("Expected error for invalid PEM, got: %v", err)
	}
}
