package main

import (
	"log"

	"github.com/jaimeteb/chatto/extensions"
	"github.com/jaimeteb/chatto/fsm"
	"github.com/jaimeteb/chatto/query"
)

func validateAnswer1(req *extensions.ExecuteExtensionRequest) (res *extensions.ExecuteExtensionResponse) {
	ans := req.FSM.Slots["answer_1"]
	fsmDomain := req.Domain

	if !(ans == "1" || ans == "2" || ans == "3") {
		return &extensions.ExecuteExtensionResponse{
			FSM: &fsm.FSM{
				State: fsmDomain.StateTable["question_1"],
				Slots: req.FSM.Slots,
			},
			Answers: []query.Answer{{Text: "Select one of the options"}},
		}
	}

	return &extensions.ExecuteExtensionResponse{
		FSM: req.FSM,
		Answers: []query.Answer{{Text: "Question 2:\n" +
			"What is the capital of the state of Utah?\n" +
			"1. Salt Lake City\n" +
			"2. Jefferson City\n" +
			"3. Cheyenne"}},
	}
}

func validateAnswer2(req *extensions.ExecuteExtensionRequest) (res *extensions.ExecuteExtensionResponse) {
	ans := req.FSM.Slots["answer_2"]
	fsmDomain := req.Domain

	if !(ans == "1" || ans == "2" || ans == "3") {
		return &extensions.ExecuteExtensionResponse{
			FSM: &fsm.FSM{
				State: fsmDomain.StateTable["question_2"],
				Slots: req.FSM.Slots,
			},
			Answers: []query.Answer{{Text: "Select one of the options"}},
		}
	}

	return &extensions.ExecuteExtensionResponse{
		FSM: req.FSM,
		Answers: []query.Answer{{Text: "Question 3:\n" +
			"Who painted Starry Night?\n" +
			"1. Pablo Picasso\n" +
			"2. Claude Monet\n" +
			"3. Vincent Van Gogh"}},
	}
}

func calculateScore(req *extensions.ExecuteExtensionRequest) (res *extensions.ExecuteExtensionResponse) {
	ans := req.FSM.Slots["answer_1"]
	fsmDomain := req.Domain
	slt := req.FSM.Slots

	if !(ans == "1" || ans == "2" || ans == "3") {
		return &extensions.ExecuteExtensionResponse{
			FSM: &fsm.FSM{
				State: fsmDomain.StateTable["question_3"],
				Slots: req.FSM.Slots,
			},
			Answers: []query.Answer{{Text: "Select one of the options"}},
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

	return &extensions.ExecuteExtensionResponse{
		FSM:     req.FSM,
		Answers: []query.Answer{{Text: message}},
	}
}

var registeredExtensions = extensions.RegisteredExtensions{
	"val_ans_1": validateAnswer1,
	"val_ans_2": validateAnswer2,
	"score":     calculateScore,
}

func main() {
	if err := extensions.ServeREST(registeredExtensions); err != nil {
		log.Fatalln(err)
	}
}
