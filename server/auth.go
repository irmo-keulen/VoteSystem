package main

import (
	"fmt"
)

var (
	user1 userCred     = userCred{"test1", "pass1"}
	user2 userCred     = userCred{"test2", "pass2"}
	users userAuthCred = userAuthCred{[]userCred{user1, user2}}
)

// Used to indentify user
type userCred struct {
	Usercode  string `json:"usercode"`
	PublicKey string `json:"publickey"`
}

func (u *userCred) String() string {
	return fmt.Sprintf("{\"usercode\":\"%s\",\"publicKey\":\"%s\"}",
		u.Usercode, u.PublicKey)
}

type userAuthCred struct {
	Users []userCred
}

// // Reads file to check if username & usercode exist
// // - Return: bool (true if user exists)
// func (u *userAuthCred) authenticate(canUser user) bool {
// 	for _, val := range u.Users {
// 		if val.Username == canUser.Username && val.Usercode == canUser.Usercode {
// 			return true
// 		}
// 	}
// 	return false
// }

// TODO
//      Handle auth process.
