package main

import (
	"log"

	ext "github.com/jaimeteb/chatto/extension"
	"github.com/jaimeteb/chatto/fsm"
	"github.com/jaimeteb/chatto/query"
)

func validateAnswer1(req *ext.Request) (res *ext.Response) {
	ans := req.FSM.Slots["answer_1"]
	dom := req.DB

	if !(ans == "1" || ans == "2" || ans == "3") {
		return &ext.Response{
			FSM: &fsm.FSM{
				State: dom.StateTable["question_1"],
				Slots: req.FSM.Slots,
			},
			Answers: query.Answers("Select one of the options"),
		}
	}

	message := "Question 2:\n" +
		"What is the capital of the state of Utah?\n" +
		"1. Salt Lake City\n" +
		"2. Jefferson City\n" +
		"3. Cheyenne"

	return &ext.Response{
		FSM:     req.FSM,
		Answers: query.Answers(message),
	}
}

func validateAnswer2(req *ext.Request) (res *ext.Response) {
	ans := req.FSM.Slots["answer_2"]
	dom := req.DB

	if !(ans == "1" || ans == "2" || ans == "3") {
		return &ext.Response{
			FSM: &fsm.FSM{
				State: dom.StateTable["question_2"],
				Slots: req.FSM.Slots,
			},
			Answers: query.Answers("Select one of the options"),
		}
	}

	message := "Question 3:\n" +
		"Who painted Starry Night?\n" +
		"1. Pablo Picasso\n" +
		"2. Claude Monet\n" +
		"3. Vincent Van Gogh"

	return &ext.Response{
		FSM:     req.FSM,
		Answers: query.Answers(message),
	}
}

func calculateScore(req *ext.Request) (res *ext.Response) {
	ans := req.FSM.Slots["answer_1"]
	dom := req.DB
	slt := req.FSM.Slots

	if !(ans == "1" || ans == "2" || ans == "3") {
		return &ext.Response{
			FSM: &fsm.FSM{
				State: dom.StateTable["question_3"],
				Slots: req.FSM.Slots,
			},
			Answers: query.Answers("Select one of the options"),
		}
	}

	answer1 := slt["answer_1"]
	answer2 := slt["answer_2"]
	answer3 := slt["answer_3"]

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

	return &ext.Response{
		FSM:     req.FSM,
		Answers: query.Answers(message),
	}
}

var myExtMap = ext.RegisteredFuncs{
	"ext_val_ans_1": validateAnswer1,
	"ext_val_ans_2": validateAnswer2,
	"ext_score":     calculateScore,
}

func main() {
	if err := ext.ServeExtensionREST(myExtMap); err != nil {
		log.Fatalln(err)
	}
}
