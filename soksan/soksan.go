package soksan

/*  soksan allows you to embed a go playground in your website
    Copyright (C) 2014 Benjamin BALET

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.*/

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"compress/gzip"
	"path/filepath"
	"io"
	"strings"
	"fmt"
	"log"
	"os"
)

var (
	Compression 	bool
	HostPlayGround	string
	UserAgent		string
	SamplePath		string
)

// checkHttpError checks and reports any fatal error. Display an HTTP-500 page
func checkHttpError(err error, w http.ResponseWriter) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal("%v", err)
	}
}

// makeHandler serves GZIP content
func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// don't compress the answer if the client doesn't accept or if Compression disable
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") || !Compression {
			fn(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		fn(gzr, r)
	}
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

// Write is faster if you explicitly specify the content type
func (w gzipResponseWriter) Write(b []byte) (int, error) {
	if "" == w.Header().Get("Content-Type") {
		// If no content type, apply sniffing algorithm to un-gzipped body.
		w.Header().Set("Content-Type", http.DetectContentType(b))
	}
	return w.Writer.Write(b)
}

// requestPlayground requests the actual playground service and
// write back the response along with the HTTP status code
func requestPlayground(form url.Values, url string, w http.ResponseWriter) {
	
	//Reroute to playground and read the answer
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, nil)
	checkHttpError(err, w)
	req.Header.Set("User-Agent", UserAgent)
	resp, err := client.PostForm(url, form)
	checkHttpError(err, w)
	temp, _ := ioutil.ReadAll(resp.Body)
	
	if resp.StatusCode != 200 {
		http.Error(w, string(temp), http.StatusInternalServerError)
	} else {
		//Return answer with proper MIMETYPE
		w.Header().Set("Content-type", "application/json; charset: utf-8")
		fmt.Fprintf(w, "%s", temp)
	}
}

// compileHandler is the HTTP Handler of the compilation service
func compileHandler(w http.ResponseWriter, r *http.Request) {
	//Parse input form and insert these parameter into a new form
	form := url.Values{}
	err := r.ParseForm()
	checkHttpError(err, w)
	for k, v := range r.Form {
		form.Add(k, v[0])
	}
	requestPlayground(form, HostPlayGround + "/compile", w)
}

// fmtHandler is the HTTP Handler of the formatting service
func fmtHandler(w http.ResponseWriter, r *http.Request) {
	//Parse input form and insert these parameter into a new form
	form := url.Values{}
	err := r.ParseForm()
	checkHttpError(err, w)
	for k, v := range r.Form {
		form.Add(k, v[0])
	}
	requestPlayground(form, HostPlayGround + "/fmt", w)
}

// runHandler is the HTTP Handler of the compilation service variant
// it sends a go source file stored on the server instead of a web object content
func runHandler(w http.ResponseWriter, r *http.Request) {
	//Parse input form and get the content of the go file
	form := url.Values{}
	err := r.ParseForm()
	checkHttpError(err, w)
	
	//Check if file exists
	filename := filepath.Join(SamplePath, r.Form["file"][0])
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		http.Error(w, "Go snippet file not found", http.StatusNotFound)
		return
	}
	
	filego, err := ioutil.ReadFile(filename)
	checkHttpError(err, w)
	form.Add("version", "2")
	form.Add("body", string(filego))
	requestPlayground(form, HostPlayGround + "/compile", w)
}

// init initializes the package by registering the extra HTTP handlers
func init() {
	Compression = true
	HostPlayGround = "http://play.golang.org"
	UserAgent = "soksan"
	SamplePath = "./gocode/"
	http.HandleFunc("/fmt", makeHandler(fmtHandler))
	http.HandleFunc("/compile", makeHandler(compileHandler))
	http.HandleFunc("/run", makeHandler(runHandler))
}
