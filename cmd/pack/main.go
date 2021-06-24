package main

import (
	"os"

	"github.com/heroku/color"

	"github.com/YousefHaggyHeroku/pack/cmd"

	"github.com/YousefHaggyHeroku/pack/internal/commands"
	clilogger "github.com/YousefHaggyHeroku/pack/internal/logging"
	"github.com/buildpacks/pack"
)

func main() {
	// create logger with defaults
	logger := clilogger.NewLogWithWriters(color.Stdout(), color.Stderr())

	rootCmd, err := cmd.NewPackCommand(logger)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	ctx := commands.CreateCancellableContext()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		if _, isSoftError := err.(pack.SoftError); isSoftError {
			os.Exit(2)
		}
		os.Exit(1)
	}
}
