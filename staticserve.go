package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	webroot := flag.String("webroot", "/usr/share/doc", "path to web files")
	stripPrefix := flag.String("stripPrefix", "", "path to strip from incoming requests")
	port := flag.Int("port", 8080, "port to listen on")
	flag.Parse()

	log.Printf("Listening on port %d\n", *port)

	if *stripPrefix != "" {
		log.Printf("stripPrefix '%s'\n", *stripPrefix)
		// Handle either case of trailing slash
		stripPrefixTrim := strings.TrimRight(*stripPrefix, "/")
		http.Handle(stripPrefixTrim+"/", http.StripPrefix(stripPrefixTrim+"/", http.FileServer(http.Dir(*webroot))))
	} else {
		http.Handle("/", http.FileServer(http.Dir(*webroot)))
	}
	panic(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}
