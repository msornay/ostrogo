package main

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
)

type OServer struct {
	*http.Server
}

// ListenAndServeTLS() of http.Server configure one certificate in a tls.Config
// from a given certFile. We want to override that. Instead, we need a filename
// containing a private key for issuing certificates.
func (srv *OServer) ListenAndServeTLS(keyFile string) error {
	ln, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		return err
	}
	config := new(tls.Config)
	tlsListener := tls.NewListener(ln.(*net.TCPListener), config)
	return srv.Serve(tlsListener)
}

func ListenAndServeTLS(addr string, keyFile string, handler http.Handler) error {
	server := &OServer{&http.Server{Addr: addr, Handler: handler}}
	return server.ListenAndServeTLS(keyFile)
}

func handler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server.\n"))
}

func main() {
	http.HandleFunc("/", handler)
	log.Printf("About to listen on 10443. Go to https://127.0.0.1:10443/")
	err := ListenAndServeTLS(":10443", "key.pem", nil)
	if err != nil {
		log.Fatal(err)
	}
}
