package main

import (
	"log"
)

func main() {
	setup()
	var err error
	privKey, pubKey, err = getKeys()
	if err != nil {
		log.Fatalf("Coulnd't parse key : %s", err.Error())
	}
	err = voteProcess()
	if err != nil {
		log.Fatalf("error Completing vote process : %s", err.Error())
	}
}

// TODO sign castVote message
