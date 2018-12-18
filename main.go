package main

import (
	"crypto/tls"
	"net/http"

	"github.com/spf13/viper"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	viper.AutomaticEnv()

	mux := http.NewServeMux()
	httpHandler := NewHandler(true, viper.GetString("TEAMS_WEBHOOK_KEY"), webHook{})
	mux.HandleFunc("/", httpHandler)

	certManager := autocert.Manager{
		Prompt: autocert.AcceptTOS,
		Cache:  autocert.DirCache("certs"),
	}

	server := &http.Server{
		Addr:    ":443",
		Handler: mux,
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}

	go http.ListenAndServe(":80", certManager.HTTPHandler(nil))
	server.ListenAndServeTLS("", "")
}

type webHook struct {
}

func (w webHook) OnMessage(req Request) (Response, error) {
	return BuildResponse("Hello " + req.FromUser.Name), nil
}
