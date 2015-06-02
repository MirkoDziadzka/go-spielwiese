package main

import (
	"log"
	"net/http"
)
import "net/http/httputil"
import "crypto/tls"

// return a list of 'secure' cipher suites ordered by preference
func cipherSuites() []uint16 {
	return []uint16{
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,

		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,

		tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_RSA_WITH_AES_128_CBC_SHA,
		tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
	}
}

// create my own proxy object ... add mappings and handle teh whole stuff later
type Proxy map[string]*httputil.ReverseProxy

// backend is the ip or the hostname to which we proxy. But we
// assume that this host accept requests for the original hostname
func (p *Proxy) add(hostname, backend string) {
	proxy := &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.Host = backend
			r.URL.Scheme = "http"
			r.URL.Host = backend
			log.Printf("proxy: %v\n", r)
			log.Printf("proxy: %v\n", r.URL)
		},
	}
	(*p)[hostname] = proxy
}

func (p *Proxy) handler(w http.ResponseWriter, r *http.Request) {
	proxy, ok := (*p)[r.Host]
	if ok {
		proxy.ServeHTTP(w, r)
	} else {
		http.Error(w, "host not found", http.StatusNotFound)
	}
}

func main() {
	proxy := make(Proxy)
	proxy.add("bithalde.de", "127.0.0.1")
	proxy.add("localhost:8080", "bithalde.de")

	// cert stuff
	//cert_srv3_bithalde, err := tls.loadx509keypair("/home/mirko/srv3.bithalde.de.crt", "/home/mirko/srv3.bithalde.de.key")
	//if err != nil {
	//	log.fatal(err)
	//}
	//config := &tls.config{
	//	certificates:             []tls.certificate{cert_srv3_bithalde},
	//	ciphersuites:             ciphersuites(),
	//	minversion:               tls.versiontls10,
	//	preferserverciphersuites: true,
	//}
	//config.buildnametocertificate()

	//for key := range config.nametocertificate {
	//	log.printf("key: %s\n", key)
	//}
	//listener, err := tls.listen("tcp", ":443", config)

	//http.Serve(listener, nil)
	http.HandleFunc("/", proxy.handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
