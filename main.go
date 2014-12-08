package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
	"github.com/phyber/negroni-gzip/gzip"
)

// ResponseWriter wrapper to catch 404s
type html5mode struct {
	w http.ResponseWriter
	r *http.Request
}

var h5m html5mode
var webRoot string

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

	html5mode := flag.Bool("html5mode", false, "On HTTP 404, serve index.html. Used with AngularJS html5mode.")
	flag.StringVar(&webRoot, "d", "public", "root directory of website")
	certfile := flag.String("certfile", "", "SSL certificate filename")
	keyfile := flag.String("keyfile", "", "SSL key filename")
	flag.IntVar(&port, "p", 8080, "HTTP port")
	flag.IntVar(&portTLS, "s", 8081, "HTTPS port")
	forceTLS := flag.Bool("forceTLS", false, "Force HTTPS")
	stripPrefix := flag.String("stripPrefix", "", "path to strip from incoming requests")
	flag.Parse()

	_, err = os.Lstat(webRoot)
	if err != nil {
		log.Fatalf("error opening webroot: %s\n", err)
	} else {
		log.Printf("using webroot '%s'\n", webRoot)
	}

	mux := http.NewServeMux()

	if *stripPrefix != "" {
		log.Printf("stripPrefix '%s'\n", *stripPrefix)
		// Handle either case of trailing slash
		stripPrefixTrim := strings.TrimRight(*stripPrefix, "/")
		mux.Handle(stripPrefixTrim+"/", http.StripPrefix(stripPrefixTrim+"/", http.FileServer(http.Dir(webRoot))))
	} else {
		mux.Handle("/", http.FileServer(http.Dir(webRoot)))

	}

	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
	)
	if *forceTLS && *certfile != "" && *keyfile != "" {
		log.Printf("Force TLS enabled\n")
		n.Use(negroni.HandlerFunc(redir))
	}
	n.Use(gzip.Gzip(gzip.DefaultCompression))

	if *html5mode {
		log.Printf("HTML5 mode enabled\n")
		n.Use(negroni.HandlerFunc(html5ModeMiddleware))
	}

	n.UseHandler(mux)

	if *certfile != "" && *keyfile != "" {
		go func() {
			log.Printf("HTTPS listening on port %d\n", portTLS)
			log.Fatal(http.ListenAndServeTLS(":"+strconv.Itoa(portTLS), *certfile, *keyfile, n))
		}()
	}

	log.Printf("HTTP listening on port %d\n", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), n))
}

// This should come before static file-serving middleware
func html5ModeMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	h5m.w = rw
	h5m.r = r
	next(&h5m, r)
	//context.Clear(r)
}

func (sr *html5mode) Header() http.Header { return sr.w.Header() }
func (sr *html5mode) Write(d []byte) (int, error) {
	if _, ok := context.GetOk(sr.r, "html5modeWritten"); ok {
		return 0, nil
	} else {
		return sr.w.Write(d)
	}
}

func (sr *html5mode) WriteHeader(status int) {
	if status == 404 {
		if path.Ext(sr.r.URL.Path) == "" {
			// Server contents of index.html if request isn't for a file
			// Required for Angular.js html5mode
			sr.w.Header().Del("Content-Type")
			http.ServeFile(sr.w, sr.r, webRoot+"/index.html")
			context.Set(sr.r, "html5modeWritten", true)
			return
		}
	}
	sr.w.WriteHeader(status)
}
