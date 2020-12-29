// File is used for declaring al custom types and methods
package main

import (
	"bytes"
	"crypto/rsa"
	"crypto/sha512"
	"fmt"
)

var (
	filenamePublic  = "./pub_key"
	filenamePrivate = "./priv_key"
	keyUrl          = "http://localhost:8000/api/pubkey"
	getVoteUrl      = "http://localhost:8000/api/getvote"
	userCode        = "1234HelloWorld!"
	ck              = "\u2713"
	privKey         []byte
	pubKey          []byte
)

// Used for constructing a message JSON
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

// Used for decoding the vote subject
// TODO add Signing process
type vote struct {
	Subject string `json:"subject"`
	Hash    []byte `json:"hash"`
}

type signedMessage struct {
	Vote []byte `json:"vote"`
	Sign []byte `json:"sign"`
}

// Compares the parsed hash with a calculated hash
func (v *vote) checkHash() bool {
	h := sha512.New()
	h.Write([]byte(v.Subject))
	if bytes.Compare(v.Hash, h.Sum(nil)) != 0 {
		return false
	}
	return true
}

// Verifies the sign of a vote.
// Returns true if sign is correct
func (v *vote) checkSign(sign []byte, pubKey *rsa.PublicKey) bool {
	return VerifySign(sign, []byte(v.Subject), pubKey)
}

// Used to identify CLI to the server
type userCred struct {
	Usercode  string `json:"usercode"`
	PublicKey []byte `json:"publickey"`
}
