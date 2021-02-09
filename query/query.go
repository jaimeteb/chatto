package query

// Question for the FSM
type Question struct {
	Sender string
	Text   string
}

// Answer from the FSM
type Answer struct {
	Text  string `json:"text"`
	Image string `json:"image"`
}
