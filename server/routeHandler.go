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

// Used to retrieve Public key from CLI-user
// - Methods Allowed : POST
// - Returns         : Public Key Server
// - TODO            : Write key to DB
func retrieveKey(w http.ResponseWriter, r *http.Request) {
	var cred userCred
	msg, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Whoops Something went wrong, Please try again."))
	}
	err = json.Unmarshal(msg, &cred)
	if err != nil {
		fmt.Println(fmt.Errorf("Error parsing data, %s", err.Error()))
		return
	}
	// privb, _ := pem.Decode(cred.PublicKeyPEM)
	if err != nil {
		fmt.Println(err)
		w.Write([]byte("An error occured"))
	}
	fmt.Printf("%v", cred)
	// Returns own public key.
	key, err := ioutil.ReadFile("./pub_key")
	if err != nil {
		panic(err)
	}
	w.Write(key)
}

// TODO:
//     write publickey to redis database (key = usercode)
