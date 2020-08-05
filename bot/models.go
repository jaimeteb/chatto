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

// Bot models a bot with a Classifier and an FSM
type Bot struct {
	Machines   map[string]*fsm.FSM
	Domain     fsm.Domain
	Classifier clf.Classifier
}

// Answer takes a user input and executes a transition on the FSM if possible
func (b Bot) Answer(mess Message) string {
	if _, ok := b.Machines[mess.Sender]; !ok {
		b.Machines[mess.Sender] = &fsm.FSM{State: 0, Slots: make(map[string]interface{})}
	}

	cmd, _ := b.Classifier.Predict(mess.Text) // Predict command from text using classifier
	return b.Machines[mess.Sender].ExecuteCmd(cmd, b.Domain)
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
