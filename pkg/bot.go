package pkg

// Bot struct
type Bot struct {
	ID        int
	Name      string
	StateSize int
	MaxHist   int
	History   History
}

// Answer processes input from user and produce output
func (b Bot) Answer(mess Message) Message {
	id := mess.Sender
	b.History.Append(id, mess)

	var resp Message
	switch b.History.Messages[mess.Sender][len(b.History.Messages[mess.Sender])-1].Text {
	case "/hello":
		resp.Sender = b.Name
		resp.Text = "Hello there!"
	case "/bye":
		resp.Sender = b.Name
		resp.Text = "Goodbye!"
	default:
		resp.Sender = b.Name
		resp.Text = "..."
	}

	b.History.Append(id, resp)
	b.History.Print(id)
	return resp
}
