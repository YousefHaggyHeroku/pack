package fakes

import (
	"github.com/YousefHaggyHeroku/pack/
	"github.com/YousefHaggyHeroku/pack/internal/inspectimage"
	"github.com/YousefHaggyHeroku/pack/logging"
)

type FakeInspectImageWriter struct {
	PrintForLocal  string
	PrintForRemote string
	ErrorForPrint  error

	ReceivedInfoForLocal   *pack.ImageInfo
	ReceivedInfoForRemote  *pack.ImageInfo
	RecievedGeneralInfo    inspectimage.GeneralInfo
	ReceivedErrorForLocal  error
	ReceivedErrorForRemote error
}

func (w *FakeInspectImageWriter) Print(
	logger logging.Logger,
	sharedInfo inspectimage.GeneralInfo,
	local, remote *pack.ImageInfo,
	localErr, remoteErr error,
) error {
	w.ReceivedInfoForLocal = local
	w.ReceivedInfoForRemote = remote
	w.ReceivedErrorForLocal = localErr
	w.ReceivedErrorForRemote = remoteErr
	w.RecievedGeneralInfo = sharedInfo

	logger.Infof("\nLOCAL:\n%s\n", w.PrintForLocal)
	logger.Infof("\nREMOTE:\n%s\n", w.PrintForRemote)

	return w.ErrorForPrint
}
