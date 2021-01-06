package bot

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/jaimeteb/chatto/clf"
	"github.com/jaimeteb/chatto/fsm"
)

var fsmYaml = `
states:
  - "off"
  - "on"
commands:
  - "turn_on"
  - "turn_off"
functions:
  - tuple:
      command: "turn_on"
      state: "off"
    transition: "on"
    message: "Turning on."
  - tuple:
      command: "turn_off"
      state: "on"
    transition: "off"
    message: "Turning off."
defaults:
  unknown: "Can't do that."
  unsure: "???"
`

var clfYaml = `
classification:
  - command: "turn_on"
    texts:
      - "turn on"
      - "on"

  - command: "turn_off"
    texts:
      - "turn off"
      - "off"
`

func writeDummyFileClf() error {
	clfFile := []byte(clfYaml)
	return ioutil.WriteFile("clf.yml", clfFile, 0644)
}

func removeDummyFileClf() error {
	return os.Remove("clf.yml")
}

func writeDummyFileFSM() error {
	fsmFile := []byte(fsmYaml)
	return ioutil.WriteFile("fsm.yml", fsmFile, 0644)
}

func removeDummyFileFSM() error {
	return os.Remove("fsm.yml")
}

func TestBot(t *testing.T) {
	if err := writeDummyFileClf(); err != nil {
		t.Errorf(err.Error())
	}
	if err := writeDummyFileFSM(); err != nil {
		t.Errorf(err.Error())
	}

	here := "."
	bc := LoadBotConfig(&here)
	domain := fsm.Create(&here)
	classifier := clf.Create(&here)
	extension := fsm.LoadExtensions(bc.Extensions)
	machines := &fsm.CacheStoreFSM{}
	endpoints := make(map[string]interface{})
	bot := Bot{"botto", machines, domain, classifier, extension, endpoints}

	resp1 := bot.Answer(Message{"foo", "on"})
	if resp1 != "Turning on." {
		t.Errorf("resp1 is incorrect, got: %v, want: %v.", resp1, "Turning on.")
	}

	resp2 := bot.Answer(Message{"foo", "on"})
	if resp2 != "Can't do that." {
		t.Errorf("resp2 is incorrect, got: %v, want: %v.", resp2, "Can't do that.")
	}

	if err := removeDummyFileClf(); err != nil {
		t.Errorf(err.Error())
	}
	if err := removeDummyFileFSM(); err != nil {
		t.Errorf(err.Error())
	}
}
