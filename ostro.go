package main

import (
	"fmt"
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

type OstrogoProxy struct {
	Client http.Client
}

// See https://golang.org/src/net/http/httputil/reverseproxy.go
func (p *OstrogoProxy) Serve(rw http.ResponseWriter, in *http.Request) {
	out := new(http.Request)
	*out = *in // Shallow copy apparently

	// It's an error for RequestURI to be set in a client request
	out.RequestURI = ""
	out.URL.Scheme = "http"
	out.URL.Host = in.Host // From the Host header
	out.URL.Path = in.URL.Path

	log.Println("Forwarding: ", out.URL)

	// XXX Query string

	resp, err := p.Client.Do(out)
	if err != nil {
		log.Println("Do: ", err)
		return
	}
	fmt.Println(resp)
	// rw.WriteHeader(resp.Status)
	// b, err := resp.Body.Read()
	// if err != nil {
	//	log.Println("Do: ", err)
	//	return
	// }
	resp.Write(rw)
}

func main() {
	p := new(OstrogoProxy)
	http.HandleFunc("/", p.Serve)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
