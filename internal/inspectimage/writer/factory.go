package writer

import (
	"fmt"

	"github.com/YousefHaggyHeroku/pack/internal/inspectimage"

	"github.com/YousefHaggyHeroku/pack"
	"github.com/YousefHaggyHeroku/pack/logging"

	"github.com/YousefHaggyHeroku/pack/internal/style"
)

type Factory struct{}

type InspectImageWriter interface {
	Print(
		logger logging.Logger,
		sharedInfo inspectimage.GeneralInfo,
		local, remote *pack.ImageInfo,
		localErr, remoteErr error,
	) error
}

func NewFactory() *Factory {
	return &Factory{}
}

func (f *Factory) Writer(kind string, bom bool) (InspectImageWriter, error) {
	if bom {
		switch kind {
		case "human-readable", "json":
			return NewJSONBOM(), nil
		case "yaml":
			return NewYAMLBOM(), nil
		}
	} else {
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
	}

	return nil, fmt.Errorf("output format %s is not supported", style.Symbol(kind))
}
