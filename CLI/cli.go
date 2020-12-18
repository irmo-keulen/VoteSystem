package main

import (
	"log"
)

type userCred struct {
	Usercode  string `json:"usercode"`
	PublicKey []byte `json:"publickey"`
}

var (
	privKey []byte
	pubKey  []byte
)

func main() {
	setup()
	var err error
	privKey, pubKey, err = getKeys(filenamePriv, filenamePub)
	if err != nil {
		log.Fatalf("Coulnd't parse key : %s", err.Error())
	}
	voteSubject, err := getVoteSub(pubKey, privKey, userCode)
	if err != nil {
		log.Fatalf("Error getting vote subject : %s", err.Error())
	}
	voteProcess(voteSubject)

}

// TODO Refactor error handling
// TODO Vote process.
