package tools

import (
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

type SessionManager struct {
	*sessions.CookieStore
}

var hashKey = securecookie.GenerateRandomKey(64)
var blockKey = securecookie.GenerateRandomKey(64)

var sessionStore = sessions.NewCookieStore(hashKey, blockKey)


func GetSessionStore() *sessions.CookieStore {
	return sessionStore
}
