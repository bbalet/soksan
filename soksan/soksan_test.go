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
	"testing"
	"net/http/httptest"
)

func TestCompilation(t *testing.T) {
	//setup lib and test handler
	UserAgent = "soksan-test"
	ts := httptest.NewServer(http.HandlerFunc(compileHandler))
	defer ts.Close()
	//instanciate client and prepare POST form
	client := &http.Client{}
	form := url.Values{}
	form.Add("version", "2")
	form.Add("body", goCodeHello)
	//request and get the BODY element of the answer
	resp, err := client.PostForm(ts.URL, form)
	if err != nil { panic(err) }
	temp, err := ioutil.ReadAll(resp.Body)
	if err != nil { panic(err) }
	body := string(temp)
	
	if body != goCodeHelloResult {
		t.Errorf("#0 result = %s ; expected %s", body, goCodeHelloResult)
	}
}

func TestFormat(t *testing.T) {
	//setup lib and test handler
	UserAgent = "soksan-test"
	ts := httptest.NewServer(http.HandlerFunc(fmtHandler))
	defer ts.Close()
	//instanciate client and prepare POST form
	client := &http.Client{}
	form := url.Values{}
	form.Add("version", "2")
	form.Add("body", goCodeHello)
	//request and get the BODY element of the answer
	resp, err := client.PostForm(ts.URL, form)
	if err != nil { panic(err) }
	temp, err := ioutil.ReadAll(resp.Body)
	if err != nil { panic(err) }
	body := string(temp)
	
	if body != goCodeHelloFormat {
		t.Errorf("#0 result = %s ; expected %s", body, goCodeHelloFormat)
	}
}

func TestContentType(t *testing.T) {
	//setup lib and test handlers
	UserAgent = "soksan-test"
	ts1 := httptest.NewServer(http.HandlerFunc(compileHandler))
	ts2 := httptest.NewServer(http.HandlerFunc(fmtHandler))
	defer ts1.Close()
	defer ts2.Close()
	//instanciate client and prepare POST form (used 2 times)
	client := &http.Client{}
	form := url.Values{}
	form.Add("version", "2")
	form.Add("body", goCodeHello)
	
	//request and test the HEADER of each answers
	resp, err := client.PostForm(ts1.URL, form)
	if err != nil { panic(err) }
	if resp.Header.Get("Content-Type") != "application/json; charset: utf-8" {
		t.Errorf("#0 HEADER / Content-Type = %s; expected %s", resp.Header.Get("Content-Type"), "application/json; charset: utf-8")
	}
	resp, err = client.PostForm(ts2.URL, form)
	if err != nil { panic(err) }
	if resp.Header.Get("Content-Type") != "application/json; charset: utf-8" {
		t.Errorf("#0 HEADER / Content-Type = %s; expected %s", resp.Header.Get("Content-Type"), "application/json; charset: utf-8")
	}
}

var goCodeHello = `package main
import "fmt"
func main() {
	fmt.Println("Hello, playground")
}`

var goCodeHelloResult = "{\"Errors\":\"\",\"Events\":[{\"Message\":\"Hello, playground\\n\",\"Delay\":0}]}\n"

var goCodeHelloFormat = `{"Body":"package main\n\nimport \"fmt\"\n\nfunc main() {\n\u0009fmt.Println(\"Hello, playground\")\n}\n","Error":""}
`
/*type TestGZip struct {
	Compression		bool
	AcceptEncoding	string
	Compressed		string
}

var testsGZip = []TestGZip{
    {true, "gzip", "gzip"},
    {true, "", ""},
	{false, "gzip", ""},
    {false, "", ""},
}


func TestCompressionOption(t *testing.T) {
	//setup lib and test handler
	UserAgent = "soksan-test"
	ts := httptest.NewServer(makeHandler(compileHandler))
	defer ts.Close()

    for i, test := range testsGZip {
		client := &http.Client{}
		req, err := http.NewRequest("POST", ts.URL, nil)
		if err != nil { panic(err) }
		req.Header.Set("Accept-Encoding", test.AcceptEncoding)	
		form := url.Values{}
		form.Add("version", "2")
		form.Add("body", goCodeHello)
		resp, err := client.PostForm(ts.URL, form)
		if err != nil { panic(err) }
		
		//Content-Encoding:gzip	
        if resp.Header.Get("Content-Encoding") != test.Compressed {
			t.Errorf("#%v", req.Header)
			t.Errorf("#%v", resp.Header)
			temp, _ := ioutil.ReadAll(resp.Body)
			t.Errorf("#%v", string(temp))
            t.Errorf("#%d: (%t, %s)=%s; expected %s", i, test.Compression, test.AcceptEncoding, resp.Header["Content-Encoding"], test.Compressed)
        }
    }
}*/