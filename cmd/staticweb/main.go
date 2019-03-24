package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

var (
	letsencrypt  = flag.String("letsencrypt", "letsencrypt", "directory for letsencrypt cache")
	staticDir    = flag.String("site", "site", "directory served")
	httpAddr     = flag.String("http", ":80", "http address")
	httpsAddr    = flag.String("https", ":443", "https address")
	readTimeout  = flag.Duration("read_timeout", 5*time.Second, "read timeout")
	writeTimeout = flag.Duration("write_timeout", 5*time.Second, "write timeout")
	idleTimeout  = flag.Duration("idle_timeout", 120*time.Second, "idle timeout")
)

func main() {
	flag.Parse()

	autocert := autocert.Manager{
		Prompt: autocert.AcceptTOS,
		Cache:  autocert.DirCache(*letsencrypt),
	}

	httpsSrv := http.Server{
		Addr:         *httpsAddr,
		ReadTimeout:  *readTimeout,
		WriteTimeout: *writeTimeout,
		IdleTimeout:  *idleTimeout,
		Handler:      http.FileServer(http.Dir(*staticDir)),
		TLSConfig:    &tls.Config{GetCertificate: autocert.GetCertificate},
	}

	go func() {
		fmt.Printf("HTTPS listening on %s\n", httpsSrv.Addr)
		err := httpsSrv.ListenAndServeTLS("", "")
		if err != nil {
			log.Fatalf("HTTPs failed: %s", err)
		}
	}()

	httpSrv := http.Server{
		Addr:         *httpAddr,
		ReadTimeout:  *readTimeout,
		WriteTimeout: *writeTimeout,
		IdleTimeout:  *idleTimeout,
		Handler:      autocert.HTTPHandler(nil),
	}

	fmt.Printf("HTTP listening on %s\n", httpSrv.Addr)
	if err := httpSrv.ListenAndServe(); err != nil {
		log.Fatalf("HTTP failed: %s", err)
	}
}
