package main

import (
	"context"
	"crypto/rsa"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var (
	voteSubject  string
	ck           = "\u2713"
	filenamePub  = "./pub_key"
	filenamePriv = "./priv_key"
	ctx          = context.Background()

	//Global DB-Connector for Public Keys
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	//Global DB-Connector for Voting
	vrdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	})
)

type userCred struct {
	Usercode  string `json:"usercode"`
	PublicKey []byte `json:"publickey"`
}
type vote struct {
	Subject string `json:"subject"`
	Hash    []byte `json:"hash"`
}

func (v *vote) byte() []byte {
	s := fmt.Sprintf("%s%s", v.Subject, v.Hash)
	return []byte(s)
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

func (v *castVote) checkSign(key *rsa.PublicKey) bool {
	return VerifySign(v.Sign, v.preSign(), key)
}
