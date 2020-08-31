package main

import (
	"github.com/jaimeteb/chatto/fsm"
)

func calculateScore(m *fsm.FSM) interface{} {
	answer1 := m.Slots["answer_1"].(string)
	answer2 := m.Slots["answer_2"].(string)
	answer3 := m.Slots["answer_3"].(string)

	score := 0
	if answer1 == "2" {
		score++
	}
	if answer2 == "1" {
		score++
	}
	if answer3 == "3" {
		score++
	}

	var message string
	switch score {
	case 0:
		message = "You got 0/3 answers right.\nBetter luck next time!"
	case 1:
		message = "You got 1/3 answers right.\nKeep trying!"
	case 2:
		message = "You got 2/3 answers right.\nPretty good!"
	case 3:
		message = "You got 3/3 answers right.\nYou are good! Congrats!"
	}
	return message
	// Type *start* to begin again
}

// Ext is exported
var Ext = fsm.FuncMap{
	"ext_score": calculateScore,
}

func main() {}
