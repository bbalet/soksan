package main

/*  soksan allows you to interact with a go playground 
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
	"log"
	"net/http"
	"net/url"
	"os"
	"compress/gzip"
	"io"
	"strings"
	"fmt"
	"encoding/json"
	
	"bitbucket.org/kardianos/osext"
	"bitbucket.org/kardianos/service"	
	"github.com/gorilla/mux"
)

type Configuration struct {
	Port        	string
	Secured     	bool
	Compression 	bool
	HostPlayGround	string
	UserAgent		string
	SamplePath		string
}

var config Configuration
var CONFIGURATION_FILE string
var PRIVATE_KEY_FILE string
var CERTIFICATE_FILE string
var DATA_FOLDER string
var XP string

var logSrv service.Logger
var name = "soksan"
var displayName = "Go playground mediator"
var desc = "soksan allows you to interact with a go playground"
var isService bool = true
var store *sessions.FilesystemStore

// main runs the program as a service or as a command line tool.
// Several verbs allows you to install, start, stop or remove the service.
// "run" verb allows you to run the program as a command line tool.
// e.g. "goServerView install" installs the service
// e.g. "goServerView run" starts the program from the console (blocking)
func main() {
	s, err := service.NewService(name, displayName, desc)
	if err != nil {
		fmt.Printf("%s unable to start: %s", displayName, err)
		return
	}
	logSrv = s

	if len(os.Args) > 1 {
		var err error
		verb := os.Args[1]
		switch verb {
		case "install":
			err = s.Install()
			if err != nil {
				fmt.Printf("Failed to install: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" installed.\n", displayName)
		case "remove":
			err = s.Remove()
			if err != nil {
				fmt.Printf("Failed to remove: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" removed.\n", displayName)
		case "run":
			isService = false
			doWork()
		case "start":
			err = s.Start()
			if err != nil {
				fmt.Printf("Failed to start: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" started.\n", displayName)
		case "stop":
			err = s.Stop()
			if err != nil {
				fmt.Printf("Failed to stop: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" stopped.\n", displayName)
		}
		return
	}
	err = s.Run(func() error {
		// start
		go doWork()
		return nil
	}, func() error {
		// stop
		stopWork()
		return nil
	})
	if err != nil {
		s.Error(err.Error())
	}
}

// doWork is the actual main entry of the application whereas main set up
// the context (console program or service)
func doWork() {
	//Load configuration
	logInfo("Load configuration")
	XP, _ = osext.ExecutableFolder()
	CONFIGURATION_FILE = XP + "/conf/config.json"
	PRIVATE_KEY_FILE = XP + "/conf/private.pem"
	CERTIFICATE_FILE = XP + "/conf/cacert.pem"
	DATA_FOLDER = XP + "/data/"

	file, err := ioutil.ReadFile(CONFIGURATION_FILE)
	if err != nil {
		logFatal("(main) Configuration file : ", err)
	}
	json.Unmarshal(file, &config)

	//Start the embedded web server
	logInfo("Start the embedded web server")

	//Define the web application routes
	r := mux.NewRouter()
	r.HandleFunc("/fmt", makeHandler(fmtHandler))
	r.HandleFunc("/compile", makeHandler(compileShareHandler))
	r.HandleFunc("/run", makeHandler(runHandler))
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(XP)))
	http.Handle("/", r)

	logInfo("Listening on %s\n", config.Port)
	if config.Secured {
		err = http.ListenAndServeTLS(config.Port, CERTIFICATE_FILE, PRIVATE_KEY_FILE, nil)
		checkError(err)
	} else {
		err = http.ListenAndServe(config.Port, nil)
		checkError(err)
	}
}

// stopWork stops the service
func stopWork() {
	logInfo("I'm Stopping!")
}

//------------------------------------------------------------------------------
// Utility functions
//------------------------------------------------------------------------------

// logInfo reports a message in the console or the system log,
// depending on the execution context (console or service)
func logInfo(logMessage string, a ...interface{}) {
	if isService {
		logSrv.Info(logMessage, a...)
	} else {
		log.Printf(logMessage, a...)
	}
}

// logInfo reports an error in the console or the system log,
// depending on the execution context (console or service)
func logFatal(logMessage string, a ...interface{}) {
	if isService {
		logSrv.Error(logMessage, a...)
	} else {
		log.Fatalf(logMessage, a...)
	}
}

// checkError checks and reports any fatal error (errors occuring before the HTTP server is listening)
func checkError(err error) {
	if err != nil {
		logFatal("%v", err)
	}
}

// checkHttpError checks and reports any fatal error. Display an HTTP-500 page
func checkHttpError(err error, w http.ResponseWriter) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logFatal("%v", err)
	}
}

// makeHandler serves GZIP content
func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
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

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	if "" == w.Header().Get("Content-Type") {
		// If no content type, apply sniffing algorithm to un-gzipped body.
		w.Header().Set("Content-Type", http.DetectContentType(b))
	}
	return w.Writer.Write(b)
}


//------------------------------------------------------------------------------
// Web handlers
//------------------------------------------------------------------------------

// fmtHandler is the HTTP Handler of the formatting service
func fmtHandler(w http.ResponseWriter, r *http.Request) {
	
	
//	checkHttpError(err, w)
	w.Header().Set("Content-type", "application/json; charset: utf-8")
}

// compileShareHandler is the HTTP Handler of the compilation service
func compileShareHandler(w http.ResponseWriter, r *http.Request) {
	form := url.Values{}
	err := r.ParseForm()
	checkHttpError(err, w)
	for k, v := range r.Form {
		form.Add(k, v)
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", config.HostPlayGround + "/compile", nil)
	checkHttpError(err, w)
	req.Header.Set("User-Agent", config.UserAgent)
	resp, err := client.PostForm(config.HostPlayGround + "/compile", form)
	checkHttpError(err, w)
	
	//config.UserAgent
	//w.Header().Set("Content-type", "application/json; charset: utf-8")
	//resp, err := http.PostForm(config.HostPlayGround + "/compile", v)
	//checkHttpError(err, w)
	w.Header().Set("Content-type", "application/json; charset: utf-8")
	
	fmt.Fprint(r, resp)
}

// runHandler is the HTTP Handler of the compilation service variant
// it sends a go source file stored on the server instead of a web object content
func runHandler(w http.ResponseWriter, r *http.Request) {
	
	
	//checkHttpError(err, w)
	w.Header().Set("Content-type", "application/json; charset: utf-8")
}

