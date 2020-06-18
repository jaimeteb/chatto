package models

// Bot struct
type Bot struct {
	ID      int
	Name    string
	History History
}

// func reslice(s []int, index int) []int {
// 	return append(s[:index], s[index+1:]...)
// }

// Answer :
// Process input from user and produce output
func (b Bot) Answer(mess Message) Message {
	b.History.Messages[mess.Sender] = append(b.History.Messages[mess.Sender], mess)

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

	b.History.Messages[mess.Sender] = append(b.History.Messages[mess.Sender], resp)
	b.History.Print(mess.Sender)
	return resp
}
