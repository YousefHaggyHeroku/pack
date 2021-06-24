package fakes

import (
	"context"

	"github.com/YousefHaggyHeroku/pack/internal/builder"
	"github.com/buildpacks/pack/config"
)

type FakeInspectableFetcher struct {
	InspectableToReturn *FakeInspectable
	ErrorToReturn       error

	CallCount int

	ReceivedName       string
	ReceivedDaemon     bool
	ReceivedPullPolicy config.PullPolicy
}

func (f *FakeInspectableFetcher) Fetch(ctx context.Context, name string, daemon bool, pullPolicy config.PullPolicy) (builder.Inspectable, error) {
	f.CallCount++

	f.ReceivedName = name
	f.ReceivedDaemon = daemon
	f.ReceivedPullPolicy = pullPolicy

	return f.InspectableToReturn, f.ErrorToReturn
}
