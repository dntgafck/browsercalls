package server

import (
	"github.com/gorilla/context"
	"net/http"
)

func contextClearMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer context.Clear(r)
	next(rw, r)
}
