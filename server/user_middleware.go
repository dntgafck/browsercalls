package server

import (
	"errors"
	"fmt"
	"github.com/gorilla/context"
	"github.com/osvaldshpengler/browsercalls/tools"
	"net/http"
	"strings"
	"time"
)

const AUTH_COOKIE_NAME = "bs_auth_session_id"

var ErrUnserializeAuth = errors.New("user_middleware: unserialize cookie error")
var ErrWrongCredentials = errors.New("user_middleware: cookie credentials mismatch")

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

	cId, cUsername, cPassword, cEmail, err := unserializeAuthCookie(authCookie)
	if nil != err {
		tools.Log.Info(err, authCookie)
		next(rw, r)
	}

	dba, err := tools.GetDbAccessor()
	if err != nil {
		tools.Log.Error(err)
	}

	var dUsername, dPassword, dEmail string
	err = dba.QueryRow("SELECT username, email, password FROM users WHERE id = $1", cId).Scan(
		&dUsername,
		&dEmail,
		&dPassword,
	)
	if nil != err {
		tools.Log.Error(err)
	}
	if dUsername != cUsername || dEmail != cEmail || cPassword != dPassword {
		tools.Log.Info(ErrWrongCredentials)
		next(rw, r)
	}

	cValue := serializeAuthCookie(cId, dUsername, dPassword, dEmail)
	cOptions := map[string]interface{}{
		"expires": time.Now().Add(time.Hour),
	}
	cm.Set(rw, AUTH_COOKIE_NAME, cValue, cOptions)

	user := &User{cId, dUsername, dPassword, dEmail}
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
