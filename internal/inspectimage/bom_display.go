package inspectimage

import (
	"github.com/buildpacks/lifecycle/buildpack"

	"github.com/YousefHaggyHeroku/pack"
	"github.com/YousefHaggyHeroku/pack/internal/dist"
)

type BOMDisplay struct {
	Remote    []BOMEntryDisplay `json:"remote" yaml:"remote"`
	Local     []BOMEntryDisplay `json:"local" yaml:"local"`
	RemoteErr string            `json:"remote_error,omitempty" yaml:"remoteError,omitempty"`
	LocalErr  string            `json:"local_error,omitempty" yaml:"localError,omitempty"`
}

type BOMEntryDisplay struct {
	Name      string                 `toml:"name" json:"name" yaml:"name"`
	Version   string                 `toml:"version,omitempty" json:"version,omitempty" yaml:"version,omitempty"`
	Metadata  map[string]interface{} `toml:"metadata" json:"metadata" yaml:"metadata"`
	Buildpack dist.BuildpackRef      `json:"buildpacks" yaml:"buildpacks" toml:"buildpacks"`
}

func NewBOMDisplay(info *pack.ImageInfo) []BOMEntryDisplay {
	if info == nil {
		return nil
	}
	return displayBOM(info.BOM)
}

func displayBOM(bom []buildpack.BOMEntry) []BOMEntryDisplay {
	var result []BOMEntryDisplay
	for _, entry := range bom {
		result = append(result, BOMEntryDisplay{
			Name:     entry.Name,
			Version:  entry.Version,
			Metadata: entry.Metadata,

			Buildpack: dist.BuildpackRef{
				BuildpackInfo: dist.BuildpackInfo{
					ID:      entry.Buildpack.ID,
					Version: entry.Buildpack.Version,
				},
				Optional: entry.Buildpack.Optional,
			},
		})
	}

	return result
}
