package strings_test

import (
	"testing"

	"github.com/YousefHaggyHeroku/pack/internal/strings"

	"github.com/sclevine/spec"

	h "github.com/buildpacks/pack/testhelpers"
)

func TestValueOrDefault(t *testing.T) {
	spec.Run(t, "Strings", func(t *testing.T, when spec.G, it spec.S) {
		var (
			assert = h.NewAssertionManager(t)
		)

		when("#ValueOrDefault", func() {
			it("returns value when value is non-empty", func() {
				output := strings.ValueOrDefault("some-value", "-")
				assert.Equal(output, "some-value")
			})

			it("returns default when value is empty", func() {
				output := strings.ValueOrDefault("", "-")
				assert.Equal(output, "-")
			})
		})
	})
}
