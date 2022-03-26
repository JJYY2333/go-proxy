/*
@Time    : 3/26/22 09:03
@Author  : Neil
@File    : user.go
*/

package socks

type Session struct {
	target string
	user   *User
}

func NewSession() *Session {
	s := new(Session)
	return s
}

// AddUser if socks don't use auth, add nil user
func (s *Session) AddUser(u *User) {
	s.user = u
}

func (s *Session) AddTarget(tgt string) {
	s.target = tgt
}

func (s *Session) GetUname() string {
	return s.user.uname
}

func (s *Session) GetTarget() string {
	return s.target
}

type User struct {
	uname string
}

func NewAuthUser(uname string) *User {
	u := new(User)
	u.uname = uname
	return u
}

func NewAnonymousUser() *User {
	u := new(User)
	u.uname = "Anonymous"
	return u
}
