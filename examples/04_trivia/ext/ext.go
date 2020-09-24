package main

import (
	"github.com/jaimeteb/chatto/fsm"
)

func validateAnswer1(m *fsm.FSM, dom *fsm.Domain, txt string) interface{} {
	if !(txt == "1" || txt == "2" || txt == "3") {
		(*m).State = dom.StateTable["question_1"]
		return "Select one of the options"
	}

	return "Question 2:\n" +
		"What is the capital of the state of Utah?\n" +
		"1. Salt Lake City\n" +
		"2. Jefferson City\n" +
		"3. Cheyenne"
}

func validateAnswer2(m *fsm.FSM, dom *fsm.Domain, txt string) interface{} {
	if !(txt == "1" || txt == "2" || txt == "3") {
		(*m).State = dom.StateTable["question_2"]
		return "Select one of the options"
	}

	return "Question 3:\n" +
		"Who painted Starry Night?\n" +
		"1. Pablo Picasso\n" +
		"2. Claude Monet\n" +
		"3. Vincent Van Gogh"
}

func calculateScore(m *fsm.FSM, dom *fsm.Domain, txt string) interface{} {
	if !(txt == "1" || txt == "2" || txt == "3") {
		(*m).State = dom.StateTable["question_3"]
		return "Select one of the options"
	}

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
// var Ext = fsm.FuncMap{
// 	"ext_val_ans_1": validateAnswer1,
// 	"ext_val_ans_2": validateAnswer2,
// 	"ext_score":     calculateScore,
// }

func main() {}
