package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jsgoecke/go-wit"

	"github.com/Jordank321/TeamHookerBot/commands"

	"golang.org/x/crypto/acme/autocert"

	. "github.com/ahmetb/go-linq"
	"github.com/peterhellberg/giphy"
	"github.com/spf13/viper"
)

func (w webHook) OnMessage(req Request) (Response, error) {
	intents := commands.GetWitIntent(req.Text)

	data, _ := json.MarshalIndent(intents, "", "    ")
	log.Println(string(data[:]))

	if len(intents) == 0 {
		return BuildTextResponse("Hello " + req.FromUser.Name), nil
	}
	From(intents).Sort(func(i interface{}, j interface{}) bool {
		return i.(wit.Outcome).Confidence < j.(wit.Outcome).Confidence
	}).ToSlice(&intents)

	bestGuess := intents[0]

	data, _ = json.MarshalIndent(bestGuess, "", "    ")
	log.Println(string(data[:]))

	hasRandomGifRequest := From(bestGuess.Entities).AnyWith(func(entity interface{}) bool {
		if entity.(KeyValue).Key.(string) == "intent" {
			return From(entity.(KeyValue).Value.([]wit.MessageEntity)).AnyWith(func(messageEntity interface{}) bool {
				return (*messageEntity.(wit.MessageEntity).Value).(string) == "GIF_RANDOM"
			})
		}
		return false
	})
	if hasRandomGifRequest {
		g := giphy.DefaultClient
		resp, err := g.Random(nil)
		panicErr(err)
		return BuildTextResponse(resp.Data.Caption + " - " + resp.Data.URL), nil
	}
	return BuildTextResponse("Idk, want some cookies?"), nil
}

func main() {
	viper.AutomaticEnv()

	teamsKey := getCheckSetting("TEAMS_WEBHOOK_KEY")
	//giphyKey := getCheckSetting("GIPHY_API_KEY")
	witAiKey := getCheckSetting("WIT_AI_KEY")

	commands.SetWitToken(witAiKey)

	mux := http.NewServeMux()
	httpHandler := NewHandler(true, teamsKey, webHook{})
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
	//err := server.ListenAndServe()
	panicErr(err)
}

type webHook struct {
}

func getCheckSetting(name string) string {
	value := viper.GetString(name)
	if len(value) == 0 {
		log.Panicf("It helps when you set a key in the env as %s", name)
	}
	return value
}
