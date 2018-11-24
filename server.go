package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"gowebapp/source/controller"
	"gowebapp/source/shared/database"
	"gowebapp/source/shared/email"
	"gowebapp/source/view"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

const (
	// DEV represents developer mode
	// when set to true, certain security features (such as https, csrf)
	// will be disabled
	DEV = true
)

func cleanup(semaphore chan os.Signal) {
	s := <-semaphore
	if s == os.Interrupt {
		log.Println("Stopping server: user interrrupt")
	} else {
		log.Println("Stopping server: process terminated")
	}
	// Stop the server(s)
	if err := httpSrv.Shutdown(context.Background()); err != nil {
		log.Println(err)
	}
	if DEV == false {
		if err := tlsSrv.Shutdown(context.Background()); err != nil {
			log.Println(err)
		}
	}
	// Close database
	if err := database.Terminate(); err != nil {
		log.Println(err)
	}
	os.Exit(1)
}

var (
	tlsSrv  *http.Server
	httpSrv *http.Server
)

func initTLSServer() {
	m := autocert.Manager{
		Cache:      autocert.DirCache("secret-dir"),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("legoracer.ddns.net"),
	}
	r := controller.GetRouter()
	tlsConf := &tls.Config{
		// Causes servers to use Go's default ciphersuite preferences,
		// which are tuned to avoid attacks. Does nothing on clients.
		PreferServerCipherSuites: true,
		// Only use safe curves which have assembly implementations
		CurvePreferences: []tls.CurveID{
			tls.X25519,
		},
		GetCertificate: m.GetCertificate,
		MinVersion:     tls.VersionTLS11,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305, // Go 1.8 only
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,   // Go 1.8 only
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256, // iOS compatibility

			// Best disabled, as they don't provide Forward Secrecy,
			// but might be necessary for some clients
			// tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			// tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		},
	}
	tlsSrv = &http.Server{
		Handler:   r,
		Addr:      ":443",
		TLSConfig: tlsConf,
		// enforce timeouts
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}

func initRedirectServer() {
	var r http.Handler
	// In dev mode all pages are available through HTTP
	if DEV {
		r = controller.GetRouter()
	} else {
		r = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Connection", "close")
			url := "https://" + r.Host + r.URL.String()
			http.Redirect(w, r, url, http.StatusMovedPermanently)
		})
	}

	httpSrv = &http.Server{
		Handler: r,
		Addr:    ":8080",
		// enforce timeouts
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

}

func jsonMustUnmarshal(path string, dest interface{}) {
	// Open file
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	// Unmarshal json data into provided interface
	err = json.Unmarshal(data, dest)
	if err != nil {
		log.Fatal(err)
	}
}

func databaseConfig() *database.MySQLInfo {
	dbinfo := &database.MySQLInfo{}
	jsonMustUnmarshal("config/dbconfig.json", dbinfo)
	return dbinfo
}

func emailConfig() *email.SMTPInfo {
	e := &email.SMTPInfo{}
	jsonMustUnmarshal("config/smtp.json", e)
	return e
}

func main() {
	// Initialize views
	view.LoadTemplates("template")
	// Initialize DB
	database.Connect(databaseConfig())
	// Initialize emails
	email.Configure(emailConfig())
	email.LoadTemplates("email")

	if DEV == false {
		initTLSServer()
		go func() {
			log.Fatal(tlsSrv.ListenAndServeTLS("", ""))
		}()
	}

	initRedirectServer()
	go func() {
		log.Fatal(httpSrv.ListenAndServe())
	}()

	log.Println("Server started")

	// capture process termination signalls
	csig := make(chan os.Signal, 2)
	signal.Notify(csig, os.Interrupt, syscall.SIGTERM)
	cleanup(csig)
}
