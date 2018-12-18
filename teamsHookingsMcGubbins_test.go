package main

import (
	"encoding/base64"
	"testing"
)

func Test_authenticateRequest(t *testing.T) {

	standardKey, _ := base64.StdEncoding.DecodeString("n6zPd4AYcZjQEfn8P8s9zP7Mc+2FrYI0Lzg6NszHNP8=")

	type args struct {
		body       []byte
		authHeader string
		key        []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Actual example",
			args{
				body:       []byte("{\"type\":\"message\",\"id\":\"1545105672086\",\"timestamp\":\"2018-12-18T04:01:12.088Z\",\"localTimestamp\":\"2018-12-18T04:01:12.088+00:00\",\"serviceUrl\":\"https://smba.trafficmanager.net/uk/\",\"channelId\":\"msteams\",\"from\":{\"id\":\"29:1xOKsq8hwMFMu1_gm4WD3sqdYZ-5CRjT1VKlFPY6RFDd4jSDW08nMwBDhlpl1kUFRL2s0CWEW2vOatviT5iH_yQ\",\"name\":\"Gary Gary\",\"aadObjectId\":\"06928b7d-cec2-4ddb-95b6-9014023bcfd5\"},\"conversation\":{\"isGroup\":true,\"id\":\"19:b4be35a4a5b944b3bcc86bb020f0367f@thread.skype;messageid=1545105672086\",\"name\":null,\"conversationType\":\"channel\"},\"recipient\":null,\"textFormat\":\"plain\",\"attachmentLayout\":null,\"membersAdded\":[],\"membersRemoved\":[],\"topicName\":null,\"historyDisclosed\":null,\"locale\":null,\"text\":\"<at>Hooker</at>\\n\",\"speak\":null,\"inputHint\":null,\"summary\":null,\"suggestedActions\":null,\"attachments\":[{\"contentType\":\"text/html\",\"contentUrl\":null,\"content\":\"<div><div><span itemscope=\\\"\\\" itemtype=\\\"http://schema.skype.com/Mention\\\" itemid=\\\"0\\\">Hooker</span></div>\\n</div>\",\"name\":null,\"thumbnailUrl\":null}],\"entities\":[{\"type\":\"clientInfo\",\"locale\":\"en-GB\",\"country\":\"GB\",\"platform\":\"Windows\"}],\"channelData\":{\"teamsChannelId\":\"19:b4be35a4a5b944b3bcc86bb020f0367f@thread.skype\",\"teamsTeamId\":\"19:b4be35a4a5b944b3bcc86bb020f0367f@thread.skype\",\"channel\":{\"id\":\"19:b4be35a4a5b944b3bcc86bb020f0367f@thread.skype\"},\"team\":{\"id\":\"19:b4be35a4a5b944b3bcc86bb020f0367f@thread.skype\"},\"tenant\":{\"id\":\"03d27877-9753-4c42-a088-b3c1ca2ed391\"}},\"action\":null,\"replyToId\":null,\"value\":null,\"name\":null,\"relatesTo\":null,\"code\":null}"),
				authHeader: "5bQV8u3oy+V8LL1T2XmrnxMMLyo8hDSBvGrLOXPpuzg=",
				key:        standardKey,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := authenticateRequest(tt.args.body, tt.args.authHeader, tt.args.key); got != tt.want {
				t.Errorf("authenticateRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
