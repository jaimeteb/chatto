package pkg

// Message struct
type Message struct {
	Sender string `json:"sender"`
	Text   string `json:"text"`
}

// // History type maps user IDs to their message histories
// type History map[string][]Message

// // UserFSM type maps user IDs to their FSMs
// type UserFSM map[string]FSM

// // Answer processes input from user and produce output
// func (u UserFSM) Answer(mess Message) Message {
// 	// id := mess.Sender

// 	log.Println(u)

// 	if userMachine, ok := u[mess.Sender]; !ok {
// 		userMachine := &FSM{State: Initial}
// 		userMachine.ExecuteCmd(mess.Text)
// 	} else {
// 		userMachine.ExecuteCmd(mess.Text)
// 	}

// 	resp := Message{
// 		Sender: "foo",
// 		Text:   "bar",
// 	}

// 	// b.History.Print(id)
// 	return resp
// }
