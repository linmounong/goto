/*
A goto server that redirect a configured key to a URL.

PATH /a general apis
	* /a/search
		* GET /a/search?q=<key>&page=1&&per=10&owner=
		* GET /a/search/<key>?page=1&&per=10&owner=
	* /a/manage
		* GET /a/manage?key=<key>
		* POST/PUT /a/manage?key=<key>
		* DELETE /a/manage?key=<key>

PATH /s static files
  * /s/<file> static file <file>

PATH /l login logout
  * /l login logout

PATH /<key> redirect
  * /<key> redirect

Rules

TODO(ynlin): details and examples.

1. A key may only cantain [a-zA-Z0-9-_.] e.g. "hello-world"/"go1.6.1". The length is at least 2.

2. A URL may only contain scheme/userinfo/host/path, host mandatory. See https://golang.org/pkg/net/url/#URL.

3. When redirecting, scheme/userinfo/host/path are replaced if configured.

4. When redirecting, trailing "/" is added to path only when necessary.

*/
package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
)

var (
	p = flag.Int("p", 8096, "HTTP port")
	s = flag.Bool("s", false, "use local static assets")
)

func main() {
	flag.Parse()

	DbOpen()
	defer DbClose()

	http.Handle("/", dumped(http.HandlerFunc(gotoHandler)))
	http.Handle("/a/", dumped(http.StripPrefix("/a", http.HandlerFunc(apiHandler))))
	http.Handle("/h/", dumped(http.StripPrefix("/h", http.HandlerFunc(hashHandler))))
	if *s {
		log.Println("using local static assets")
		http.Handle("/s/", dumped(http.StripPrefix("/s/", http.FileServer(http.Dir("s")))))
	} else {
		http.Handle("/s/", dumped(http.StripPrefix("/s/", http.FileServer(assetFS()))))
	}
	log.Println("serving on port", *p)
	http.ListenAndServe(":"+strconv.Itoa(*p), nil)
}
