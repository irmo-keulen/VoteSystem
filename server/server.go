package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

var (
	filenamePub  = "./pub_key"
	filenamePriv = "./priv_key"
)

func main() {
	setup()
	//if err != nil {
	//	fmt.Printf("Generating Keys returned the following error. err: %v", err.Error())
	//}
	err := testDBConnect()
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	r := mux.NewRouter()
	fmt.Printf("%s Setup Finished.\nServer listening on localhost:8000\n", ck)
	r.HandleFunc("/", index)
	r.HandleFunc("/api/pubkey", retrieveKey).Methods("POST")
	//r.HandleFunc("/showcase/encrypt", func(w http.ResponseWriter, r *http.Request) {
	//	readHashedBytes(w, r, privateKey)
	//})
	r.HandleFunc("/testing", decryptMessage).Methods("POST")

	log.Fatal(http.ListenAndServe("localhost:8000", r))
}

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
