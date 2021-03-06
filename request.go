package request

import (
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"
)

const Version = "0.2.0"

type FileField struct {
	FieldName string
	FileName  string
	File      io.Reader
}

type BasicAuth struct {
	Username string
	Password string
}

type Args struct {
	Client    *http.Client
	Headers   map[string]string
	Cookies   map[string]string
	Data      map[string]string
	Params    map[string]string
	Files     []FileField
	Json      interface{}
	Proxy     string
	BasicAuth BasicAuth
}

func NewArgs(c *http.Client) *Args {
	if c.Jar == nil {
		options := cookiejar.Options{
			PublicSuffixList: publicsuffix.List,
		}
		jar, _ := cookiejar.New(&options)
		c.Jar = jar
	}

	return &Args{
		Client:    c,
		Headers:   defaultHeaders,
		Cookies:   nil,
		Data:      nil,
		Params:    nil,
		Files:     nil,
		Json:      nil,
		Proxy:     "",
		BasicAuth: BasicAuth{},
	}
}

func newURL(u string, params map[string]string) string {
	if params == nil {
		return u
	}

	p := url.Values{}
	for k, v := range params {
		p.Set(k, v)
	}
	if strings.Contains(u, "?") {
		return u + "&" + p.Encode()
	}
	return u + "?" + p.Encode()
}

func newBody(a *Args) (body io.Reader, contentType string, err error) {
	if a.Data == nil && a.Files == nil && a.Json == nil {
		return nil, "", nil
	}
	if a.Files != nil {
		return newMultipartBody(a)
	} else if a.Json != nil {
		return newJsonBody(a)
	}

	d := url.Values{}
	for k, v := range a.Data {
		d.Set(k, v)
	}
	return strings.NewReader(d.Encode()), "", nil
}

func newRequest(method string, url string, a *Args) (resp *Response, err error) {
	body, contentType, err := newBody(a)
	u := newURL(url, a.Params)
	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return nil, err
	}
	applyHeaders(a, req, contentType)
	applyCookies(a, req)
	applyProxy(a)
	applyCheckRdirect(a)

	if a.BasicAuth.Username != "" {
		req.SetBasicAuth(a.BasicAuth.Username, a.BasicAuth.Password)
	}

	s, err := a.Client.Do(req)
	resp = &Response{s, nil}
	return
}

// Get issues a GET to the specified URL.
//
// Caller should close resp.Body when done reading from it.
func Get(url string, a *Args) (resp *Response, err error) {
	resp, err = newRequest("GET", url, a)
	return
}

// Head issues a HEAD to the specified URL.
//
// Caller should close resp.Body when done reading from it.
func Head(url string, a *Args) (resp *Response, err error) {
	resp, err = newRequest("HEAD", url, a)
	return
}

// Post issues a POST to the specified URL.
//
// Caller should close resp.Body when done reading from it.
func Post(url string, a *Args) (resp *Response, err error) {
	resp, err = newRequest("POST", url, a)
	return
}

// Put issues a PUT to the specified URL.
//
// Caller should close resp.Body when done reading from it.
func Put(url string, a *Args) (resp *Response, err error) {
	resp, err = newRequest("PUT", url, a)
	return
}

// Patch issues a PATCH to the specified URL.
//
// Caller should close resp.Body when done reading from it.
func Patch(url string, a *Args) (resp *Response, err error) {
	resp, err = newRequest("PATCH", url, a)
	return
}

// Delete issues a DELETE to the specified URL.
//
// Caller should close resp.Body when done reading from it.
func Delete(url string, a *Args) (resp *Response, err error) {
	resp, err = newRequest("DELETE", url, a)
	return
}

// Options issues a OPTIONS to the specified URL.
//
// Caller should close resp.Body when done reading from it.
func Options(url string, a *Args) (resp *Response, err error) {
	resp, err = newRequest("OPTIONS", url, a)
	return
}
