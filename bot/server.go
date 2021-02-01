package bot

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	cmn "github.com/jaimeteb/chatto/common"
	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func (b Bot) restEndpointHandler(w http.ResponseWriter, r *http.Request) {
	mess, err := b.Clients.REST.RecieveMessage(w, r)
	if err != nil {
		log.Error(err)
		return
	}

	resp := b.Answer(mess)

	if err := SendMessages(resp, &b.Clients.REST, mess.Sender, w); err != nil {
		log.Error(err)
		return
	}
}

func (b Bot) telegramEndpointHandler(w http.ResponseWriter, r *http.Request) {
	mess, err := b.Clients.Telegram.RecieveMessage(w, r)
	if err != nil {
		log.Error(err)
		return
	}

	resp := b.Answer(mess)

	if err := SendMessages(resp, &b.Clients.Telegram, mess.Sender, w); err != nil {
		log.Error(err)
		return
	}
}

func (b Bot) twilioEndpointHandler(w http.ResponseWriter, r *http.Request) {
	mess, err := b.Clients.Twilio.RecieveMessage(w, r)
	if err != nil {
		log.Error(err)
		return
	}

	resp := b.Answer(mess)

	if err := SendMessages(resp, &b.Clients.Twilio, mess.Sender, w); err != nil {
		log.Error(err)
		return
	}
}

func (b Bot) slackEndpointHandler(w http.ResponseWriter, r *http.Request) {
	mess, err := b.Clients.Slack.RecieveMessage(w, r)
	if err != nil {
		log.Error(err)
		return
	} else if (mess == cmn.Message{}) {
		return
	}

	resp := b.Answer(mess)

	if err := SendMessages(resp, &b.Clients.Slack, mess.Sender, w); err != nil {
		log.Error(err)
		return
	}
}

func (b Bot) detailsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	senderObj := b.Machines.Get(vars["sender"])

	js, err := json.Marshal(senderObj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (b Bot) predictHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var mess cmn.Message

	err := decoder.Decode(&mess)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	inputText := mess.Text
	prediction, prob := b.Classifier.Predict(inputText)
	ans := Prediction{inputText, prediction, prob}

	js, err := json.Marshal(ans)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (b Bot) slackSocketmodeHandler() {
	client := b.Clients.Slack.Socketmode

	go func() {
		for evt := range client.Events {
			switch evt.Type {
			case socketmode.EventTypeConnecting:
				log.Info("Connecting to Slack with Socket Mode...")
			case socketmode.EventTypeConnectionError:
				log.Error("Connection to Slack failed. Retrying later...")
			case socketmode.EventTypeConnected:
				log.Info("Connected to Slack with Socket Mode")
			case socketmode.EventTypeEventsAPI:
				eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
				if !ok {
					log.Warnf("Ignored %+v\n", evt)

					continue
				}

				log.Infof("Event received: %+v\n", eventsAPIEvent)

				client.Ack(*evt.Request)

				switch eventsAPIEvent.Type {
				case slackevents.CallbackEvent:
					innerEvent := eventsAPIEvent.InnerEvent
					switch ev := innerEvent.Data.(type) {
					case *slackevents.MessageEvent:
						resp := b.Answer(cmn.Message{Sender: ev.User, Text: ev.Text})

						messages, _, err := Messages(resp)
						if err != nil {
							log.Error(err)

							continue
						}

						for n := range messages {
							sendErr := b.Clients.Slack.SendMessage(messages[n], ev.User)
							if sendErr != nil {
								log.Errorf("Failed posting message: %v", err)

								continue
							}
						}
					case *slackevents.AppMentionEvent:
						resp := b.Answer(cmn.Message{Sender: ev.User, Text: ev.Text})

						messages, _, err := Messages(resp)
						if err != nil {
							log.Error(err)

							continue
						}

						for n := range messages {
							sendErr := b.Clients.Slack.SendMessage(messages[n], ev.User)
							if sendErr != nil {
								log.Errorf("Failed posting message: %v", err)

								continue
							}
						}
					}
				default:
					log.Debugf("Unsupported Events API event received")
				}
			case socketmode.EventTypeInteractive:
				// TODO: Support interactions.
			case socketmode.EventTypeSlashCommand:
				// TODO: Support slash commands.
			default:
				log.Debugf("Unexpected event type received: %s", evt.Type)
			}
		}
	}()

	client.Run()
}

// ServeBot function
func ServeBot(path *string, port *int) {
	bot := LoadBot(path)

	// log.Info("\n" + LOGO)
	log.Info("Server started")

	r := mux.NewRouter()

	// Integration Endpoints
	r.HandleFunc("/endpoints/rest", bot.restEndpointHandler).Methods("POST")
	r.HandleFunc("/endpoints/telegram", bot.telegramEndpointHandler).Methods("POST")
	r.HandleFunc("/endpoints/twilio", bot.twilioEndpointHandler).Methods("POST")
	r.HandleFunc("/endpoints/slack", bot.slackEndpointHandler).Methods("POST")

	if bot.Clients.Slack.Socketmode != nil {
		go bot.slackSocketmodeHandler()
	}

	// Prediction and Sender Endpoints
	r.HandleFunc("/predict", bot.predictHandler).Methods("POST")
	r.HandleFunc("/senders/{sender}", bot.detailsHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *port), r))
}
