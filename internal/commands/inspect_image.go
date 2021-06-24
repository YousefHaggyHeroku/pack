package commands

import (
	"github.com/spf13/cobra"

	"github.com/YousefHaggyHeroku/pack/internal/inspectimage"

	"github.com/YousefHaggyHeroku/pack/internal/inspectimage/writer"

	"github.com/YousefHaggyHeroku/pack/internal/config"
	"github.com/buildpacks/pack/logging"
)

//go:generate mockgen -package testmocks -destination testmocks/mock_inspect_image_writer_factory.go github.com/YousefHaggyHeroku/pack/internal/commands InspectImageWriterFactory
type InspectImageWriterFactory interface {
	Writer(kind string, BOM bool) (writer.InspectImageWriter, error)
}

type InspectImageFlags struct {
	BOM          bool
	OutputFormat string
}

func InspectImage(
	logger logging.Logger,
	writerFactory InspectImageWriterFactory,
	cfg config.Config,
	client PackClient,
) *cobra.Command {
	var flags InspectImageFlags
	cmd := &cobra.Command{
		Use:     "inspect <image-name>",
		Args:    cobra.ExactArgs(1),
		Aliases: []string{"inspect-image"},
		Short:   "Show information about a built app image",
		Example: "pack inspect buildpacksio/pack",
		RunE: logError(logger, func(cmd *cobra.Command, args []string) error {
			img := args[0]

			sharedImageInfo := inspectimage.GeneralInfo{
				Name:            img,
				RunImageMirrors: cfg.RunImages,
			}

			w, err := writerFactory.Writer(flags.OutputFormat, flags.BOM)
			if err != nil {
				return err
			}

			remote, remoteErr := client.InspectImage(img, false)
			local, localErr := client.InspectImage(img, true)

			if err := w.Print(logger, sharedImageInfo, local, remote, localErr, remoteErr); err != nil {
				return err
			}
			return nil
		}),
	}
	AddHelpFlag(cmd, "inspect")
	cmd.Flags().BoolVar(&flags.BOM, "bom", false, "print bill of materials")
	cmd.Flags().StringVarP(&flags.OutputFormat, "output", "o", "human-readable", "Output format to display builder detail (json, yaml, toml, human-readable).\nOmission of this flag will display as human-readable.")
	return cmd
}
