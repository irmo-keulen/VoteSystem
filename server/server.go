package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	filenamePub, filenamePriv := "./pub_key", "./priv_key"
	err := genRsaKeyPair(filenamePub, filenamePriv)
	if err != nil {
		fmt.Printf("Generating Keys returned the following error. err: %v", err.Error())
	}
	r := mux.NewRouter()

	r.HandleFunc("/", index)
	r.HandleFunc("/api/pubkey", sendPubKey).Methods("GET")
	r.HandleFunc("/api/pubkey", retrieveKey).Methods("POST")

	log.Fatal(http.ListenAndServe("localhost:8000", r))
}

// Returns Hello, World. Exclusively for testing purposes.
func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World"))
}

func sendPubKey(w http.ResponseWriter, r *http.Request) {
	key, err := ioutil.ReadFile("./pub_key")
	if err != nil {
		panic(err)
	}
	w.Write(key)
}

// Used to retrieve Public key from CLI-user
// - Methods Allowed : POST
// - Returns         : nil
// - TODO            : Write key to DB
func retrieveKey(w http.ResponseWriter, r *http.Request) {
	msg, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Whoops Something went wrong, Try again."))
	}
	keyPem, _ := pem.Decode(msg)
	if keyPem.Type != "RSA PUBLIC KEY" {
		w.Write([]byte("Field type isn't correct"))
	}
	parsedKey, err := x509.ParsePKCS1PublicKey(keyPem.Bytes)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(parsedKey)
}
