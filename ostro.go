package main

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
)

type certIssuer struct {
	CACert *tls.Certificate
}

func (ci *certIssuer) getCertificate(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	log.Printf("getCertificate called.")
	// derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
	return nil, nil
}

// Looks like http.ListenAndServeTLS() but here the certificate & the private
// key are going to be used to issue new certificates as they are requested.
func ListenAndServeTLS(addr string, certFile string, keyFile string, handler http.Handler) error {

	server := http.Server{Addr: addr, Handler: handler}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	var issuer *certIssuer
	if certFile != "" || keyFile != "" {
		caCert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return err
		}
		issuer = &certIssuer{CACert: &caCert}
	}

	config := &tls.Config{GetCertificate: issuer.getCertificate}
	tlsListener := tls.NewListener(ln.(*net.TCPListener), config)

	return server.Serve(tlsListener)
}

func handler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server.\n"))
}

func main() {
	http.HandleFunc("/", handler)
	log.Printf("About to listen on 10443. Go to https://127.0.0.1:10443/")
	err := ListenAndServeTLS(":10443", "cert_test.pem", "key_test.pem", nil)
	if err != nil {
		log.Fatal(err)
	}
}
