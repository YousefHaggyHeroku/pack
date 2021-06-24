package commands

import (
	"github.com/spf13/cobra"

	builderwriter "github.com/YousefHaggyHeroku/pack/internal/builder/writer"
	"github.com/YousefHaggyHeroku/pack/internal/config"
	"github.com/buildpacks/pack/logging"
)

func NewBuilderCommand(logger logging.Logger, cfg config.Config, client PackClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "builder",
		Aliases: []string{"builders"},
		Short:   "Interact with builders",
		RunE:    nil,
	}

	cmd.AddCommand(BuilderCreate(logger, cfg, client))
	cmd.AddCommand(BuilderInspect(logger, cfg, client, builderwriter.NewFactory()))
	cmd.AddCommand(BuilderSuggest(logger, client))
	AddHelpFlag(cmd, "builder")
	return cmd
}
