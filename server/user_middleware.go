package server

import (
	"net/http"
	"fmt"
	"strings"
	"errors"
	"github.com/osvaldshpengler/browsercalls/tools"
	"github.com/gorilla/context"
)

const AUTH_COOKIE_NAME = "bs_auth_session_id"

var ErrUnserializeAuth = errors.New("user_middleware: unserialize cookie error")

type User struct {
	Id       int
	Username string
	Email    string
	Password string
}

func userMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	cm := tools.GetCookieManager()
	authCookie, err := cm.Get(r, AUTH_COOKIE_NAME)
	if nil != err {
		next(rw, r)
		return
	}

	id, username, password, email, err := unserializeAuthCookie(authCookie)
	if nil != err {
		next(rw, r)
	}

	user := &User{id, username, email, password}

	context.Set(r, "user", user)
}

func serializeAuthCookie(id int, username, password, email string) string {
	return fmt.Sprintf("%d|%s|%s|%s", id, username, password, email)
}

func unserializeAuthCookie(value string) (id int, username, password, email string, err error) {
	unserialized := strings.Split(value, "|")
	if 3 != len(unserialized) {
		err = ErrUnserializeAuth
		return
	}

	id = int(unserialized[0])
	username = unserialized[1]
	password = unserialized[2]
	email = unserialized[3]

	return
}