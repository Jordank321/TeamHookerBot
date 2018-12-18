package commands

import (
	"log"

	wit "github.com/jsgoecke/go-wit"
)

var client *wit.Client

func SetWitToken(token string) {
	client = wit.NewClient(token)
}

func GetWitIntent(msg string) []wit.Outcome {
	// Process a text message
	request := &wit.MessageRequest{}
	request.Query = msg
	result, err := client.Message(request)
	if err != nil {
		log.Printf(err.Error())
	}
	return result.Outcomes
}
