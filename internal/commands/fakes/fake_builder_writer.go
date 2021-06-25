package fakes

import (
	"github.com/YousefHaggyHeroku/pack/
	"github.com/YousefHaggyHeroku/pack/internal/builder/writer"
	"github.com/YousefHaggyHeroku/pack/internal/config"
	"github.com/YousefHaggyHeroku/pack/logging"
)

type FakeBuilderWriter struct {
	PrintForLocal  string
	PrintForRemote string
	ErrorForPrint  error

	ReceivedInfoForLocal   *pack.BuilderInfo
	ReceivedInfoForRemote  *pack.BuilderInfo
	ReceivedErrorForLocal  error
	ReceivedErrorForRemote error
	ReceivedBuilderInfo    writer.SharedBuilderInfo
	ReceivedLocalRunImages []config.RunImage
}

func (w *FakeBuilderWriter) Print(
	logger logging.Logger,
	localRunImages []config.RunImage,
	local, remote *pack.BuilderInfo,
	localErr, remoteErr error,
	builderInfo writer.SharedBuilderInfo,
) error {
	w.ReceivedInfoForLocal = local
	w.ReceivedInfoForRemote = remote
	w.ReceivedErrorForLocal = localErr
	w.ReceivedErrorForRemote = remoteErr
	w.ReceivedBuilderInfo = builderInfo
	w.ReceivedLocalRunImages = localRunImages

	logger.Infof("\nLOCAL:\n%s\n", w.PrintForLocal)
	logger.Infof("\nREMOTE:\n%s\n", w.PrintForRemote)

	return w.ErrorForPrint
}
