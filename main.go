package main

import (
	teams "github.com/ericdaugherty/msteams-webhook-go"
	"github.com/spf13/viper"
)

func main() {
	viper.AutomaticEnv()
	teams.NewHandler(true, viper.GetString("TEAMS_WEBHOOK_KEY"), webHook{})

}

type webHook struct {
}

func (w webHook) OnMessage(req teams.Request) (teams.Response, error) {
	return teams.BuildResponse("Hello " + req.FromUser.Name), nil
}
