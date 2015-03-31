## Staticserve

Simple web server in Go

Usage:

```
./staticserve -webroot=/var/www/html -port=8080
```

Point your browser at `http://hostname:8080`

### HTML5Mode

When the `-html5mode` flag is supplied, the web server will return the contents of the root path (i.e. index.html) whenever:

- a 404 error occurs
- the path is not a file (i.e. does not have a filename extension)

This allows JS frameworks such as AngularJS to use HTML5 history API for a better user experience. One caveat is the case where you really do want a 404 to be returned (e.g. for SEO). In this case the `-404Path` flag can be supplied which specifies a path that will generate a hard 404 error e.g. `-404Path=/404.html`. The AngularJS app can then point to this path for any URLs that need a hard 404.
