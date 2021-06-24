package commands

import (
	"context"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	pubbldpkg "github.com/YousefHaggyHeroku/pack/buildpackage"
	pubcfg "github.com/YousefHaggyHeroku/pack/config"
	"github.com/YousefHaggyHeroku/pack/internal/config"
	"github.com/YousefHaggyHeroku/pack/internal/style"
	"github.com/YousefHaggyHeroku/pack/logging"
	"github.com/buildpacks/pack"
)

// BuildpackPackageFlags define flags provided to the BuildpackPackage command
type BuildpackPackageFlags struct {
	PackageTomlPath   string
	Format            string
	Publish           bool
	Policy            string
	BuildpackRegistry string
	Path              string
}

// BuildpackPackager packages buildpacks
type BuildpackPackager interface {
	PackageBuildpack(ctx context.Context, options pack.PackageBuildpackOptions) error
}

// PackageConfigReader reads BuildpackPackage configs
type PackageConfigReader interface {
	Read(path string) (pubbldpkg.Config, error)
}

// BuildpackPackage packages (a) buildpack(s) into OCI format, based on a package config
func BuildpackPackage(logger logging.Logger, cfg config.Config, client BuildpackPackager, packageConfigReader PackageConfigReader) *cobra.Command {
	var flags BuildpackPackageFlags
	cmd := &cobra.Command{
		Use:     "package <name> --config <config-path>",
		Short:   "Package a buildpack in OCI format.",
		Args:    cobra.ExactValidArgs(1),
		Example: "pack buildpack package my-buildpack --config ./package.toml\npack buildpack package my-buildpack.cnb --config ./package.toml --f file",
		Long: "buildpack package allows users to package (a) buildpack(s) into OCI format, which can then to be hosted in " +
			"image repositories or persisted on disk as a '.cnb' file. You can also package a number of buildpacks " +
			"together, to enable easier distribution of a set of buildpacks. " +
			"Packaged buildpacks can be used as inputs to `pack build` (using the `--buildpack` flag), " +
			"and they can be included in the configs used in `pack builder create` and `pack buildpack package`. For more " +
			"on how to package a buildpack, see: https://buildpacks.io/docs/buildpack-author-guide/package-a-buildpack/.",
		RunE: logError(logger, func(cmd *cobra.Command, args []string) error {
			if err := validateBuildpackPackageFlags(&flags); err != nil {
				return err
			}

			stringPolicy := flags.Policy
			if stringPolicy == "" {
				stringPolicy = cfg.PullPolicy
			}
			pullPolicy, err := pubcfg.ParsePullPolicy(stringPolicy)
			if err != nil {
				return errors.Wrap(err, "parsing pull policy")
			}
			bpPackageCfg := pubbldpkg.DefaultConfig()
			var bpPath string
			if flags.Path != "" {
				if bpPath, err = filepath.Abs(flags.Path); err != nil {
					return errors.Wrap(err, "resolving buildpack path")
				}
				bpPackageCfg.Buildpack.URI = bpPath
			}
			relativeBaseDir := ""
			if flags.PackageTomlPath != "" {
				bpPackageCfg, err = packageConfigReader.Read(flags.PackageTomlPath)
				if err != nil {
					return errors.Wrap(err, "reading config")
				}

				relativeBaseDir, err = filepath.Abs(filepath.Dir(flags.PackageTomlPath))
				if err != nil {
					return errors.Wrap(err, "getting absolute path for config")
				}
			}
			name := args[0]
			if flags.Format == pack.FormatFile {
				switch ext := filepath.Ext(name); ext {
				case pack.CNBExtension:
				case "":
					name += pack.CNBExtension
				default:
					logger.Warnf("%s is not a valid extension for a packaged buildpack. Packaged buildpacks must have a %s extension", style.Symbol(ext), style.Symbol(pack.CNBExtension))
				}
			}
			if err := client.PackageBuildpack(cmd.Context(), pack.PackageBuildpackOptions{
				RelativeBaseDir: relativeBaseDir,
				Name:            name,
				Format:          flags.Format,
				Config:          bpPackageCfg,
				Publish:         flags.Publish,
				PullPolicy:      pullPolicy,
				Registry:        flags.BuildpackRegistry,
			}); err != nil {
				return err
			}

			action := "created"
			if flags.Publish {
				action = "published"
			}

			logger.Infof("Successfully %s package %s", action, style.Symbol(name))
			return nil
		}),
	}

	cmd.Flags().StringVarP(&flags.PackageTomlPath, "config", "c", "", "Path to package TOML config")
	cmd.Flags().StringVarP(&flags.Format, "format", "f", "", `Format to save package as ("image" or "file")`)
	cmd.Flags().BoolVar(&flags.Publish, "publish", false, `Publish to registry (applies to "--format=image" only)`)
	cmd.Flags().StringVar(&flags.Policy, "pull-policy", "", "Pull policy to use. Accepted values are always, never, and if-not-present. The default is always")
	cmd.Flags().StringVarP(&flags.Path, "path", "p", "", "Path to the Buildpack that needs to be packaged")
	cmd.Flags().StringVarP(&flags.BuildpackRegistry, "buildpack-registry", "r", "", "Buildpack Registry name")

	AddHelpFlag(cmd, "package")
	return cmd
}

func validateBuildpackPackageFlags(p *BuildpackPackageFlags) error {
	if p.Publish && p.Policy == pubcfg.PullNever.String() {
		return errors.Errorf("--publish and --pull-policy never cannot be used together. The --publish flag requires the use of remote images.")
	}
	if p.PackageTomlPath != "" && p.Path != "" {
		return errors.Errorf("--config and --path cannot be used together. Please specify the relative path to the Buildpack directory in the package config file.")
	}

	return nil
}
