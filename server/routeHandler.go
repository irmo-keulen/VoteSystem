package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Returns Hello, World. Exclusively for testing purposes.
func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World"))
}

// Returns public key
// - Methods Allowed : GET
// - Returns 		 : string("hello, World")
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
	cred := userCred{}
	msg, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Whoops Something went wrong, Please try again."))
	}
	err = json.Unmarshal(msg, &cred)
	if err != nil {
		fmt.Println(fmt.Errorf("Error parsing data, %s", err.Error()))
		return
	}
	fmt.Printf("%v", cred.String())
}
