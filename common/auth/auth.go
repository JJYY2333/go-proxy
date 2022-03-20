/*
@Time    : 3/20/22 10:19
@Author  : Neil
@File    : auth.go
*/

package auth

import "log"

type Authenticator interface {
	Check(uname, passwd string) bool
}

type DummyAuth struct {
	record map[string]string
}

func NewDummyAuth() *DummyAuth {
	d := new(DummyAuth)
	d.record = map[string]string{"Neil": "Neil"}
	return d
}

func (a *DummyAuth) Check(uname, passwd string) bool {
	if len(uname) == 0 || len(passwd) == 0 {
		log.Printf("invalid uname or passwd")
		return false
	}

	recordPwd, ok := a.record[uname]
	if !ok {
		log.Printf("user not found: %s", uname)
		return false
	}

	if passwd != recordPwd {
		log.Printf("wrong password: %s", passwd)
		return false
	}

	return true
}