package server

import "github.com/codegangsta/negroni"

type Server struct {
	*negroni.Negroni
}
