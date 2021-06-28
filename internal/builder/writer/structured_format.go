package writer

import (
	"fmt"

	"github.com/YousefHaggyHeroku/pack/internal/style"

	"github.com/YousefHaggyHeroku/pack"
	pubbldr "github.com/YousefHaggyHeroku/pack/builder"
	"github.com/YousefHaggyHeroku/pack/internal/builder"
	"github.com/YousefHaggyHeroku/pack/internal/config"
	"github.com/YousefHaggyHeroku/pack/internal/dist"
	"github.com/YousefHaggyHeroku/pack/logging"
)

type InspectOutput struct {
	SharedBuilderInfo
	RemoteInfo *BuilderInfo `json:"remote_info" yaml:"remote_info" toml:"remote_info"`
	LocalInfo  *BuilderInfo `json:"local_info" yaml:"local_info" toml:"local_info"`
}

type RunImage struct {
	Name           string `json:"name" yaml:"name" toml:"name"`
	UserConfigured bool   `json:"user_configured,omitempty" yaml:"user_configured,omitempty" toml:"user_configured,omitempty"`
}

type Lifecycle struct {
	builder.LifecycleInfo `yaml:"lifecycleinfo,inline"`
	BuildpackAPIs         builder.APIVersions `json:"buildpack_apis" yaml:"buildpack_apis" toml:"buildpack_apis"`
	PlatformAPIs          builder.APIVersions `json:"platform_apis" yaml:"platform_apis" toml:"platform_apis"`
}

type Stack struct {
	ID     string   `json:"id" yaml:"id" toml:"id"`
	Mixins []string `json:"mixins,omitempty" yaml:"mixins,omitempty" toml:"mixins,omitempty"`
}

type BuilderInfo struct {
	Description            string                  `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
	CreatedBy              builder.CreatorMetadata `json:"created_by" yaml:"created_by" toml:"created_by"`
	Stack                  Stack                   `json:"stack" yaml:"stack" toml:"stack"`
	Lifecycle              Lifecycle               `json:"lifecycle" yaml:"lifecycle" toml:"lifecycle"`
	RunImages              []RunImage              `json:"run_images" yaml:"run_images" toml:"run_images"`
	Buildpacks             []dist.BuildpackInfo    `json:"buildpacks" yaml:"buildpacks" toml:"buildpacks"`
	pubbldr.DetectionOrder `json:"detection_order" yaml:"detection_order" toml:"detection_order"`
}

type StructuredFormat struct {
	MarshalFunc func(interface{}) ([]byte, error)
}

func (w *StructuredFormat) Print(
	logger logging.Logger,
	localRunImages []config.RunImage,
	local, remote *pack.BuilderInfo,
	localErr, remoteErr error,
	builderInfo SharedBuilderInfo,
) error {
	if localErr != nil {
		return fmt.Errorf("preparing output for %s: %w", style.Symbol(builderInfo.Name), localErr)
	}

	if remoteErr != nil {
		return fmt.Errorf("preparing output for %s: %w", style.Symbol(builderInfo.Name), remoteErr)
	}

	outputInfo := InspectOutput{SharedBuilderInfo: builderInfo}

	if local != nil {
		stack := Stack{ID: local.Stack}

		if logger.IsVerbose() {
			stack.Mixins = local.Mixins
		}

		outputInfo.LocalInfo = &BuilderInfo{
			Description: local.Description,
			CreatedBy:   local.CreatedBy,
			Stack:       stack,
			Lifecycle: Lifecycle{
				LifecycleInfo: local.Lifecycle.Info,
				BuildpackAPIs: local.Lifecycle.APIs.Buildpack,
				PlatformAPIs:  local.Lifecycle.APIs.Platform,
			},
			RunImages:      runImages(local.RunImage, localRunImages, local.RunImageMirrors),
			Buildpacks:     local.Buildpacks,
			DetectionOrder: local.Order,
		}
	}

	if remote != nil {
		stack := Stack{ID: remote.Stack}

		if logger.IsVerbose() {
			stack.Mixins = remote.Mixins
		}

		outputInfo.RemoteInfo = &BuilderInfo{
			Description: remote.Description,
			CreatedBy:   remote.CreatedBy,
			Stack:       stack,
			Lifecycle: Lifecycle{
				LifecycleInfo: remote.Lifecycle.Info,
				BuildpackAPIs: remote.Lifecycle.APIs.Buildpack,
				PlatformAPIs:  remote.Lifecycle.APIs.Platform,
			},
			RunImages:      runImages(remote.RunImage, localRunImages, remote.RunImageMirrors),
			Buildpacks:     remote.Buildpacks,
			DetectionOrder: remote.Order,
		}
	}

	if outputInfo.LocalInfo == nil && outputInfo.RemoteInfo == nil {
		return fmt.Errorf("unable to find builder %s locally or remotely", style.Symbol(builderInfo.Name))
	}

	var (
		output []byte
		err    error
	)
	if output, err = w.MarshalFunc(outputInfo); err != nil {
		return fmt.Errorf("untested, unexpected failure while marshaling: %w", err)
	}

	logger.Info(string(output))

	return nil
}

func runImages(runImage string, localRunImages []config.RunImage, buildRunImages []string) []RunImage {
	var images = []RunImage{}

	for _, i := range localRunImages {
		if i.Image == runImage {
			for _, m := range i.Mirrors {
				images = append(images, RunImage{Name: m, UserConfigured: true})
			}
		}
	}

	if runImage != "" {
		images = append(images, RunImage{Name: runImage})
	}

	for _, m := range buildRunImages {
		images = append(images, RunImage{Name: m})
	}

	return images
}
