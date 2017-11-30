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

	"github.com/elazarl/go-bindata-assetfs"
)

var (
	p = flag.Int("p", 8096, "HTTP port")
)

func main() {
	flag.Parse()

	DbOpen()
	defer DbClose()

	http.Handle("/", dumped(http.HandlerFunc(gotoHandler)))
	http.Handle("/a/", dumped(http.StripPrefix("/a", http.HandlerFunc(apiHandler))))
	http.Handle("/h/", dumped(http.StripPrefix("/h", http.HandlerFunc(hashHandler))))
	http.Handle("/s/", http.StripPrefix("/s/", http.FileServer(
		&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo, Prefix: "s"})))
	log.Println("serving on port", *p)
	http.ListenAndServe(":"+strconv.Itoa(*p), nil)
}
