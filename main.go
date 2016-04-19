package main

import "github.com/osvaldshpengler/browsercalls/server"

func main() {
	s := server.NewServer()

	s.Run(":8000");
}
