package pack

import (
	"context"

	"github.com/buildpacks/pack/config"

	"github.com/pkg/errors"

	"github.com/YousefHaggyHeroku/pack/internal/buildpackage"
	"github.com/YousefHaggyHeroku/pack/internal/dist"
	"github.com/YousefHaggyHeroku/pack/internal/style"
)

var (
	// Version is the version of `pack`. It is injected at compile time.
	Version = "0.0.0"
)

func extractPackagedBuildpacks(ctx context.Context, pkgImageRef string, fetcher ImageFetcher, publish bool, pullPolicy config.PullPolicy) (mainBP dist.Buildpack, depBPs []dist.Buildpack, err error) {
	pkgImage, err := fetcher.Fetch(ctx, pkgImageRef, !publish, pullPolicy)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "fetching image")
	}

	mainBP, depBPs, err = buildpackage.ExtractBuildpacks(pkgImage)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "extracting buildpacks from %s", style.Symbol(pkgImageRef))
	}

	return mainBP, depBPs, nil
}
