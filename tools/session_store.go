package tools

import (
	"github.com/gorilla/sessions"
	"os"
)

type SessionManager struct {
	*sessions.CookieStore
}

var hashKey = os.Getenv("BC_APP_SECRET")
var blockKey = os.Getenv("BC_ENC_KEY")

var sessionStore = sessions.NewCookieStore([]byte(hashKey), []byte(blockKey))

func GetSessionStore() *sessions.CookieStore {
	return sessionStore
}
