package bot

import (
	"bufio"
	"fmt"
	"os"

	"github.com/jaimeteb/chatto/clf"
	"github.com/jaimeteb/chatto/fsm"
)

// CLI runs a bot in a command line interface
func CLI() {
	domain := fsm.Create()
	classifier := clf.Create()
	machines := make(map[string]*fsm.FSM)
	machines["cli"] = &fsm.FSM{State: 0}
	bot := Bot{machines, domain, classifier}

	fmt.Println("CLI started")

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("you:\t| ")
		cmd, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}

		fmt.Print("botto:\t| ")
		resp := bot.Answer(Message{"cli", cmd})
		fmt.Println(resp)
	}
}
