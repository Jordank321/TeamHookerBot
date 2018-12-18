package main

import (
	"log"
	"net/http"

	"golang.org/x/crypto/acme/autocert"

	"github.com/spf13/viper"
)

func main() {
	viper.AutomaticEnv()
	key := viper.GetString("TEAMS_WEBHOOK_KEY")

	if len(key) == 0 {
		log.Panic("It helps when you set a key in the env as TEAMS_WEBHOOK_KEY")
	}

	mux := http.NewServeMux()
	httpHandler := NewHandler(true, key, webHook{})
	mux.HandleFunc("/", httpHandler)

	certManager := autocert.Manager{
		Prompt: autocert.AcceptTOS,
		Cache:  autocert.DirCache("certs"),
		//HostPolicy: autocert.HostWhitelist("jordankelwick.com", "www.jordankelwick.com"),
	}

	server := &http.Server{
		Addr:      ":https",
		Handler:   mux,
		TLSConfig: certManager.TLSConfig(),
	}

	go http.ListenAndServe(":80", certManager.HTTPHandler(nil))
	err := server.ListenAndServeTLS("", "")
	panicErr(err)
}

type webHook struct {
}

func (w webHook) OnMessage(req Request) (Response, error) {
	return BuildTextResponse("Hello " + req.FromUser.Name), nil
}
