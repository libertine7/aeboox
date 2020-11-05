package main

import (
	"crypto/tls"
	"github.com/allegro/bigcache"
	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json2"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
	"time"
)

func basicAuthApi(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()

		if !ok || !(user == "x" && pass == "3") {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized.", 401)
			return
		}

		h.ServeHTTP(w, r)
	})
}

var Cache *bigcache.BigCache

func main() {
	Cache, _ = bigcache.NewBigCache(bigcache.Config{
		Shards:             1024,
		LifeWindow:         48 * time.Hour,
		MaxEntriesInWindow: 1000 * 10 * 60,
		MaxEntrySize:       500,
		Verbose:            false,
		HardMaxCacheSize:   512,
		OnRemove:           nil,
		OnRemoveWithReason: nil,
	})

	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("api.apps4reading.com"),
		Cache:      autocert.DirCache("certs"),
	}

	apiServer := rpc.NewServer()
	apiServer.RegisterCodec(json2.NewCustomCodec(&rpc.CompressionSelector{}), "application/json")
	apiServer.RegisterService(new(EbooxAnalytics), "")
	http.Handle("/jsonrpc/", basicAuthApi(apiServer))

	server := &http.Server{Addr: ":https", TLSConfig: &tls.Config{GetCertificate: certManager.GetCertificate}}
	go http.ListenAndServe(":http", certManager.HTTPHandler(nil))
	log.Fatal(server.ListenAndServeTLS("", "")) //Key and cert are coming from Let's Encrypt
}
