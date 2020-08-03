package bot

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jaimeteb/chatto/clf"
	"github.com/jaimeteb/chatto/fsm"
)

// CLI runs a bot in a command line interface
func CLI(path *string) {
	domain := fsm.Create(path)
	classifier := clf.Create(path)

	machines := make(map[string]*fsm.FSM)
	bot := Bot{machines, domain, classifier}

	log.Println("CLI started")

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("you:\t| ")
		cmd, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}

		fmt.Print("botto:\t| ")
		resp := bot.Answer(Message{"cli", strings.TrimSuffix(cmd, "\n")})
		fmt.Println(resp)
	}
}
