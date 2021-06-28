package layer

import (
	"github.com/YousefHaggyHeroku/pack/internal/archive"
)

func CreateSingleFileTar(tarFile, path, txt string, twf archive.TarWriterFactory) error {
	tarBuilder := archive.TarBuilder{}
	tarBuilder.AddFile(path, 0644, archive.NormalizedDateTime, []byte(txt))
	return tarBuilder.WriteToPath(tarFile, twf)
}
