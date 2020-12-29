package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	err := setup()
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	r := mux.NewRouter()
	fmt.Printf("%s Setup Finished.\nServer listening on localhost:8000\n", ck)
	r.HandleFunc("/", index)
	r.HandleFunc("/api/pubkey", retrieveKey).Methods("POST")
	r.HandleFunc("/api/getvote", getVote).Methods("POST")
	r.HandleFunc("/vote/cast", handleVote).Methods("POST")

	log.Fatal(http.ListenAndServe("localhost:8000", r))
}
