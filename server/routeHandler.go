package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type userCred struct {
	Usercode  string `json:"usercode"`
	PublicKey string `json:"publickey"`
}

// Returns Hello, World. Exclusively for testing purposes.
func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World"))
}

// Used to retrieve Public key from CLI-user
// - Methods Allowed : POST
// - Returns         : Public Key Server
func retrieveKey(w http.ResponseWriter, r *http.Request) {
	var cred userCred
	msg, err := ioutil.ReadAll(r.Body)
	if err != nil {
		_, _ = w.Write([]byte("Whoops Something went wrong, Please try again."))
	}
	err = json.Unmarshal(msg, &cred)
	if err != nil {
		fmt.Println(fmt.Errorf("Error parsing data, %s", err.Error()))
		_, _ = w.Write([]byte(`{"http-code":500}`))
		return
	}
	err = rdb.Set(ctx, cred.Usercode, cred.PublicKey, 0).Err()
	if err != nil {
		panic(err)
	}
	// Returns own public key.
	key, err := ioutil.ReadFile("./pub_key")
	if err != nil {
		panic(err)
	}
	_, _ = w.Write(key)
}

func decryptMessage(w http.ResponseWriter, r *http.Request) {
	privKey, err := ioutil.ReadFile(filenamePriv)
	msg, err := ioutil.ReadAll(r.Body)
	fmt.Println(string(msg))
	if err != nil {
		fmt.Println("Error")
	}
	fmt.Println(string(DecryptWithPrivateKey(msg, BytesToPrivateKey(privKey))))

}
