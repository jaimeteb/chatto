package bot

import (
	"log"
	"os"

	"github.com/jaimeteb/chatto/clf"
	"github.com/jaimeteb/chatto/fsm"
)

// Message models and incoming/outgoing message
type Message struct {
	Sender string      `json:"sender"`
	Text   interface{} `json:"text"`
}

// TelegramMessageIn models a telegram incoming message
type TelegramMessageIn struct {
	UpdateID int                    `json:"update_id"`
	Message  TelegramMessageInInner `json:"message"`
}

// TelegramMessageInInner models a telegram incoming message inner struct
type TelegramMessageInInner struct {
	MessageID int                        `json:"message_id"`
	From      TelegramMessageInInnerFrom `json:"from"`
	Date      int                        `json:"date"`
	Text      string                     `json:"text"`
}

// TelegramMessageInInnerFrom models a telegram incoming message inner struct
type TelegramMessageInInnerFrom struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
}

// TwilioMessageIn models an incoming Twilio message
type TwilioMessageIn struct {
	From             string `form:"From"`
	Body             string `form:"Body"`
	To               string `form:"To"`
	MediaURL         string `form:"MediaUrl"`
	MediaContentType string `form:"MediaContentType"`
	MessageSid       string `form:"MessageSid"`
	SmsStatus        string `form:"SmsStatus"`
	AccountSid       string `form:"AccountSid"`
	Sid              string `form:"Sid"`
	SmsSid           string `form:"SmsSid"`
	SmsMessageSid    string `form:"SmsMessageSid"`
	NumMedia         int    `form:"NumMedia"`
	NumSegments      int    `form:"NumSegments"`
	APIVersion       string `form:"ApiVersion"`
}

// Bot models a bot with a Classifier and an FSM
type Bot struct {
	// Machines   map[string]*fsm.FSM
	Machines   fsm.StoreFSM
	Domain     fsm.Domain
	Classifier clf.Classifier
	Extension  fsm.Extension
	Clients    map[string]interface{}
}

// Prediction models a classifier prediction and its orignal string
type Prediction struct {
	Original    string  `json:"original"`
	Predicted   string  `json:"predicted"`
	Probability float64 `json:"probability"`
}

// Answer takes a user input and executes a transition on the FSM if possible
func (b Bot) Answer(mess Message) interface{} {
	if !b.Machines.Exists(mess.Sender) {
		b.Machines.Set(
			mess.Sender,
			&fsm.FSM{
				State: 0,
				Slots: make(map[string]string),
			},
		)
	}

	inputMessage := mess.Text.(string)
	cmd, _ := b.Classifier.Predict(inputMessage)

	m := b.Machines.Get(mess.Sender)
	resp := m.ExecuteCmd(cmd, inputMessage, b.Domain, b.Extension)
	b.Machines.Set(mess.Sender, m)

	return resp
}

// LoadBot loads all configurations and returns a Bot
func LoadBot(path *string) Bot {
	domain := fsm.Create(path)
	classifier := clf.Create(path)

	// Load Extensions
	extension, err := fsm.LoadExtension(path)
	if err != nil {
		log.Println("Using bot without extensions.")
	}

	// Load clients
	clients := LoadClients(path)

	var machines fsm.StoreFSM
	// REDIS
	if redisHost := os.Getenv("REDIS_HOST"); redisHost != "" {
		machines = &fsm.RedisStoreFSM{R: fsm.RDB}
		log.Println("Registered RedisStoreFSM")
	} else {
		machines = &fsm.CacheStoreFSM{}
		log.Println("Registered CacheStoreFSM")
	}

	return Bot{machines, domain, classifier, extension, clients}
}

// LOGO for Chatto
const LOGO = `
                           *******                          
                  *************************                 
             *********                *********             
          *******                           ******.         
        *****                                  ******       
      *****                                       *****     
    *****                                           *****   
   ****                                              .****  
  ****         ********,             *********         **** 
 ****       .******.******         ******.******       .****
 ****       ****       ****       ****       ****       ****
****                                                    ****
****                                                     ***
****                                                    ****
 ****                                                   ****
 ****                  ****       ****                 ****.
  ****                 *****     ****                  **** 
   ****.                 ***********                 *****  
    *****                                           ****    
      *****                                       *****     
        ******                                 ******       
          .******                          .******          
              *********               *********             
                  .***********************.                 
`
