package commands_test

import (
	"bytes"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	"github.com/spf13/cobra"

	"github.com/YousefHaggyHeroku/pack/internal/commands"
	"github.com/YousefHaggyHeroku/pack/internal/commands/fakes"
	"github.com/YousefHaggyHeroku/pack/internal/commands/testmocks"
	"github.com/YousefHaggyHeroku/pack/internal/config"
	ilogging "github.com/YousefHaggyHeroku/pack/internal/logging"
	"github.com/buildpacks/pack/logging"
	h "github.com/buildpacks/pack/testhelpers"
)

func TestBuildpackCommand(t *testing.T) {
	spec.Run(t, "BuildpackCommand", testBuildpackCommand, spec.Parallel(), spec.Report(report.Terminal{}))
}

func testBuildpackCommand(t *testing.T, when spec.G, it spec.S) {
	var (
		cmd        *cobra.Command
		logger     logging.Logger
		outBuf     bytes.Buffer
		mockClient *testmocks.MockPackClient
	)

	it.Before(func() {
		logger = ilogging.NewLogWithWriters(&outBuf, &outBuf)
		mockController := gomock.NewController(t)
		mockClient = testmocks.NewMockPackClient(mockController)
		cmd = commands.NewBuildpackCommand(logger, config.Config{}, mockClient, fakes.NewFakePackageConfigReader())
		cmd.SetOut(logging.GetWriterForLevel(logger, logging.InfoLevel))
	})

	when("buildpack", func() {
		it("prints help text", func() {
			cmd.SetArgs([]string{})
			h.AssertNil(t, cmd.Execute())
			output := outBuf.String()
			h.AssertContains(t, output, "Interact with buildpacks")
			for _, command := range []string{"Usage", "package", "register", "yank", "pull", "inspect"} {
				h.AssertContains(t, output, command)
			}
		})
	})
}
