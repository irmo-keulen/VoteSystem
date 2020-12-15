package main

import (
	"crypto/rsa"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Used to indentify user
type userCred struct {
	Usercode       string `json:usercode`
	Username       string `json:username`
	Voted          bool   `json:voted`
	publicKey      *rsa.PublicKey
	PublicKeyBytes []byte `json:publickey`
}

func (u *userCred) String() string {
	return fmt.Sprintf("{\"usercode\":\"%s\",\"username\":\"%s\",\"voted\":%t,\"publicKey\":\"%s\"}",
		u.Usercode, u.Username, u.Voted, u.PublicKeyString)
}

func main() {
	filenamePub, filenamePriv := "./pub_key", "./priv_key"
	err := genRsaKeyPair(filenamePub, filenamePriv)
	if err != nil {
		fmt.Printf("Generating Keys returned the following error. err: %v",
			err.Error())
	}

	// For test purposes
	fmt.Println(checkUserCred(user{"test1", "pass1"}))

	r := mux.NewRouter()
	fmt.Printf("Setup Finished.\nServer listening on localhost:8000\n")
	r.HandleFunc("/", index)
	r.HandleFunc("/api/pubkey", sendPubKey).Methods("GET")
	r.HandleFunc("/api/pubkey", retrieveKey).Methods("POST")

	log.Fatal(http.ListenAndServe("localhost:8000", r))
}
