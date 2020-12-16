package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type userCred struct {
	Usercode  string `json:"usercode"`
	PublicKey []byte `json:"publickey"`
}

type message struct {
	Msg []byte `json:"message"`
}

var (
	filenamePub  = "./pub_key"
	filenamePriv = "./priv_key"
	keyUrl       = "http://localhost:8000/api/pubkey"
)

func main() {
	setup()
	pubKey, err := getPubKey(filenamePub)
	if err != nil {
		log.Fatal(err) // TODO check if file exists, otherwise create new key pair.
	}
	user := userCred{"testingCode", pubKey}
	pKeyServer, err := exchangeKey(user, "http://localhost:8000/api/pubkey")
	if err != nil {
		log.Fatal(err)
	}
	msg := EncryptWithPublicKey([]byte("Hello, World"), pKeyServer)
	req, err := http.NewRequest("POST", "http://localhost:8000/testing", bytes.NewBuffer(msg))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	fmt.Println(body)
}

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

// TODO
//Refactor Crypto (using new code base)
//Implement byte writer to file (pub/priv key)
