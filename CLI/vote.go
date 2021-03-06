package main

import (
	"bufio"
	"bytes"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Used for casting a vote.
// Sign will be an encrypted hash value

// Prompts user for the vote
// Right now this only a for or against vote
func voteProcess() error {
	// TODO Encrypt vote process
	voteSubject, err := getVoteSub(pubKey, userCode)
	if err != nil {
		return fmt.Errorf("error retreiving vote subject : %s", err.Error())
	}
	fmt.Printf("The vote today is about %s\nAre you for(1) or against(2)?", voteSubject)
	reader := bufio.NewReader(os.Stdin)
	vote := false

mainLoop:
	for {
		val, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Got the following error: %s\nTrying again.", err.Error())
		}
		switch val {
		case "1\n":
			vote = true
			break mainLoop
		case "2\n":
			vote = false
			break mainLoop
		default:
			fmt.Printf("Got %s, which cant be decoded.\nPlease use either 1/2", val)
		}
	}
	postVote(vote)
	return nil
}

// Handles the voting process connection wise
func postVote(voteVal bool) {
	// Encrypt request
	h := sha512.New()
	h.Write([]byte(fmt.Sprintf("%s%t", userCode, voteVal)))
	vote := castVote{
		UserCode: userCode,
		VoteVal:  voteVal,
		Hash:     h.Sum(nil)}
	var err error
	vote.Sign, err = SignMessage(vote.preSign(), BytesToPrivateKey(privKey))
	if err != nil {
		log.Fatalf("error posting vote, %s", err.Error())
	}
	castVoteJSON, err := json.Marshal(vote)
	req, err := http.NewRequest("POST", "http://localhost:8000/vote/cast", bytes.NewBuffer(castVoteJSON))
	if err != nil {
		log.Fatalf("error generating request : %s", err.Error())
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("error sending post request : %s", err.Error())
	}
	rsp, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(rsp))
	return
}

// Retrieves vote subject from server.
// This is a wrapper around the getVote func.
// Handles exchanging keys etc.
func getVoteSub(pubKey []byte, userCode string) (string, error) {
	user := userCred{userCode, pubKey}
	pKeyServer, err := exchangeKey(user, keyUrl)
	if err != nil {
		log.Fatal(err)
	}
	if pKeyServer == nil {
		fmt.Printf("No key was parsed")
	}
	sub, err := getVote(pKeyServer)
	if err != nil {
		return "", fmt.Errorf("error retrieving voting subject : %s", err.Error())
	}

	return sub, nil
}
