package main

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

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
		fmt.Println(fmt.Errorf("error parsing data : %s", err.Error()))
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
		fmt.Printf("error reading publicKey : %s", err.Error())
	}
	_, _ = w.Write(key)
}

func getVote(w http.ResponseWriter, r *http.Request) {
	privKey, err := ioutil.ReadFile(filenamePriv)
	if err != nil {
		log.Fatalf("error reading private key : %s", err.Error())
	}
	encMsg, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error during getVote : %s", err.Error())
	}

	msg, err := decryptMsg(encMsg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error during decrypting message : %s", err.Error())
	}
	k, err := rdb.Get(ctx, msg).Result()
	if err != nil {
		fmt.Fprintf(os.Stderr, "User uknown : %s", err.Error())
		w.Write([]byte("No Public key is found, please send to /api/pubkey"))
		r.Body.Close()
		return
	}
	h := sha512.New()
	h.Write([]byte(voteSubject))
	sign, err := SignMessage([]byte(voteSubject), BytesToPrivateKey(privKey))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error generating sign : %s", err.Error())
	}
	v := vote{voteSubject, h.Sum(nil)}
	voteJSON, _ := json.Marshal(v)
	encVoteJSOn := EncryptWithPublicKey(voteJSON, BytesToPublicKey([]byte(k)))
	sMsg := signedMessage{encVoteJSOn, sign}
	val, err := json.Marshal(sMsg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't parse json")
	}
	_, _ = w.Write(val)
	r.Body.Close()
}

func handleVote(w http.ResponseWriter, r *http.Request) {
	vote := castVote{}
	voteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading cast vote. Err: %s", err.Error())
		r.Body.Close()
		return
	}
	err = json.Unmarshal(voteReq, &vote)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error retreiving data from body. Err: %s", err.Error())
		r.Body.Close()
	}
	pubKey, err := rdb.Get(ctx, vote.UserCode).Result()
	if !vote.checkSign(BytesToPublicKey([]byte(pubKey))) {
		w.Write([]byte("Sign isn't correct\nNo vote has been cast"))
	}
	fmt.Printf("%t\n", vote.VoteVal)
	err = vrdb.Set(ctx, vote.UserCode, fmt.Sprintf("%t", vote.VoteVal), 0).Err()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing vote to database. Err: %s", err.Error())
		w.Write([]byte("Error handling vote, please try again"))
	}
	w.Write([]byte("You have voted"))
	r.Body.Close()
	return
}
