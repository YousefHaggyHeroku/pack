package commands_test

import (
	"bytes"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	"github.com/spf13/cobra"

	"github.com/YousefHaggyHeroku/pack/internal/commands"
	"github.com/YousefHaggyHeroku/pack/internal/commands/testmocks"
	"github.com/YousefHaggyHeroku/pack/internal/config"
	ilogging "github.com/YousefHaggyHeroku/pack/internal/logging"
	"github.com/YousefHaggyHeroku/pack/logging"
	h "github.com/YousefHaggyHeroku/pack/testhelpers"
)

func TestBuilderCommand(t *testing.T) {
	spec.Run(t, "BuilderCommand", testBuilderCommand, spec.Parallel(), spec.Report(report.Terminal{}))
}

func testBuilderCommand(t *testing.T, when spec.G, it spec.S) {
	var (
		cmd    *cobra.Command
		logger logging.Logger
		outBuf bytes.Buffer
	)

	it.Before(func() {
		logger = ilogging.NewLogWithWriters(&outBuf, &outBuf)
		mockController := gomock.NewController(t)
		mockClient := testmocks.NewMockPackClient(mockController)
		cmd = commands.NewBuilderCommand(logger, config.Config{}, mockClient)
		cmd.SetOut(logging.GetWriterForLevel(logger, logging.InfoLevel))
	})

	when("builder", func() {
		it("prints help text", func() {
			cmd.SetArgs([]string{})
			h.AssertNil(t, cmd.Execute())
			output := outBuf.String()
			h.AssertContains(t, output, "Interact with builders")
			h.AssertContains(t, output, "Usage:")
			for _, command := range []string{"create", "suggest", "inspect"} {
				h.AssertContains(t, output, command)
				h.AssertNotContains(t, output, command+"-builder")
			}
		})
	})
}
