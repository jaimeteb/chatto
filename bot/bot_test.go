package bot_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jaimeteb/chatto/bot"
	"github.com/jaimeteb/chatto/channels"
	"github.com/jaimeteb/chatto/channels/messages"
	"github.com/jaimeteb/chatto/channels/mockchannels"
	"github.com/jaimeteb/chatto/clf"
	"github.com/jaimeteb/chatto/extension"
	"github.com/jaimeteb/chatto/fsm"
	"github.com/jaimeteb/chatto/query"
	"github.com/jaimeteb/chatto/testutils"
	log "github.com/sirupsen/logrus"
)

func TestBot_endpointHandler(t *testing.T) {
	testBot, restChnl, twilioChnl, telegramChnl, slackChnl, err := newTestBot(t)
	if err != nil {
		t.Fatal(err)
	}

	ts := httptest.NewServer(testBot.Router)
	defer ts.Close()

	type args struct {
		endpoint    string
		message     []byte
		mockReceive *gomock.Call
		mockSend    *gomock.Call
	}
	tests := []struct {
		name    string
		bot     *bot.Bot
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "rest endpoint test",
			bot:  testBot,
			args: args{
				endpoint:    fmt.Sprintf("%s/endpoints/rest", ts.URL),
				message:     []byte(`{"sender": "42", "text": "on"}`),
				mockReceive: restChnl.EXPECT().ReceiveMessage(gomock.Any()).Return(&messages.Receive{Question: &query.Question{Sender: "42", Text: "on"}}, nil),
				mockSend:    restChnl.EXPECT().SendMessage(gomock.Any()).Return(nil),
			},
			want: `[{"text":"Turning on.","image":""}]`,
		},
		{
			name: "twilio endpoint test",
			bot:  testBot,
			args: args{
				endpoint:    fmt.Sprintf("%s/endpoints/twilio", ts.URL),
				message:     []byte(`{"sender": "42", "text": "off"}`),
				mockReceive: twilioChnl.EXPECT().ReceiveMessage(gomock.Any()).Return(&messages.Receive{Question: &query.Question{Sender: "42", Text: "off"}}, nil),
				mockSend:    twilioChnl.EXPECT().SendMessage(gomock.Any()).Return(nil),
			},
			want: `[{"text":"Turning off.","image":""},{"text":"❌","image":""}]`,
		},
		{
			name: "telegram endpoint test",
			bot:  testBot,
			args: args{
				endpoint:    fmt.Sprintf("%s/endpoints/telegram", ts.URL),
				message:     []byte(`{"sender": "42", "text": "on"}`),
				mockReceive: telegramChnl.EXPECT().ReceiveMessage(gomock.Any()).Return(&messages.Receive{Question: &query.Question{Sender: "42", Text: "on"}}, nil),
				mockSend:    telegramChnl.EXPECT().SendMessage(gomock.Any()).Return(nil),
			},
			want: `[{"text":"Turning on.","image":""}]`,
		},
		{
			name: "slack endpoint test",
			bot:  testBot,
			args: args{
				endpoint:    fmt.Sprintf("%s/endpoints/slack", ts.URL),
				message:     []byte(`{"sender": "42", "text": "on"}`),
				mockReceive: slackChnl.EXPECT().ReceiveMessage(gomock.Any()).Return(&messages.Receive{Question: &query.Question{Sender: "42", Text: "on"}}, nil),
				mockSend:    slackChnl.EXPECT().SendMessage(gomock.Any()).Return(nil),
			},
			want: `[{"text":"Can't do that.","image":""}]`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := http.Post(tt.args.endpoint, "application/json", bytes.NewBuffer(tt.args.message))
			if (err != nil) != tt.wantErr {
				t.Errorf("Bot.endpointHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if (err != nil) != tt.wantErr {
				t.Errorf("Bot.endpointHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(string(got), tt.want) {
				t.Errorf("Bot.endpointHandler() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestBot_Extensions(t *testing.T) {
	botPort, err := strconv.Atoi(testutils.GetFreePort(t))
	if err != nil {
		t.Fatal(err)
	}

	extensionPort := testutils.GetFreePort(t)

	testutils.RunGoExtension(t, testutils.Examples00TestPath, extensionPort)

	bc, err := bot.LoadConfig(testutils.Examples00TestPath, botPort)
	if err != nil {
		t.Fatal(err)
	}
	bc.Extensions.URL = fmt.Sprintf("http://127.0.0.1:%s", extensionPort)

	testBot, _, _, _, _, err := newTestBot(t)
	if err != nil {
		t.Fatalf("failed to load bot: %s", err)
	}

	_, err = testBot.Answer(&query.Question{
		Sender: "tester",
		Text:   "hello",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestBot_Answer(t *testing.T) {
	testBot, _, _, _, _, err := newTestBot(t)
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		question *query.Question
	}
	tests := []struct {
		name    string
		bot     *bot.Bot
		args    args
		want    []query.Answer
		wantErr bool
	}{
		{
			name: "turn on the thing",
			bot:  testBot,
			args: args{
				question: &query.Question{
					Sender: "42",
					Text:   "on",
				},
			},
			want: []query.Answer{{
				Text: "Turning on.",
			}},
		},
		{
			name: "turn off the thing",
			bot:  testBot,
			args: args{
				question: &query.Question{
					Sender: "42",
					Text:   "off",
				},
			},
			want: []query.Answer{
				{
					Text: "Turning off.",
				},
				{
					Text: "❌",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.bot.Answer(tt.args.question)
			if (err != nil) != tt.wantErr {
				t.Errorf("Bot.Answer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bot.Answer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBot_Predict(t *testing.T) {
	testBot, _, _, _, _, err := newTestBot(t)
	if err != nil {
		t.Fatal(err)
	}

	ts := httptest.NewServer(testBot.Router)
	defer ts.Close()

	predictEndpoint := fmt.Sprintf("%s/predict", ts.URL)

	type args struct {
		inputText []byte
	}
	tests := []struct {
		name    string
		bot     *bot.Bot
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test on",
			bot:  testBot,
			args: args{
				inputText: []byte(`{"text": "on"}`),
			},
			want: `{"original":"on","predicted":"turn_on","probability":0.999999999985}`,
		},
		{
			name: "test off",
			bot:  testBot,
			args: args{
				inputText: []byte(`{"text": "off"}`),
			},
			want: `{"original":"off","predicted":"turn_off","probability":0.999999999985}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := http.Post(predictEndpoint, "application/json", bytes.NewBuffer(tt.args.inputText))
			if (err != nil) != tt.wantErr {
				t.Errorf("Bot.predictHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if (err != nil) != tt.wantErr {
				t.Errorf("Bot.predictHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(string(got), tt.want) {
				t.Errorf("Bot.predictHandler() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestBot_Details(t *testing.T) {
	testBot, _, _, _, _, err := newTestBot(t)
	if err != nil {
		t.Fatal(err)
	}

	ts := httptest.NewServer(testBot.Router)
	defer ts.Close()

	detailsEndpoint := fmt.Sprintf("%s/senders", ts.URL)

	testBot.Store.Set("marcopolo", &fsm.FSM{})

	type args struct {
		sender string
	}
	tests := []struct {
		name    string
		bot     *bot.Bot
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test unknown",
			bot:  testBot,
			args: args{
				sender: "atlantis",
			},
			want: `sender does not exist
`,
		},
		{
			name: "test known",
			bot:  testBot,
			args: args{
				sender: "marcopolo",
			},
			want: `{"state":0,"slots":null}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := http.Get(fmt.Sprintf("%s/%s", detailsEndpoint, tt.args.sender))
			if (err != nil) != tt.wantErr {
				t.Errorf("Bot.endpointHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if (err != nil) != tt.wantErr {
				t.Errorf("Bot.endpointHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(string(got), tt.want) {
				t.Errorf("Bot.endpointHandler() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestBot_Run(t *testing.T) {
	botPort, err := strconv.Atoi(testutils.GetFreePort(t))
	if err != nil {
		t.Fatal(err)
	}

	bc, err := bot.LoadConfig(testutils.Examples05SimplePath, botPort)
	if err != nil {
		t.Fatal(err)
	}

	b, err := bot.New(bc)
	if err != nil {
		t.Fatalf("failed to load bot: %s", err)
	}

	go b.Run()
}

func newTestBot(t *testing.T) (*bot.Bot, *mockchannels.MockChannel, *mockchannels.MockChannel,
	*mockchannels.MockChannel, *mockchannels.MockChannel, error) {
	botConfig := &bot.Config{
		Name:       "chatto",
		Extensions: extension.Config{},
		Store:      fsm.StoreConfig{},
		Port:       0,
		Path:       testutils.Examples05SimplePath,
	}

	b := &bot.Bot{
		Name:   botConfig.Name,
		Store:  fsm.NewStore(botConfig.Store),
		Config: botConfig,
	}

	// Load Channels
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	b.Channels = &channels.Channels{}

	restChnl := mockchannels.NewMockChannel(ctrl)
	b.Channels.REST = restChnl

	twilioChnl := mockchannels.NewMockChannel(ctrl)
	b.Channels.Twilio = twilioChnl

	telegramChnl := mockchannels.NewMockChannel(ctrl)
	b.Channels.Telegram = telegramChnl

	slackChnl := mockchannels.NewMockChannel(ctrl)
	b.Channels.Slack = slackChnl

	// Load FSM
	fsmConfig, err := fsm.LoadConfig(botConfig.Path)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	b.Domain = fsm.New(fsmConfig)

	// Load Classifier
	classifConfig, err := clf.LoadConfig(botConfig.Path)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	b.Classifier = clf.New(classifConfig)

	// Load Extensions
	ext, err := extension.New(botConfig.Extensions)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	b.Extension = ext

	// Register HTTP handlers
	b.RegisterRoutes()

	log.Infof("My name is '%v'", b.Name)

	return b, restChnl, twilioChnl, telegramChnl, slackChnl, nil
}
