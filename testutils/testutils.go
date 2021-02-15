package testutils

import (
	"bytes"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"syscall"
	"testing"

	log "github.com/sirupsen/logrus"
)

var (
	Examples00InvalidPath = "../_examples/404/"
	Examples00TestPath    = "../_examples/00_test/"
	Examples01MoodbotPath = "../_examples/01_moodbot/"
	Examples02MiscPath    = "../_examples/02_misc/"
	Examples03PokemonPath = "../_examples/03_pokemon/"
	Examples04TriviaPath  = "../_examples/04_trivia/"
	Examples05SimplePath  = "../_examples/05_simple/"
)

// GetFreePort returns an available port to use
func GetFreePort(t *testing.T) string {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	_ = listener.Close()

	return strings.Split(listener.Addr().String(), ":")[1]
}

// RunGoExtension for running extentions with unit tests
func RunGoExtension(t *testing.T, path, port string) {
	extension := fmt.Sprintf("%s/ext/ext.go", strings.TrimRight(path, "/"))

	cmd := exec.Command("go", "run", extension, "-port", port)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	t.Cleanup(func() {
		if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); err != nil {
			log.Error(err)
		}
	})
}
