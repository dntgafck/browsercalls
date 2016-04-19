package server

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"github.com/osvaldshpengler/browsercalls/tools"
)

func NewServer() *negroni.Negroni {
	errLogger, _ := tools.Log.Get("error")
	infoLogger, _ := tools.Log.Get("info")

	loggerMW := &negroni.Logger{infoLogger}
	recoverMW := &negroni.Recovery{
		Logger:     errLogger,
		PrintStack: true,
		StackAll:   false,
		StackSize:  1024 * 8,
	}

	n := negroni.New(
		negroni.HandlerFunc(contextClearMiddleware),
		recoverMW,
		loggerMW,
		negroni.NewStatic(http.Dir(os.Getenv("BC_APP_PATH")+"public")),
		negroni.HandlerFunc(userMiddleware),
	)

	lc := &loginController{}

	router := mux.NewRouter()
	router.HandleFunc("/", handleHome)
	router.HandleFunc("/login", lc.handleLogin)
	router.HandleFunc("/register", lc.handleRegister)

	n.UseHandler(router)

	return n
}

func handleHome(rw http.ResponseWriter, r *http.Request) {
	u, ok := context.GetOk(r, "user")
	if !ok {
		http.Redirect(rw, r, "/login", http.StatusFound)
		return
	}
	user, _ := u.(*User)
	fmt.Fprintf(rw, "Hi there, I love %s!", user.Username)
}
