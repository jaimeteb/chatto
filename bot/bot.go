package bot

import (
	"github.com/jaimeteb/chatto/clf"
	"github.com/jaimeteb/chatto/fsm"
)

// Message models and incoming/outgoing message
type Message struct {
	Sender string `json:"sender"`
	Text   string `json:"text"`
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

// Bot models a bot with a Classifier and an FSM
type Bot struct {
	// Machines   map[string]*fsm.FSM
	Machines   fsm.StoreFSM
	Domain     fsm.Domain
	Classifier clf.Classifier
	Extension  fsm.Extension
	Endpoints  map[string]interface{}
}

// Prediction models a classifier prediction and its orignal string
type Prediction struct {
	Original    string  `json:"original"`
	Predicted   string  `json:"predicted"`
	Probability float64 `json:"probability"`
}

// Answer takes a user input and executes a transition on the FSM if possible
func (b Bot) Answer(mess Message) string {
	// if _, ok := b.Machines[mess.Sender]; !ok {
	// 	b.Machines[mess.Sender] = &fsm.FSM{State: 0, Slots: make(map[string]interface{})}
	// }

	// cmd, _ := b.Classifier.Predict(mess.Text) // Predict command from text using classifier
	// return b.Machines[mess.Sender].ExecuteCmd(cmd, mess.Text, b.Domain, b.Extension)

	if !b.Machines.Exists(mess.Sender) {
		// b.Machines[mess.Sender] = &fsm.FSM{
		// 	State: i,
		// 	Slots: make(map[string]interface{}),
		// }
		b.Machines.SetState(mess.Sender, 0) // What if I don't set a state
	}

	cmd, _ := b.Classifier.Predict(mess.Text)
	return cmd
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
