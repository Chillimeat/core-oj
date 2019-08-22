package main

import (
	"flag"
	"log"

	server "github.com/Myriad-Dreamin/core-oj/compiler/compilerx/src"
)

var (
	port = flag.String("port", ":23366", "port to listen")
)

func init() {
	flag.Parse()
}

func main() {
	if srv, err := server.NewServer(); err != nil {
		log.Fatal(err)
	} else if err := srv.ListenAndServe(*port); err != nil {
		log.Fatal(err)
	}
}
