package main

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Handles setting up the environment
// Does the following things:
// Checks if private key file exists
// Creates a new key pair if key doesn't exists
// Writes the keyBytes to key files
func setup() {
	fpPub, _ := filepath.Abs(filenamePublic)
	fpPriv, _ := filepath.Abs(filenamePrivate)
	fmt.Printf("Generating keys: \n- Privpath: \t%s\n- Pubpath:\t%s\n", fpPriv, fpPub)
	if _, err := os.Stat(fpPriv); !os.IsNotExist(err) { // Checks if private key already exists
		fmt.Printf("%v PrivKey Path Exists.\nskipping generating keys.\n", ck)
		return
	}
	privKey, pubKey := GenerateKeyPair(4096)
	privFile, err := os.OpenFile(filenamePrivate, os.O_CREATE|os.O_WRONLY, os.ModePerm)
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

	pubFile, err := os.OpenFile(filenamePublic, os.O_CREATE|os.O_WRONLY, os.ModePerm)
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

func getVote(pubKeyServer *rsa.PublicKey) (string, error) {
	voteSub := vote{}
	msg := EncryptWithPublicKey([]byte(userCode), pubKeyServer)
	req, err := http.NewRequest("POST", getVoteUrl, bytes.NewBuffer(msg))
	if err != nil {
		return "", fmt.Errorf("error creating request 1 : %s", err.Error())
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request : %s", err.Error())
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
	sMsg := signedMessage{}
	json.Unmarshal(body, &sMsg)
	str, err := decryptMsg(sMsg.Vote)
	if err != nil {
		return "", fmt.Errorf("error decrypting get vote : %s", err.Error())
	}
	err = json.Unmarshal([]byte(str), &voteSub)
	if err != nil {
		return str, fmt.Errorf("error unmarshalling JSON : %s", err.Error())
	}

	if !voteSub.checkHash() {
		return voteSub.Subject, fmt.Errorf("hash isn't correct")
	}
	if !voteSub.checkSign(sMsg.Sign, pubKeyServer) {
		return voteSub.Subject, fmt.Errorf("sign isn't correct")
	}
	return voteSub.Subject, nil
}

// Handles getting the pub/priv key.
// Looks for the filenames declared in cli.go
func getKeys() (privKey []byte, pubKey []byte, err error) {
	pubKey, err = getKey(filenamePublic)
	if err != nil {
		return privKey, pubKey, fmt.Errorf("error parsing pubKey : %s", err.Error())
	}
	privKey, err = getKey(filenamePrivate)
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
