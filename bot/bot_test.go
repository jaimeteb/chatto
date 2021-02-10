package bot

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/jaimeteb/chatto/channels"
	"github.com/jaimeteb/chatto/query"
)

var examples00TestPath = "../examples/00_test/"

func TestBot1(t *testing.T) {
	bot, err := LoadBot(&examples00TestPath)
	if err != nil {
		t.Errorf("failed to load bot: %s", err)
	}

	if bot.Name != "test_bot" {
		t.Errorf("bot name is incorrect, got: %v, want: %v.", bot.Name, "test_bot")
	}

	ans, err := bot.Answer(&query.Question{
		Sender: "bar",
		Text:   "on",
	})
	if err != nil {
		t.Errorf("failed to get answer from bot: %s", err)
	}

	if len(ans) != 1 && ans[0].Text != "Turning on." {
		t.Errorf("answer is incorrect, got: %v, want: %v.", ans, "Turning on.")
	}
}

func TestBot2(t *testing.T) {
	bot, err := LoadBot(&examples00TestPath)
	if err != nil {
		t.Fatalf("failed to load bot: %s", err)
	}

	_, err = bot.Answer(&query.Question{
		Sender: "baz",
		Text:   "on",
	})
	if err != nil {
		t.Fatal(err)
	}

	jsonStr := []byte(`{"sender": "42", "text": "on"}`)
	req, _ := http.NewRequest("POST", "", bytes.NewBuffer(jsonStr))
	w := httptest.NewRecorder()
	bot.restEndpointHandler(w, req)

	jsonStr2 := []byte(`{"update_id": 1, "message": {"message_id": 0, "from": 
	{"id": 42, "first_name": "", "username": ""}, "date": 0, "text": "off"}}`)
	req2, _ := http.NewRequest("POST", "", bytes.NewBuffer(jsonStr2))
	w2 := httptest.NewRecorder()
	bot.telegramEndpointHandler(w2, req2)

	formData := url.Values{
		"From":             {"42"},
		"Body":             {"?"},
		"To":               {""},
		"MediaUrl":         {""},
		"MediaContentType": {""},
		"MessageSid":       {""},
		"SmsStatus":        {""},
		"AccountSid":       {""},
		"Sid":              {""},
		"SmsSid":           {""},
		"SmsMessageSid":    {""},
		"NumMedia":         {"0"},
		"NumSegments":      {"0"},
		"ApiVersion":       {""},
	}
	req3, _ := http.NewRequest("POST", "", strings.NewReader(formData.Encode()))
	w3 := httptest.NewRecorder()
	bot.twilioEndpointHandler(w3, req3)

	req4, _ := http.NewRequest("GET", "/senders/42", nil)
	w4 := httptest.NewRecorder()
	bot.detailsHandler(w4, req4)

	jsonStr5 := []byte(`{"text": "."}`)
	req5, _ := http.NewRequest("POST", "", bytes.NewBuffer(jsonStr5))
	w5 := httptest.NewRecorder()
	bot.predictHandler(w5, req5)

	jsonStr6 := []byte(`{"event": {"channel": "43", "text": "on"}}`)
	req6, _ := http.NewRequest("POST", "", bytes.NewBuffer(jsonStr6))
	w6 := httptest.NewRecorder()
	bot.slackEndpointHandler(w6, req6)

	jsonStr7 := []byte(`{"challenge": "challenge"}`)
	req7, _ := http.NewRequest("POST", "", bytes.NewBuffer(jsonStr7))
	w7 := httptest.NewRecorder()
	bot.slackEndpointHandler(w7, req7)
}

// func TestBotNoClientsAndImages(t *testing.T) {
// 	path := "../examples/01_moodbot/"

// 	bot, err := LoadBot(&path)
// 	if err != nil {
// 		t.Errorf("failed to load bot: %s", err)
// 	}

// 	wREST := httptest.NewRecorder()
// 	messages := []fsm.Message{
// 		fsm.Message{
// 			Text: "only text",
// 		},
// 		fsm.Message{
// 			Text:  "text and image",
// 			Image: "https://i.imgur.com/8MU0IUT.jpeg",
// 		},
// 	}

// 	receive, err := bot.Channels.REST.ReceiveMessage(wREST, wREST.Result().Request)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	bot.Channels.REST.SendMessage()
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

func TestServeBot(t *testing.T) {
	port := 9999

	go ServeBot(&examples00TestPath, &port)
}

func TestExtFromBot(t *testing.T) {
	bot, err := LoadBot(&examples00TestPath)
	if err != nil {
		t.Fatalf("failed to load bot: %s", err)
	}

	bot.Channels = &channels.Channels{}

	_, err = bot.Answer(&query.Question{
		Sender: "ext_tester",
		Text:   "hello",
	})
	if err != nil {
		t.Fatal(err)
	}
}
