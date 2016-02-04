package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
)

func printRaw(req *http.Request) {
	raw, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Println("DumpRequest: ", err)
		return
	}
	fmt.Printf("%s\n", raw)
}

type Proxy struct {
	Client http.Client
}

// See https://golang.org/src/net/http/httputil/reverseproxy.go
func (p *Proxy) Serve(rw http.ResponseWriter, in *http.Request) {
	out := new(http.Request)
	*out = *in // Shallow copy apparently

	// It's an error for RequestURI to be set in a client request
	out.RequestURI = ""
	out.URL.Scheme = "http"
	out.URL.Host = in.Host // From the Host header

	// XXX Refuse localhost & co.

	out.URL.Path = in.URL.Path

	log.Println("Forwarding: ", out.URL)

	// XXX Query string

	resp, err := p.Client.Do(out)
	if err != nil {
		log.Println("Do: ", err)
		return
	}

	var rh http.Header = rw.Header()
	for key, values := range resp.Header {
		for _, v := range values {
			rh.Add(key, v)
		}
	}

	rw.WriteHeader(resp.StatusCode)

	b := make([]byte, 16)
	defer resp.Body.Close()
	for {
		n, err := resp.Body.Read(b)
		rw.Write(b[:n])
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println("Body.Read: ", err)
			return
		}
	}
}

func main() {
	p := new(Proxy)
	http.HandleFunc("/", p.Serve)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
