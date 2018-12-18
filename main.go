package main

import (
	"net/http"

	"github.com/spf13/viper"
)

func main() {
	viper.AutomaticEnv()
	httpHandler := NewHandler(true, viper.GetString("TEAMS_WEBHOOK_KEY"), webHook{})
	http.ListenAndServe(":80", http.HandlerFunc(httpHandler))
}

type webHook struct {
}

func (w webHook) OnMessage(req Request) (Response, error) {
	return BuildResponse("Hello " + req.FromUser.Name), nil
}
