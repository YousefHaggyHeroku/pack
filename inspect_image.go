package pack

import (
	"context"
	"strings"

	"github.com/YousefHaggyHeroku/pack/config"

	"github.com/Masterminds/semver"
	"github.com/buildpacks/lifecycle"
	"github.com/buildpacks/lifecycle/launch"
	"github.com/pkg/errors"

	"github.com/YousefHaggyHeroku/pack/internal/dist"
	"github.com/YousefHaggyHeroku/pack/internal/image"
)

// ImageInfo is a collection of metadata describing
// an app image built using Cloud Native Buildpacks.
type ImageInfo struct {
	// Stack Identifier used when building this image
	StackID string

	// List of buildpacks that passed detection, ran their build
	// phases and made a contribution to this image.
	Buildpacks []lifecycle.Buildpack

	// Base includes two references to the run image,
	// - the Run Image ID,
	// - the hash of the last layer in the app image that belongs to the run image.
	// A way to visualize this is given an image with n layers:
	//
	// last layer in run image
	//          v
	// [1, ..., k, k+1, ..., n]
	//              ^
	//   first layer added by buildpacks
	//
	// the first 1 to k layers all belong to the run image,
	// the last k+1 to n layers are added by buildpacks.
	// the sum of all of these is our app image.
	Base lifecycle.RunImageMetadata

	// BOM or Bill of materials, contains dependency and
	// version information provided by each buildpack.
	BOM []lifecycle.BOMEntry

	// Stack includes the run image name, and a list of image mirrors,
	// where the run image is hosted.
	Stack lifecycle.StackMetadata

	// Processes lists all processes contributed by buildpacks.
	Processes ProcessDetails
}

// ProcessDetails is a collection of all start command metadata
// on an image.
type ProcessDetails struct {
	// An Images default start command.
	DefaultProcess *launch.Process

	// List of all start commands contributed by buildpacks.
	OtherProcesses []launch.Process
}

// Deserialize just the subset of fields we need to avoid breaking changes
type layersMetadata struct {
	RunImage lifecycle.RunImageMetadata `json:"runImage" toml:"run-image"`
	Stack    lifecycle.StackMetadata    `json:"stack" toml:"stack"`
}

const (
	platformAPIEnv            = "CNB_PLATFORM_API"
	cnbProcessEnv             = "CNB_PROCESS_TYPE"
	launcherEntrypoint        = "/cnb/lifecycle/launcher"
	windowsLauncherEntrypoint = `c:\cnb\lifecycle\launcher.exe`
	entrypointPrefix          = "/cnb/process/"
	windowsEntrypointPrefix   = `c:\cnb\process\`
	defaultProcess            = "web"
	fallbackPlatformAPI       = "0.3"
	windowsPrefix             = "c:"
)

// InspectImage reads the Label metadata of an image. It initializes a ImageInfo object
// using this metadata, and returns it.
// If daemon is true, first the local registry will be searched for the image.
// Otherwise it assumes the image is remote.
func (c *Client) InspectImage(name string, daemon bool) (*ImageInfo, error) {
	img, err := c.imageFetcher.Fetch(context.Background(), name, daemon, config.PullNever)
	if err != nil {
		if errors.Cause(err) == image.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	var layersMd layersMetadata
	if _, err := dist.GetLabel(img, lifecycle.LayerMetadataLabel, &layersMd); err != nil {
		return nil, err
	}

	var buildMD lifecycle.BuildMetadata
	if _, err := dist.GetLabel(img, lifecycle.BuildMetadataLabel, &buildMD); err != nil {
		return nil, err
	}

	minimumBaseImageReferenceVersion := semver.MustParse("0.5.0")
	actualLauncherVersion, err := semver.NewVersion(buildMD.Launcher.Version)

	if err == nil && actualLauncherVersion.LessThan(minimumBaseImageReferenceVersion) {
		layersMd.RunImage.Reference = ""
	}

	stackID, err := img.Label(lifecycle.StackIDLabel)
	if err != nil {
		return nil, err
	}

	platformAPI, err := img.Env(platformAPIEnv)
	if err != nil {
		return nil, errors.Wrap(err, "reading platform api")
	}

	if platformAPI == "" {
		platformAPI = fallbackPlatformAPI
	}

	platformAPIVersion, err := semver.NewVersion(platformAPI)
	if err != nil {
		return nil, errors.Wrap(err, "parsing platform api version")
	}

	var defaultProcessType string
	if platformAPIVersion.LessThan(semver.MustParse("0.4")) {
		defaultProcessType, err = img.Env(cnbProcessEnv)
		if err != nil || defaultProcessType == "" {
			defaultProcessType = defaultProcess
		}
	} else {
		inspect, _, err := c.docker.ImageInspectWithRaw(context.TODO(), name)
		if err != nil {
			return nil, errors.Wrap(err, "reading image")
		}

		entrypoint := inspect.Config.Entrypoint
		if len(entrypoint) > 0 && entrypoint[0] != launcherEntrypoint && entrypoint[0] != windowsLauncherEntrypoint {
			process := entrypoint[0]
			if strings.HasPrefix(process, windowsPrefix) {
				process = strings.TrimPrefix(process, windowsEntrypointPrefix)
				process = strings.TrimSuffix(process, ".exe") // Trim .exe for Windows support
			} else {
				process = strings.TrimPrefix(process, entrypointPrefix)
			}

			defaultProcessType = process
		}
	}

	var processDetails ProcessDetails
	for _, proc := range buildMD.Processes {
		proc := proc
		if proc.Type == defaultProcessType {
			processDetails.DefaultProcess = &proc
			continue
		}
		processDetails.OtherProcesses = append(processDetails.OtherProcesses, proc)
	}

	return &ImageInfo{
		StackID:    stackID,
		Stack:      layersMd.Stack,
		Base:       layersMd.RunImage,
		BOM:        buildMD.BOM,
		Buildpacks: buildMD.Buildpacks,
		Processes:  processDetails,
	}, nil
}
