package main

import (
	"fmt"
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
	err = testDBConnect()
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	r := mux.NewRouter()
	fmt.Printf("%s Setup Finished.\nServer listening on localhost:8000\n", ck)
	r.HandleFunc("/", index)
	r.HandleFunc("/api/pubkey", retrieveKey).Methods("POST")

	log.Fatal(http.ListenAndServe("localhost:8000", r))
}

// TODO:
//      Implement vote process.
