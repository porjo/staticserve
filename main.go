package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/codegangsta/negroni"
	"github.com/phyber/negroni-gzip/gzip"
)

func redir(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.TLS == nil {
		http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
	} else {
		next(w, r)
	}
}

func main() {

	var err error
	var port, portTLS int

	webRoot := flag.String("d", "public", "root directory of website")
	certfile := flag.String("certfile", "cert.pem", "SSL certificate filename")
	keyfile := flag.String("keyfile", "key.pem", "SSL key filename")
	flag.IntVar(&port, "p", 8080, "HTTP port")
	flag.IntVar(&portTLS, "s", 8081, "HTTPS port")
	stripPrefix := flag.String("stripPrefix", "", "path to strip from incoming requests")
	flag.Parse()

	_, err = os.Lstat(*webRoot)
	if err != nil {
		log.Fatalf("error opening webroot: %s\n", *webRoot, err)
	} else {
		log.Printf("using webroot '%s'\n", *webRoot)
	}

	mux := http.NewServeMux()

	if *stripPrefix != "" {
		log.Printf("stripPrefix '%s'\n", *stripPrefix)
		// Handle either case of trailing slash
		stripPrefixTrim := strings.TrimRight(*stripPrefix, "/")
		mux.Handle(stripPrefixTrim+"/", http.StripPrefix(stripPrefixTrim+"/", http.FileServer(http.Dir(*webRoot))))
	} else {
		mux.Handle("/", http.FileServer(http.Dir(*webRoot)))

	}

	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
	)
	//n.Use(negroni.HandlerFunc(redir))
	n.Use(gzip.Gzip(gzip.DefaultCompression))
	n.UseHandler(mux)
	go func() {
		log.Printf("HTTPS listening on port %d\n", portTLS)
		log.Fatal(http.ListenAndServeTLS(":"+strconv.Itoa(portTLS), *certfile, *keyfile, n))
	}()
	log.Printf("HTTP listening on port %d\n", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), n))
}
