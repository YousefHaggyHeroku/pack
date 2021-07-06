package pack

import (
	"context"

	"github.com/pkg/errors"

	"github.com/YousefHaggyHeroku/pack/image"
	"github.com/YousefHaggyHeroku/pack/internal/buildpackage"
	"github.com/YousefHaggyHeroku/pack/internal/dist"
	"github.com/YousefHaggyHeroku/pack/internal/style"
)

var (
	// Version is the version of `pack`. It is injected at compile time.
	Version = "0.0.0"
)

func extractPackagedBuildpacks(ctx context.Context, pkgImageRef string, fetcher ImageFetcher, fetchOptions image.FetchOptions) (mainBP dist.Buildpack, depBPs []dist.Buildpack, err error) {
	pkgImage, err := fetcher.Fetch(ctx, pkgImageRef, fetchOptions)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "fetching image")
	}

	mainBP, depBPs, err = buildpackage.ExtractBuildpacks(pkgImage)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "extracting buildpacks from %s", style.Symbol(pkgImageRef))
	}

	return mainBP, depBPs, nil
}
