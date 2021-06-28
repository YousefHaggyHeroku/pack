package commands

import (
	"github.com/spf13/cobra"

	"github.com/YousefHaggyHeroku/pack/internal/commands/stack"

	"github.com/YousefHaggyHeroku/pack/logging"
)

func SuggestStacks(logger logging.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "suggest-stacks",
		Args:    cobra.NoArgs,
		Short:   "Display list of recommended stacks",
		Example: "pack suggest-stacks",
		Run: func(*cobra.Command, []string) {
			logger.Warn("Command 'pack suggest-stacks' has been deprecated, please use 'pack stack suggest' instead")
			stack.Suggest(logger)
		},
		Hidden: true,
	}

	AddHelpFlag(cmd, "suggest-stacks")
	return cmd
}
