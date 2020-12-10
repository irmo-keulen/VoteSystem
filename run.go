package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	keys := genRsaKeyPair()
	route(mux.NewRouter())
}

func route(r *mux.Router) {
	r.HandleFunc("/api/pubkey", sendPubKey)
}

func sendPubKey(http.ResponseWriter, *http.Request) {

}
