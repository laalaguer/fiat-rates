package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func newRouter(proxyAndPort string, adminPassword string, accessKey string) *mux.Router {
	r := mux.NewRouter()

	// // Get the hello handler
	// r.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintf(w, "Hello, %q", html.EscapeString("world"))
	// }).Methods("GET")

	// Fiat Exchange Rate Getter
	r.HandleFunc("/fiat", func(w http.ResponseWriter, r *http.Request) {
		FiatPriceHandler(w, r, proxyAndPort)
	}).Methods("GET")

	// Fiat Exchange Rate Refresh from upstream
	r.HandleFunc("/refresh", func(w http.ResponseWriter, r *http.Request) {
		RefreshPriceHandler(w, r, proxyAndPort, adminPassword, accessKey)
	}).Methods("GET")

	// Get all supported symbols and bases
	r.HandleFunc("/supported", func(w http.ResponseWriter, r *http.Request) {
		SupportedSymbolsHandler(w, r)
	})

	// r.HandleFunc("/cache/debug", func(w http.ResponseWriter, r *http.Request) {
	// 	PriceCacheHandler(w, r, proxyAndPort)
	// }).Methods("GET")

	return r
}

func main() {
	var adminPassword string
	var accessKey string
	var ipAndPort string
	var proxyAndPort string
	flag.StringVar(&adminPassword, "password", "", "Admin password to control the service")
	flag.StringVar(&accessKey, "key", "", "Access key to call exchangeratesapi")
	flag.StringVar(&ipAndPort, "listen", "", "Local IP and port this program runs on")
	flag.StringVar(&proxyAndPort, "proxy", "", "Proxy IP and port for outgoing requests")
	flag.Parse()

	if len(ipAndPort) == 0 || len(accessKey) == 0 || len(adminPassword) == 0 {
		fmt.Println("Usage: -password [admin password] -key [access key] -listen [ip:port] -proxy [proxy:port](optional)")
		flag.PrintDefaults()
		os.Exit(1)
	}

	r := newRouter(proxyAndPort, adminPassword, accessKey)
	http.ListenAndServe(ipAndPort, r)
}
