// +build acceptance

package assertions

import (
	"fmt"
	"regexp"
	"testing"

	h "github.com/YousefHaggyHeroku/pack/testhelpers"
)

type LifecycleOutputAssertionManager struct {
	testObject *testing.T
	assert     h.AssertionManager
	output     string
}

func NewLifecycleOutputAssertionManager(t *testing.T, output string) LifecycleOutputAssertionManager {
	return LifecycleOutputAssertionManager{
		testObject: t,
		assert:     h.NewAssertionManager(t),
		output:     output,
	}
}

func (l LifecycleOutputAssertionManager) ReportsRestoresCachedLayer(layer string) {
	l.testObject.Helper()
	l.testObject.Log("restores the cache")

	l.assert.MatchesAll(
		l.output,
		regexp.MustCompile(fmt.Sprintf(`(?i)Restoring data for "%s" from cache`, layer)),
		regexp.MustCompile(fmt.Sprintf(`(?i)Restoring metadata for "%s" from app image`, layer)),
	)
}

func (l LifecycleOutputAssertionManager) ReportsExporterReusingUnchangedLayer(layer string) {
	l.testObject.Helper()
	l.testObject.Log("exporter reuses unchanged layers")

	l.assert.Matches(l.output, regexp.MustCompile(fmt.Sprintf(`(?i)Reusing layer '%s'`, layer)))
}

func (l LifecycleOutputAssertionManager) ReportsCacheReuse(layer string) {
	l.testObject.Helper()
	l.testObject.Log("cacher reuses unchanged layers")

	l.assert.Matches(l.output, regexp.MustCompile(fmt.Sprintf(`(?i)Reusing cache layer '%s'`, layer)))
}

func (l LifecycleOutputAssertionManager) ReportsCacheCreation(layer string) {
	l.testObject.Helper()
	l.testObject.Log("cacher adds layers")

	l.assert.Matches(l.output, regexp.MustCompile(fmt.Sprintf(`(?i)Adding cache layer '%s'`, layer)))
}

func (l LifecycleOutputAssertionManager) ReportsSkippingBuildpackLayerAnalysis() {
	l.testObject.Helper()
	l.testObject.Log("skips buildpack layer analysis")

	l.assert.Matches(l.output, regexp.MustCompile(`(?i)Skipping buildpack layer analysis`))
}

func (l LifecycleOutputAssertionManager) IncludesSeparatePhases() {
	l.testObject.Helper()

	l.assert.ContainsAll(l.output, "[detector]", "[analyzer]", "[builder]", "[exporter]")
}

func (l LifecycleOutputAssertionManager) IncludesLifecycleImageTag(tag string) {
	l.testObject.Helper()

	l.assert.Contains(l.output, tag)
}
