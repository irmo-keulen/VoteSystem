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

type voteVal bool

type votemessage struct {
	userCode string
	vote     voteVal
	hash     []byte
}

// Used for casting a vote.
// Sign will be an encrypted hash value
type castVote struct {
	UserCode string `json:"user_code"`
	VoteVal  bool   `json:"vote_val"`
	Hash     []byte `json:"hash"`
	Sign     []byte `json:"sign"`
}

// Creates a byte array of the values userCode and voteVal
// returns a Byte Array
func (v *castVote) preSign() []byte {
	return []byte(fmt.Sprintf("%s%t", v.UserCode, v.VoteVal))
}

// Prompts user for the vote
// Right now this only a for or against vote
func voteProcess(voteSubject string) bool {
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
	return true
}

// Handles the voting process connection wise
func postVote(voteVal bool) {
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
