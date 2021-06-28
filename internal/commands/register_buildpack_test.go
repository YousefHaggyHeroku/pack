package commands_test

import (
	"bytes"
	"testing"

	"github.com/YousefHaggyHeroku/pack/internal/commands"

	"github.com/YousefHaggyHeroku/pack"

	"github.com/golang/mock/gomock"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	"github.com/spf13/cobra"

	"github.com/YousefHaggyHeroku/pack/internal/commands/testmocks"
	"github.com/YousefHaggyHeroku/pack/internal/config"
	ilogging "github.com/YousefHaggyHeroku/pack/internal/logging"
	"github.com/YousefHaggyHeroku/pack/logging"
	h "github.com/YousefHaggyHeroku/pack/testhelpers"
)

func TestRegisterBuildpackCommand(t *testing.T) {
	spec.Run(t, "Commands", testRegisterBuildpackCommand, spec.Parallel(), spec.Report(report.Terminal{}))
}

func testRegisterBuildpackCommand(t *testing.T, when spec.G, it spec.S) {
	var (
		command        *cobra.Command
		logger         logging.Logger
		outBuf         bytes.Buffer
		mockController *gomock.Controller
		mockClient     *testmocks.MockPackClient
		cfg            config.Config
	)

	it.Before(func() {
		logger = ilogging.NewLogWithWriters(&outBuf, &outBuf)
		mockController = gomock.NewController(t)
		mockClient = testmocks.NewMockPackClient(mockController)
		cfg = config.Config{}

		command = commands.RegisterBuildpack(logger, cfg, mockClient)
	})

	it.After(func() {})

	when("#RegisterBuildpackCommand", func() {
		when("no image is provided", func() {
			it("fails to run", func() {
				err := command.Execute()
				h.AssertError(t, err, "accepts 1 arg")
			})
		})

		when("image name is provided", func() {
			var (
				buildpackImage string
			)

			it.Before(func() {
				buildpackImage = "buildpack/image"
			})

			it("should work for required args", func() {
				opts := pack.RegisterBuildpackOptions{
					ImageName: buildpackImage,
					Type:      "github",
					URL:       "https://github.com/buildpacks/registry-index",
					Name:      "official",
				}

				mockClient.EXPECT().
					RegisterBuildpack(gomock.Any(), opts).
					Return(nil)

				command.SetArgs([]string{buildpackImage})
				h.AssertNil(t, command.Execute())
			})

			when("config.toml exists", func() {
				it("should consume registry config values", func() {
					cfg = config.Config{
						DefaultRegistryName: "berneuse",
						Registries: []config.Registry{
							{
								Name: "berneuse",
								Type: "github",
								URL:  "https://github.com/berneuse/buildpack-registry",
							},
						},
					}
					command = commands.RegisterBuildpack(logger, cfg, mockClient)
					opts := pack.RegisterBuildpackOptions{
						ImageName: buildpackImage,
						Type:      "github",
						URL:       "https://github.com/berneuse/buildpack-registry",
						Name:      "berneuse",
					}

					mockClient.EXPECT().
						RegisterBuildpack(gomock.Any(), opts).
						Return(nil)

					command.SetArgs([]string{buildpackImage})
					h.AssertNil(t, command.Execute())
				})

				it("should handle config errors", func() {
					cfg = config.Config{
						DefaultRegistryName: "missing registry",
					}
					command = commands.RegisterBuildpack(logger, cfg, mockClient)
					command.SetArgs([]string{buildpackImage})

					err := command.Execute()
					h.AssertNotNil(t, err)
				})
			})

			it("should support buildpack-registry flag", func() {
				buildpackRegistry := "override"
				cfg = config.Config{
					DefaultRegistryName: "default",
					Registries: []config.Registry{
						{
							Name: "default",
							Type: "github",
							URL:  "https://github.com/default/buildpack-registry",
						},
						{
							Name: "override",
							Type: "github",
							URL:  "https://github.com/override/buildpack-registry",
						},
					},
				}
				opts := pack.RegisterBuildpackOptions{
					ImageName: buildpackImage,
					Type:      "github",
					URL:       "https://github.com/override/buildpack-registry",
					Name:      "override",
				}
				mockClient.EXPECT().
					RegisterBuildpack(gomock.Any(), opts).
					Return(nil)

				command = commands.RegisterBuildpack(logger, cfg, mockClient)
				command.SetArgs([]string{buildpackImage, "--buildpack-registry", buildpackRegistry})
				h.AssertNil(t, command.Execute())
			})
		})
	})
}
