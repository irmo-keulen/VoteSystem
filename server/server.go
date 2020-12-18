package main

import (
	"fmt"
	"log"
	"net/http"

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
	writeVote()
	r := mux.NewRouter()
	fmt.Printf("%s Setup Finished.\nServer listening on localhost:8000\n", ck)
	r.HandleFunc("/", index)
	r.HandleFunc("/api/pubkey", retrieveKey).Methods("POST")
	r.HandleFunc("/api/getvote", getVote).Methods("POST")
	r.HandleFunc("/vote/cast", handleVote).Methods("POST")

	log.Fatal(http.ListenAndServe("localhost:8000", r))
}

// TODO Vote Process
