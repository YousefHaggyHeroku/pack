package commands

import (
	"github.com/spf13/cobra"

	"github.com/YousefHaggyHeroku/pack/internal/style"
	"github.com/buildpacks/pack"

	"github.com/YousefHaggyHeroku/pack/internal/config"
	"github.com/buildpacks/pack/logging"
)

// Deprecated: Use BuildpackRegister instead
func RegisterBuildpack(logger logging.Logger, cfg config.Config, client PackClient) *cobra.Command {
	var opts pack.RegisterBuildpackOptions
	var flags BuildpackRegisterFlags

	cmd := &cobra.Command{
		Use:     "register-buildpack <image>",
		Hidden:  true,
		Args:    cobra.ExactArgs(1),
		Short:   "Register the buildpack to a registry",
		Example: "pack register-buildpack my-buildpack",
		RunE: logError(logger, func(cmd *cobra.Command, args []string) error {
			deprecationWarning(logger, "register-buildpack", "buildpack register")
			registry, err := config.GetRegistry(cfg, flags.BuildpackRegistry)
			if err != nil {
				return err
			}
			opts.ImageName = args[0]
			opts.Type = registry.Type
			opts.URL = registry.URL
			opts.Name = registry.Name

			if err := client.RegisterBuildpack(cmd.Context(), opts); err != nil {
				return err
			}
			logger.Infof("Successfully registered %s", style.Symbol(opts.ImageName))
			return nil
		}),
	}
	cmd.Flags().StringVarP(&flags.BuildpackRegistry, "buildpack-registry", "r", "", "Buildpack Registry name")
	AddHelpFlag(cmd, "register-buildpack")
	return cmd
}
