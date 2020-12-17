package main

import (
	"fmt"
	"log"
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
		log.Fatal(err)
	}
	user := userCred{"testingCode", pubKey}
	pKeyServer, err := exchangeKey(user, "http://localhost:8000/api/pubkey")
	if err != nil {
		log.Fatal(err)
	}
	if pKeyServer == nil {
		fmt.Printf("")
	}
}
