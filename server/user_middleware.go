package server

import (
	"database/sql"
	"errors"
	"github.com/gorilla/context"
	"github.com/osvaldshpengler/browsercalls/tools"
	"net/http"
	"encoding/gob"
)

const SESSION_COOKIE_NAME = "bs_auth_session_id"

var ErrUserValidation = errors.New("user_middleware: cookie credentials mismatch")

type User struct {
	Id       int
	Username string
	Password string
}

func init() {
	gob.Register(&User{})
}

func userMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	sessStore := tools.GetSessionStore()
	session, err := sessStore.Get(r, SESSION_COOKIE_NAME)
	if nil != err || session.IsNew {
		next(rw, r)
		return
	}

	val := session.Values["user"]
	u, ok := val.(*User)
	if !ok {
		next(rw, r)
		return
	}
	if err = validateUser(u); nil != err {
		if err == sql.ErrNoRows || err == ErrUserValidation {
			next(rw, r)
			return
		} else {
			tools.Log.Error(err)
		}
	}

	context.Set(r, "user", u)
	session.Save(r, rw)

	next(rw, r)
}

func validateUser(u *User) error {
	dba, err := tools.GetDbAccessor()
	if err != nil {
		return err
	}

	var username, password string
	err = dba.QueryRow("SELECT username, password FROM users WHERE id = $1", u.Id).Scan(
		&username,
		&password,
	)
	if nil != err {
		return err
	}

	if u.Username != username || u.Password != password {
		return ErrUserValidation
	}

	return nil
}

func initUserSession(rw http.ResponseWriter, r *http.Request, u *User) error {
	sessStore := tools.GetSessionStore()
	session, _ := sessStore.Get(r, SESSION_COOKIE_NAME)
	session.Values["user"] = u
	return session.Save(r, rw)
}
