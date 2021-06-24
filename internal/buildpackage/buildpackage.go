package buildpackage

import (
	"github.com/YousefHaggyHeroku/pack/internal/dist"
)

const MetadataLabel = "io.buildpacks.buildpackage.metadata"

type Metadata struct {
	dist.BuildpackInfo
	Stacks []dist.Stack `toml:"stacks" json:"stacks"`
}
