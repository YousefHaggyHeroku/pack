package writer

import (
	"fmt"

	"github.com/YousefHaggyHeroku/pack"
	"github.com/YousefHaggyHeroku/pack/internal/config"
	"github.com/YousefHaggyHeroku/pack/logging"

	"github.com/YousefHaggyHeroku/pack/internal/style"
)

type Factory struct{}

type BuilderWriter interface {
	Print(
		logger logging.Logger,
		localRunImages []config.RunImage,
		local, remote *pack.BuilderInfo,
		localErr, remoteErr error,
		builderInfo SharedBuilderInfo,
	) error
}

type SharedBuilderInfo struct {
	Name      string `json:"builder_name" yaml:"builder_name" toml:"builder_name"`
	Trusted   bool   `json:"trusted" yaml:"trusted" toml:"trusted"`
	IsDefault bool   `json:"default" yaml:"default" toml:"default"`
}

type BuilderWriterFactory interface {
	Writer(kind string) (BuilderWriter, error)
}

func NewFactory() *Factory {
	return &Factory{}
}

func (f *Factory) Writer(kind string) (BuilderWriter, error) {
	switch kind {
	case "human-readable":
		return NewHumanReadable(), nil
	case "json":
		return NewJSON(), nil
	case "yaml":
		return NewYAML(), nil
	case "toml":
		return NewTOML(), nil
	}

	return nil, fmt.Errorf("output format %s is not supported", style.Symbol(kind))
}
