package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var auth bool
var keyBytes []byte
var webhook WebHook

// WebHook represnts the interface needed to handle Microsoft Teams WebHook Requests.
type WebHook interface {
	OnMessage(Request) (Response, error)
}

// Request data representing an inbound WebHook request from Microsoft Teams.
type Request struct {
	Type           string `json:"type"`
	ID             string `json:"id"`
	Timestamp      string `json:"timestamp"`
	LocalTimestamp string `json:"localTimestamp"`
	ServiceURL     string `json:"serviceUrl"`
	ChannelID      string `json:"channelId"`
	FromUser       User   `json:"from"`
	Conversation   struct {
		ID string `json:"id"`
	} `json:"conversation"`
	RecipientUser User   `json:"recipient"`
	TextFormat    string `json:"textFormat"`
	Text          string `json:"text"`
	Attachments   []struct {
		ContentType string `json:"contentType"`
		Content     string `json:"Content"`
	} `json:"attachments"`
	Entities    []interface{} `json:"entities"`
	ChannelData struct {
		TeamsChannelID string `json:"teamsChannelId"`
		TeamsTeamID    string `json:"teamsTeamId"`
	}
}

// Response represents the data to return to Microsoft Teams.
type Response struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// User represents data for a Microsoft Teams user.
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// NewHandler initializes and returns a Lambda handler to process incoming requests.
func NewHandler(authenticateRequests bool, key string, wh WebHook) func(w http.ResponseWriter, r *http.Request) {
	auth = authenticateRequests
	keyBytes, _ = base64.StdEncoding.DecodeString(key)
	webhook = wh
	return handler
}

func handler(w http.ResponseWriter, lreq *http.Request) {

	if auth {
		authenticated := authenticateRequest(lreq)
		if !authenticated {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Invalid Authentication"))
			return
		}
	}

	var treq = Request{}
	bodyBytes, err := ioutil.ReadAll(lreq.Body)
	panicErr(err)

	log.Print("Request body: " + string(bodyBytes))
	err = json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&treq)
	panicErr(err)

	tresp, err := webhook.OnMessage(treq)
	panicErr(err)

	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(tresp)
	panicErr(err)

	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}

func authenticateRequest(lreq *http.Request) bool {
	bodyBytes, err := ioutil.ReadAll(lreq.Body)
	panicErr(err)
	authHeader := lreq.Header.Get("Authorization")
	messageMAC, _ := base64.StdEncoding.DecodeString(strings.TrimPrefix(authHeader, "HMAC "))

	mac := hmac.New(sha256.New, keyBytes)
	mac.Write(bodyBytes)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}

// BuildResponse is a helper method to build a Response
func BuildResponse(text string) Response {
	return Response{
		Type: "message",
		Text: text,
	}
}

func panicErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
