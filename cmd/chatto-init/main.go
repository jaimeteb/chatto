package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

var clfFile string = `classification:
  - command: "greet"
    texts:
      - "hey"
      - "hello"
      - "hi"
      - "hello there"
      - "good morning"
      - "good evening"
      - "moin"
      - "hey there"
      - "let's go"
      - "hey dude"
      - "goodmorning"
      - "goodevening"
      - "good afternoon"

  - command: "good"
    texts:
      - "perfect"
      - "great"
      - "amazing"
      - "feeling like a king"
      - "wonderful"
      - "I am feeling very good"
      - "I am great"
      - "I am amazing"
      - "I am going to save the world"
      - "super stoked"
      - "extremely good"
      - "so so perfect"
      - "so good"
      - "so perfect"

  - command: "bad"
    texts:
      - "my day was horrible"
      - "I am sad"
      - "I don't feel very well"
      - "I am disappointed"
      - "super sad"
      - "I'm so sad"
      - "sad"
      - "very sad"
      - "unhappy"
      - "not good"
      - "not very good"
      - "extremely sad"
      - "so saad"
      - "so sad "

  - command: "yes"
    texts:
      - "yes"
      - "indeed"
      - "of course"
      - "that sounds good"
      - "correct "

  - command: "no"
    texts:
      - "no"
      - "never"
      - "I don't think so"
      - "don't like that"
      - "no way"
`

var fsmFile string = `states:
  - "initial"
  - "ask_mood"
  - "say_good"
  - "say_bad"
  - "say_bye"

commands:
  - "greet"
  - "good"
  - "bad"
  - "yes"
  - "no"

transitions:
  - from:
      - initial
    into: ask_mood
    command: greet
    answers: 
      - text: "Hello! How are you?"

  - from:
      - ask_mood
    into: initial
    command: good
    answers: 
      - text: "Great! :)"

  - from:
      - ask_mood
    into: say_bad
    command: bad
    answers:
      - text: "Oh don't be sad :("
        image: https://i.imgur.com/8MU0IUT.jpeg
      - text: "Did that help?"
    # extension: dont_feel_bad

  - from:
      - say_bad
    into: initial
    command: "yes"
    answers:
      - text: "I'm glad! :)"

  - from:
      - say_bad
    into: initial
    command: "no"
    answers: 
      - text: "Oh I'm sorry"

defaults:
  unknown: "Unknown command, try again please."
  unsure: "Not sure I understood, try again please."
  error: "An error occurred."
`

var extFile string = `package main

import (
	"log"

	"github.com/jaimeteb/chatto/extensions"
	"github.com/jaimeteb/chatto/query"
)

func dontFeelBad(req *extension.ExecuteExtensionRequest) (res *extension.ExecuteExtensionResponse) {
	return &extension.ExecuteExtensionResponse{
		FSM: req.FSM,
		Answers: []query.Answer{
			{
				Text: "Oh don't be sad :(",
			},
			{
				Text:  "Did that help?",
				Image: "https://i.imgur.com/8MU0IUT.jpeg",
			},
		},
	}
}

var registeredExtensions = extension.RegisteredExtensions{
	"dont_feel_bad": dontFeelBad,
}

func main() {
	if err := extension.ServeREST(registeredExtensions); err != nil {
		log.Fatalln(err)
	}
}
`

var botFile string = `extensions:
  type: REST
  url: http://localhost:8770
`

var chnFile string = `telegram:
  bot_key:

twilio:
  account_sid:
  auth_token:
  number:

slack:
  token:
  app_token:
`

func main() {
	filePath := flag.String("path", ".", "Where to write initial files.")
	flag.Parse()

	if *filePath != "." {
		if _, err := os.Stat(*filePath); os.IsNotExist(err) {
			if err := os.MkdirAll(*filePath, 0755); err != nil {
				fmt.Printf("Couldn't create directory: %v", err)
				return
			}
		}
		if _, err := os.Stat(path.Join(*filePath, "ext")); os.IsNotExist(err) {
			if err := os.MkdirAll(path.Join(*filePath, "ext"), 0755); err != nil {
				fmt.Printf("Couldn't create directory: %v", err)
				return
			}
		}
	}

	fileMap := map[string]string{
		"clf.yml":     clfFile,
		"fsm.yml":     fsmFile,
		"bot.yml":     botFile,
		"chn.yml":     chnFile,
		"ext/main.go": extFile,
	}

	for fileName, fileContent := range fileMap {
		if err := ioutil.WriteFile(path.Join(*filePath, fileName), []byte(fileContent), 0600); err != nil {
			fmt.Printf("Couldn't write %s file: %v\n", fileName, err)
			return
		}
	}
	fmt.Println("Initial project files written successfully.")
}
