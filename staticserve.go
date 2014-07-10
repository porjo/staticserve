package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
)

func main() {
	webroot := flag.String("webroot", "/usr/share/doc", "path to web files")
	port := flag.Int("port", 8080, "port to listen on")
	flag.Parse()

	log.Printf("Listening on port %d\n", *port)
	panic(http.ListenAndServe(":"+strconv.Itoa(*port), http.FileServer(http.Dir(*webroot))))
}
