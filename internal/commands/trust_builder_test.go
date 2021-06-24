package commands_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/heroku/color"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	"github.com/spf13/cobra"

	"github.com/YousefHaggyHeroku/pack/internal/commands"
	"github.com/YousefHaggyHeroku/pack/internal/config"
	ilogging "github.com/YousefHaggyHeroku/pack/internal/logging"
	"github.com/buildpacks/pack/logging"
	h "github.com/buildpacks/pack/testhelpers"
)

func TestTrustBuilderCommand(t *testing.T) {
	color.Disable(true)
	defer color.Disable(false)
	spec.Run(t, "Commands", testTrustBuilderCommand, spec.Random(), spec.Report(report.Terminal{}))
}

func testTrustBuilderCommand(t *testing.T, when spec.G, it spec.S) {
	var (
		command      *cobra.Command
		logger       logging.Logger
		outBuf       bytes.Buffer
		tempPackHome string
		configPath   string
	)

	it.Before(func() {
		var err error

		logger = ilogging.NewLogWithWriters(&outBuf, &outBuf)
		tempPackHome, err = ioutil.TempDir("", "pack-home")
		h.AssertNil(t, err)
		configPath = filepath.Join(tempPackHome, "config.toml")
		command = commands.TrustBuilder(logger, config.Config{}, configPath)
	})

	it.After(func() {
		h.AssertNil(t, os.RemoveAll(tempPackHome))
	})

	when("#TrustBuilder", func() {
		when("no builder is provided", func() {
			it("prints usage", func() {
				command.SetArgs([]string{})
				h.AssertError(t, command.Execute(), "accepts 1 arg(s)")
			})
		})

		when("builder is provided", func() {
			when("builder is not already trusted", func() {
				it("updates the config", func() {
					command.SetArgs([]string{"some-builder"})
					h.AssertNil(t, command.Execute())

					b, err := ioutil.ReadFile(configPath)
					h.AssertNil(t, err)
					h.AssertContains(t, string(b), `[[trusted-builders]]
  name = "some-builder"`)
				})
			})

			when("builder is already trusted", func() {
				it("does nothing", func() {
					command.SetArgs([]string{"some-already-trusted-builder"})
					h.AssertNil(t, command.Execute())
					oldContents, err := ioutil.ReadFile(configPath)
					h.AssertNil(t, err)

					command.SetArgs([]string{"some-already-trusted-builder"})
					h.AssertNil(t, command.Execute())

					newContents, err := ioutil.ReadFile(configPath)
					h.AssertNil(t, err)
					h.AssertEq(t, newContents, oldContents)
				})
			})

			when("builder is a suggested builder", func() {
				it("does nothing", func() {
					h.AssertNil(t, ioutil.WriteFile(configPath, []byte(""), os.ModePerm))

					command.SetArgs([]string{"paketobuildpacks/builder:base"})
					h.AssertNil(t, command.Execute())
					oldContents, err := ioutil.ReadFile(configPath)
					h.AssertNil(t, err)
					h.AssertEq(t, string(oldContents), "")
				})
			})
		})
	})
}
