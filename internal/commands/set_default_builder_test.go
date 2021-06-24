package commands_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/heroku/color"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	"github.com/spf13/cobra"

	"github.com/YousefHaggyHeroku/pack/internal/commands"
	"github.com/YousefHaggyHeroku/pack/internal/commands/testmocks"
	"github.com/YousefHaggyHeroku/pack/internal/config"
	ilogging "github.com/YousefHaggyHeroku/pack/internal/logging"
	"github.com/buildpacks/pack"
	"github.com/buildpacks/pack/logging"
	h "github.com/buildpacks/pack/testhelpers"
)

func TestSetDefaultBuilderCommand(t *testing.T) {
	color.Disable(true)
	defer color.Disable(false)
	spec.Run(t, "Commands", testSetDefaultBuilderCommand, spec.Random(), spec.Report(report.Terminal{}))
}

func testSetDefaultBuilderCommand(t *testing.T, when spec.G, it spec.S) {
	var (
		command        *cobra.Command
		logger         logging.Logger
		outBuf         bytes.Buffer
		mockController *gomock.Controller
		mockClient     *testmocks.MockPackClient
		tempPackHome   string
	)

	it.Before(func() {
		mockController = gomock.NewController(t)
		mockClient = testmocks.NewMockPackClient(mockController)
		logger = ilogging.NewLogWithWriters(&outBuf, &outBuf)

		var err error
		tempPackHome, err = ioutil.TempDir("", "pack-home")
		h.AssertNil(t, err)
		command = commands.SetDefaultBuilder(logger, config.Config{}, filepath.Join(tempPackHome, "config.toml"), mockClient)
	})

	it.After(func() {
		mockController.Finish()
		h.AssertNil(t, os.RemoveAll(tempPackHome))
	})

	when("#SetDefaultBuilder", func() {
		when("no builder provided", func() {
			it.Before(func() {
				mockClient.EXPECT().InspectBuilder(gomock.Any(), false).Return(&pack.BuilderInfo{}, nil).AnyTimes()
			})

			it("display suggested builders", func() {
				command.SetArgs([]string{})
				h.AssertNil(t, command.Execute())
				h.AssertContains(t, outBuf.String(), "Suggested builders:")
			})
		})

		when("empty builder name is provided", func() {
			it.Before(func() {
				mockClient.EXPECT().InspectBuilder(gomock.Any(), false).Return(&pack.BuilderInfo{}, nil).AnyTimes()
			})

			it("display suggested builders", func() {
				command.SetArgs([]string{})
				h.AssertNil(t, command.Execute())
				h.AssertContains(t, outBuf.String(), "Suggested builders:")
			})
		})

		when("valid builder is provider", func() {
			when("in local", func() {
				it("sets default builder", func() {
					imageName := "some/image"
					mockClient.EXPECT().InspectBuilder(imageName, true).Return(&pack.BuilderInfo{
						Stack: "test.stack.id",
					}, nil)

					command.SetArgs([]string{imageName})
					h.AssertNil(t, command.Execute())
					h.AssertContains(t, outBuf.String(), fmt.Sprintf("Builder '%s' is now the default builder", imageName))
				})
			})

			when("in remote", func() {
				it("sets default builder", func() {
					imageName := "some/image"

					localCall := mockClient.EXPECT().InspectBuilder(imageName, true).Return(nil, nil)

					mockClient.EXPECT().InspectBuilder(imageName, false).Return(&pack.BuilderInfo{
						Stack: "test.stack.id",
					}, nil).After(localCall)

					command.SetArgs([]string{imageName})
					h.AssertNil(t, command.Execute())
					h.AssertContains(t, outBuf.String(), fmt.Sprintf("Builder '%s' is now the default builder", imageName))
				})
			})
		})

		when("invalid builder is provided", func() {
			it("error is presented", func() {
				imageName := "nonbuilder/image"

				mockClient.EXPECT().InspectBuilder(imageName, true).Return(
					nil,
					fmt.Errorf("failed to inspect image %s", imageName))

				command.SetArgs([]string{imageName})

				h.AssertNotNil(t, command.Execute())
				h.AssertContains(t, outBuf.String(), "ERROR: failed to inspect image nonbuilder/image")
			})
		})

		when("non-existent builder is provided", func() {
			it("error is present", func() {
				imageName := "nonexisting/image"

				localCall := mockClient.EXPECT().InspectBuilder(imageName, true).Return(
					nil,
					nil)

				mockClient.EXPECT().InspectBuilder(imageName, false).Return(
					nil,
					nil).After(localCall)

				command.SetArgs([]string{imageName})

				h.AssertNotNil(t, command.Execute())
				h.AssertContains(t, outBuf.String(), "ERROR: builder 'nonexisting/image' not found")
			})
		})
	})
}
