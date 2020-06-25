package pkg

// Bot struct
type Bot struct {
	ID        int
	Name      string
	StateSize int
	MaxHist   int
	History   History
	Chain     Chain
}

// Answer processes input from user and produce output
func (b Bot) Answer(mess Message) Message {
	id := mess.Sender
	b.History.Append(id, mess)

	// lastSender :=

	respCode := b.Chain.Predict(b.History.Messages[id])

	resp := Message{
		Sender: b.Name,
		Text:   respCode,
	}

	// switch b.History.Messages[mess.Sender][len(b.History.Messages[mess.Sender])-1].Text {
	// case "/hello":
	// 	resp.Sender = b.Name
	// 	resp.Text = "Hello there!"
	// case "/bye":
	// 	resp.Sender = b.Name
	// 	resp.Text = "Goodbye!"
	// default:
	// 	resp.Sender = b.Name
	// 	resp.Text = "..."
	// }

	b.History.Append(id, resp)
	b.History.Print(id)
	return resp
}
