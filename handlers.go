package main

import (
	"errors"
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	urllib "net/url"
	"strconv"
	"strings"
)

const (
	lenRandKey = 7 // same as t.cn
)

var (
	dump = flag.Bool("dump", false, "dump http request")
)

func dumped(h http.Handler) http.Handler {
	if !*dump {
		return h
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		d, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Println(err)
		}
		log.Println("===request===")
		log.Println(string(d))
		h.ServeHTTP(w, r)
		log.Println("===request done===")
	})
}

// Redirects "/h/foo" to "/s#/foo"
func hashHandler(w http.ResponseWriter, r *http.Request) {
	u := r.URL
	u.Fragment = u.Path
	u.Path = "/s"
	http.Redirect(w, r, u.String(), http.StatusFound)
}

func redirectToManage(w http.ResponseWriter, r *http.Request, key string) {
	http.Redirect(w, r, "/s/#/"+key, http.StatusFound)
}

// If it's "/", go to the index.
// If it's "/a-not-registered-key", go to the manage page.
// Otherwise, do redirect.
func gotoHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	// parts[0]==""
	if len(parts) < 2 {
		redirectToManage(w, r, "")
		return
	}
	key := parts[1]
	var u *urllib.URL
	info := DbFind(key)
	if info != nil {
		var err error
		if u, err = urllib.Parse(info.Url); err != nil {
			log.Println(err)
		}
	}
	if u == nil {
		redirectToManage(w, r, key)
		return
	}
	redirected := r.URL
	// Only redirect http to https.
	if redirected.Scheme == "http" && u.Scheme == "https" {
		redirected.Scheme = u.Scheme
	}
	if u.User != nil {
		redirected.User = u.User
	}
	redirected.Host = u.Host
	redirected.Path = u.Path
	if u.RawQuery != "" {
		if redirected.RawQuery == "" {
			redirected.RawQuery = u.RawQuery
		} else {
			redirected.RawQuery = u.RawQuery + "&" + redirected.RawQuery
		}
	}
	if u.Fragment != "" {
		redirected.Fragment = u.Fragment
	}
	if len(parts) > 2 {
		if !strings.HasSuffix(redirected.Path, "/") {
			redirected.Path += "/"
		}
		redirected.Path += strings.Join(parts[2:], "/")
	}
	redirectedLink := redirected.String()
	log.Println(key, "->", redirectedLink)
	go DbIncr(key)
	http.Redirect(w, r, redirectedLink, http.StatusFound)
}

func intArg(r *http.Request, key string, defaultVal int) (int, error) {
	valStr := r.FormValue(key)
	if valStr == "" {
		return defaultVal, nil
	}
	return strconv.Atoi(valStr)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	username := getUsername(r)
	parts := strings.SplitN(r.URL.Path, "/", 3)
	// parts[0]==""
	if len(parts) < 2 {
		writeErrorResponse(w, "action not specified", http.StatusBadRequest)
		return
	}
	action := parts[1]
	switch action {
	case "ok":
		writeOkResponse(w, r)
	case "list":
		fallthrough
	case "search":
		q := r.FormValue("q")
		if q == "" && len(parts) >= 3 {
			q = parts[2]
		}
		urlPattern := r.FormValue("url")
		owner := r.FormValue("owner")
		page, err := intArg(r, "page", 1)
		if err != nil || page <= 0 {
			writeErrorResponse(w, "invalid page", http.StatusBadRequest)
			return
		}
		per, err := intArg(r, "per", 10)
		if err != nil || per <= 0 {
			writeErrorResponse(w, "invalid per", http.StatusBadRequest)
			return
		}
		results := DbSearch(q, urlPattern, owner, page, per)
		writeLinksResponse(w, r, results)

	case "manage":
		key := r.FormValue("key")
		if len(key) < 2 {
			if r.Method != http.MethodPost || key != "" {
				writeErrorResponse(w, "too short", http.StatusBadRequest)
				return
			}
			// Generate a random key if POST requests.
			key = randKey()
		}
		switch r.Method {
		case http.MethodPost:
			fallthrough
		case http.MethodPut:
			url := r.FormValue("url")
			if url == "" {
				writeErrorResponse(w, "missing url", http.StatusBadRequest)
				return
			}
			if err := checkKey(key); err != nil {
				writeErrorResponse(w, err.Error(), http.StatusBadRequest)
				return
			}
			url, err := checkAndNormalizeUrl(url)
			if err != nil {
				writeErrorResponse(w, err.Error(), http.StatusBadRequest)
				return
			}
			if info, err := DbUpdateOrCreate(username, key, url); err != nil {
				writeErrorResponse(w, err.Error(), http.StatusUnauthorized)
			} else {
				writeLinkResponse(w, r, info)
			}
			return

		case http.MethodGet:
			info := DbFind(key)
			writeLinkResponse(w, r, info)
			return

		case http.MethodDelete:
			DbRemove(key)
			writeOkResponse(w, r)
			return

		default:
			writeErrorResponse(w, "method not supported", http.StatusNotImplemented)
			return
		}
	default:
		writeErrorResponse(w, "invalid action "+action, http.StatusBadRequest)
		return
	}
}

func randKey() string {
	key := randString(lenRandKey)
	for DbFind(key) != nil {
		key = randString(lenRandKey)
	}
	return key
}

func setJsonHeaders(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json; charset=utf8")
}

func writeLinksResponse(w http.ResponseWriter, r *http.Request, infos []GotoInfo) {
	setJsonHeaders(w)
	resp := &Response{
		Ok:       true,
		Links:    []*Link{},
		Username: getUsername(r),
	}
	for _, info := range infos {
		resp.Links = append(resp.Links, &Link{
			Key:      info.Key,
			Url:      info.Url,
			Owner:    info.Owner,
			UseCount: info.UseCount,
		})
	}
	if err := resp.marshal(w); err != nil {
		log.Println("marshal error", err)
		writeErrorResponse(w, err.Error(), http.StatusInternalServerError)
	}
}

func writeLinkResponse(w http.ResponseWriter, r *http.Request, info *GotoInfo) {
	setJsonHeaders(w)
	resp := &Response{
		Ok:       true,
		Username: getUsername(r),
	}
	if info != nil {
		resp.Link = &Link{
			Key:      info.Key,
			Url:      info.Url,
			Owner:    info.Owner,
			UseCount: info.UseCount,
		}
	}
	if err := resp.marshal(w); err != nil {
		log.Println("marshal error", err)
		writeErrorResponse(w, err.Error(), http.StatusInternalServerError)
	}
}

func writeOkResponse(w http.ResponseWriter, r *http.Request) {
	setJsonHeaders(w)
	resp := &Response{
		Ok:       true,
		Username: getUsername(r),
	}
	if err := resp.marshal(w); err != nil {
		log.Println("marshal error", err)
		writeErrorResponse(w, err.Error(), http.StatusInternalServerError)
	}
}

func writeErrorResponse(w http.ResponseWriter, errMsg string, code int) {
	log.Println("error", errMsg)
	setJsonHeaders(w)
	w.WriteHeader(code)
	r := &Response{
		ErrMsg: errMsg,
	}
	if err := r.marshal(w); err != nil {
		// Just give up writing back.
		log.Println("marshal error", err)
	}
}

// Checks if key matches [a-zA-Z-_.]* and has a length of at least 2.
func checkKey(key string) error {
	if len(key) < 2 {
		return errors.New("too short")
	}
	for i := 0; i < len(key); i++ {
		c := key[i]
		if !('A' <= c && c <= 'z') && !('0' <= c && c <= '9') && c != '-' && c != '_' && c != '.' {
			return errors.New("invalid char " + string(c))
		}
	}
	return nil
}

// Checks if url is valid. Scheme/User/Host/Path/RawQuery/Fragment will be used, others ignored.
// Http is assumed if Scheme is empty. Only http and https are allowed.
//
// Examples:
//   "/foo/bar" -> "/foo/bar"
//   "foo/bar" -> "http://foo/bar"
//   "https://foo/bar" -> "https://foo/bar"
func checkAndNormalizeUrl(url string) (string, error) {
	u, err := urllib.Parse(url)
	if err != nil {
		return "", err
	}
	if u.Scheme == "" && strings.HasPrefix(url, "/") {
		return url, nil
	}
	if u.Scheme == "" {
		u, err = urllib.Parse("http://" + url)
		if err != nil {
			return "", err
		}
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return "", errors.New("only http/https is allowed")
	}
	if u.Host == "" {
		return "", errors.New("host is mandatory")
	}
	return u.String(), nil
}
