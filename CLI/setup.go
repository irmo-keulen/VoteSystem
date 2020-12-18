package main

import (
	"bytes"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var (
	filenamePub  = "./pub_key"
	filenamePriv = "./priv_key"
	keyUrl       = "http://localhost:8000/api/pubkey"
	getVoteUrl   = "http://localhost:8000/api/getvote"
	userCode     = "1234HelloWorld!"
)

// Handles setting up the environment
// Does the following things:
// Checks if private key file exists
// Creates a new key pair if key doesn't exists
// Writes the keyBytes to key files
func setup() {
	fpPub, _ := filepath.Abs(filenamePub)
	fpPriv, _ := filepath.Abs(filenamePriv)
	fmt.Printf("Generating keys: \n- Privpath: \t%s\n- Pubpath:\t%s\n", fpPriv, fpPub)
	if _, err := os.Stat(fpPriv); !os.IsNotExist(err) { // Checks if private key already exists
		fmt.Printf("%v PrivKey Path Exists.\nskipping generating keys.\n", ck)
		return
	}
	privKey, pubKey := GenerateKeyPair(4096)
	privFile, err := os.OpenFile(filenamePriv, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		fmt.Println("Error generating privKey file")
		log.Fatal(err)
	}
	_, err = privFile.Write(PrivateKeyToBytes(privKey))
	if err != nil {
		fmt.Println("Error writing to privFile")
		log.Fatal(err)
	}
	defer privFile.Close()

	pubFile, err := os.OpenFile(filenamePub, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		fmt.Println("Error generating pubKey file")
		log.Fatal(err)
	}
	_, err = pubFile.Write(PublicKeyToBytes(pubKey))
	if err != nil {
		fmt.Println("Error writing to pubFile")
		log.Fatal(err)
	}
}

func getVote(privkey *rsa.PrivateKey, pubKeyServer *rsa.PublicKey) (string, error) {
	type vote struct {
		Subject string `json:"subject"`
		Hash    []byte `json:"hash"`
	}
	voteSub := vote{}
	msg := EncryptWithPublicKey([]byte("testingCode"), pubKeyServer)
	req, err := http.NewRequest("POST", getVoteUrl, bytes.NewBuffer(msg))
	if err != nil {
		return "", fmt.Errorf("error creating request : %s", err.Error())
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error creating request : %s", err.Error())
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return string(body), fmt.Errorf("error finishing POST request : %s", err.Error())
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error closing response body : %s", err.Error())
		}
	}()
	str, err := decryptMsg(body)
	if err != nil {
		return "", fmt.Errorf("error decrypting get vote : %s", err.Error())
	}
	err = json.Unmarshal([]byte(str), &voteSub)
	if err != nil {
		return str, fmt.Errorf("error unmarshalling JSON : %s", err.Error())
	}
	h := sha512.New()
	fmt.Println(voteSub.Subject)
	h.Write([]byte(voteSub.Subject))
	if bytes.Compare(voteSub.Hash, h.Sum(nil)) != 0 {
		fmt.Fprintf(os.Stderr, "\nExpected: %v\n", h.Sum(nil))
		fmt.Fprintf(os.Stderr, "\nGot: %v\n", voteSub.Hash)
		return voteSub.Subject, fmt.Errorf("hash isn't correct")
	}
	return voteSub.Subject, nil
}

// Handles getting the pub/priv key.
func getKeys(privPath string, pubpath string) (privKey []byte, pubKey []byte, err error) {
	pubKey, err = getKey(filenamePub)
	if err != nil {
		return privKey, pubKey, fmt.Errorf("error parsing pubKey : %s", err.Error())
	}
	privKey, err = getKey(filenamePriv)
	if err != nil {
		return privKey, pubKey, fmt.Errorf("error parsing privKey : %s", err.Error())
	}
	return privKey, pubKey, nil

}

// Reads the key file and returns public key
// Returns as byte slice
func getKey(filename string) ([]byte, error) {
	keyFile, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to read file. Err: %s", err.Error())
	}
	key, err := ioutil.ReadAll(keyFile)
	if err != nil {
		return nil, fmt.Errorf("error reading file Err: %s", err.Error())
	}
	return key, nil
}

func getVoteSub(pubKey []byte, privKey []byte, userCode string) (string, error) {
	user := userCred{userCode, pubKey}
	pKeyServer, err := exchangeKey(user, keyUrl)
	if err != nil {
		log.Fatal(err)
	}
	if pKeyServer == nil {
		fmt.Printf("No key was parsed")
	}
	sub, err := getVote(BytesToPrivateKey(privKey), pKeyServer)
	if err != nil {
		return "", fmt.Errorf("error retrieving voting subject : %s", err.Error())
	}
	return sub, nil
}
