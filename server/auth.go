package main

var (
	user1 user         = user{"test1", "pass1"}
	user2 user         = user{"test2", "pass2"}
	users userAuthCred = userAuthCred{[]user{user1, user2}}
)

type user struct {
	Username string
	Usercode string
}

type userAuthCred struct {
	Users []user
}

// Reads file to check if username & usercode exist
// - Return: bool (true if user exists)
func (u *userAuthCred) authenticate(canUser user) bool {
	for _, val := range u.Users {
		if val.Username == canUser.Username && val.Usercode == canUser.Usercode {
			return true
		}
	}
	return false
}
