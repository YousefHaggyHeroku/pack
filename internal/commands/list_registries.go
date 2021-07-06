package commands

import (
	"github.com/spf13/cobra"

	"github.com/YousefHaggyHeroku/pack/internal/config"
	"github.com/YousefHaggyHeroku/pack/logging"
)

// Deprecated: Use config registries list instead
func ListBuildpackRegistries(logger logging.Logger, cfg config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list-registries",
		Args:    cobra.NoArgs,
		Hidden:  true,
		Short:   "List buildpack registries",
		Example: "pack list-registries",
		RunE: logError(logger, func(cmd *cobra.Command, args []string) error {
			deprecationWarning(logger, "list-registries", "config registries list")
			listRegistries(args, logger, cfg)

			return nil
		}),
	}

	return cmd
}
