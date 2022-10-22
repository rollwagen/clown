package main

import (
	"os"

	"github.com/rollwagen/clown/cmd"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)

	cmd.Execute()
}
