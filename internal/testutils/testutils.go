package testutils

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/jaimeteb/chatto/internal/clf/wordvectors"
	log "github.com/sirupsen/logrus"
)

var (
	Examples00InvalidPath = "../examples/404/"
	Examples00TestPath    = "../examples/00_test/"
	Examples01MoodbotPath = "../examples/01_moodbot/"
	Examples02MiscPath    = "../examples/02_misc/"
	Examples03PokemonPath = "../examples/03_pokemon/"
	Examples04TriviaPath  = "../examples/04_trivia/"
	Examples05SimplePath  = "../examples/05_simple/"
	TestWordVectors       = "../internal/testutils/testvec"
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

// RunGoExtension for running extensions with unit tests
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

func GetTestWordVectors(skipOOV bool) (*wordvectors.VectorMap, error) {
	return wordvectors.NewVectorMap(&wordvectors.Config{
		WordVectorsFile: TestWordVectors,
		Truncate:        1,
		SkipOOV:         skipOOV,
	})
}

func RemoveGobFiles() {
	files, err := filepath.Glob("./**/**.gob")
	if err != nil {
		log.Error(err)
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			log.Error(err)
		}
	}
}
