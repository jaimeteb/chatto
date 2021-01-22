package main

import (
	"log"

	"github.com/jaimeteb/chatto/ext"
	"github.com/jaimeteb/chatto/fsm"
)

func validateAnswer1(req *ext.Request) (res *ext.Response) {
	txt := req.Txt
	dom := req.Dom

	if !(txt == "1" || txt == "2" || txt == "3") {
		return &ext.Response{
			FSM: &fsm.FSM{
				State: dom.StateTable["question_1"],
				Slots: req.FSM.Slots,
			},
			Res: "Select one of the options",
		}
	}

	return &ext.Response{
		FSM: req.FSM,
		Res: "Question 2:\n" +
			"What is the capital of the state of Utah?\n" +
			"1. Salt Lake City\n" +
			"2. Jefferson City\n" +
			"3. Cheyenne",
	}
}

func validateAnswer2(req *ext.Request) (res *ext.Response) {
	txt := req.Txt
	dom := req.Dom

	if !(txt == "1" || txt == "2" || txt == "3") {
		return &ext.Response{
			FSM: &fsm.FSM{
				State: dom.StateTable["question_2"],
				Slots: req.FSM.Slots,
			},
			Res: "Select one of the options",
		}
	}

	return &ext.Response{
		FSM: req.FSM,
		Res: "Question 3:\n" +
			"Who painted Starry Night?\n" +
			"1. Pablo Picasso\n" +
			"2. Claude Monet\n" +
			"3. Vincent Van Gogh",
	}
}

func calculateScore(req *ext.Request) (res *ext.Response) {
	txt := req.Txt
	dom := req.Dom
	slt := req.FSM.Slots

	if !(txt == "1" || txt == "2" || txt == "3") {
		return &ext.Response{
			FSM: &fsm.FSM{
				State: dom.StateTable["question_3"],
				Slots: req.FSM.Slots,
			},
			Res: "Select one of the options",
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
		FSM: req.FSM,
		Res: message,
	}
}

var myExtMap = ext.ExtensionMap{
	"ext_val_ans_1": validateAnswer1,
	"ext_val_ans_2": validateAnswer2,
	"ext_score":     calculateScore,
}

func main() {
	if err := ext.ServeExtensionREST(myExtMap); err != nil {
		log.Fatalln(err)
	}
}
