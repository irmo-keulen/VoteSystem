package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"unicode/utf8"

	"github.com/hashicorp/go-multierror"
)

// Generates new RSA Private/public Key (struct).
// Returns: multierror.Error (if no errors are caught this will be nil)
func genRsaKeyPair(filenamePub, filenamePriv string) error {
	fpPub, _ := filepath.Abs(filenamePub)
	fpPriv, _ := filepath.Abs(filenamePriv)
	tmp1 := "\u2713\u2715"
	tmp2, _ := utf8.DecodeRuneInString(tmp1)
	ck := string(tmp2)

	fmt.Printf("Generating keys: \n- Privpath: \t%s\n- Pubpath:\t%s\n", fpPriv, fpPub)

	if _, err := os.Stat(fpPriv); !os.IsNotExist(err) { // Checks if private key already exists
		fmt.Printf("%v PrivKey Path Exists.\nskipping generating keys.\n", ck)
		return nil
	}
	var err *multierror.Error
	privKey, rsaErr := rsa.GenerateKey(rand.Reader, 4096)
	multierror.Append(rsaErr, err)
	multierror.Append(savePEMKey(filenamePriv, privKey), err)                // Saves private Key
	multierror.Append(savePublicPEMKey(filenamePub, privKey.PublicKey), err) // Saves Public Key
	fmt.Printf("%v Done Generating Keys", ck)
	return err.ErrorOrNil() //Returns nil if no errors has been caught
}

// Saves PrivateKey to file, if file doens't exist a file will be created.
// @Param: Filename (string) => name of the file created, can be used to specify a path aswel.
// @Error: Error during write to file.
func savePEMKey(filename string, key *rsa.PrivateKey) error {
	pemFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Failed to Create/Read file. Err: %v", err.Error())
	}
	var privateKey = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}
	err = pem.Encode(pemFile, privateKey)
	if err != nil {
		return fmt.Errorf("Failed to write privKey to file. Err: %s", err.Error())
	}
	return nil
}

// Saves PubKey to file, if file doens't exist a file will be created.
// @Param: Filename (string) => name of the file created, can be used to specify a path aswel.
// @Panic: Error during write to file.
func savePublicPEMKey(filename string, pubKey rsa.PublicKey) error {
	asn1Bytes, err := asn1.Marshal(pubKey)
	if err != nil {
		return fmt.Errorf("Failed to mashal pubKey. Err: %s", err)
	}
	var pemkey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}
	pemFile, _ := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	err = pem.Encode(pemFile, pemkey)
	defer pemFile.Close()
	return nil
}
