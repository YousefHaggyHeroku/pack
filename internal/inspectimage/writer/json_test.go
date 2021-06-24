package writer_test

import (
	"bytes"
	"testing"

	"github.com/buildpacks/lifecycle/buildpack"
	"github.com/buildpacks/lifecycle/launch"
	"github.com/buildpacks/lifecycle/platform"
	"github.com/heroku/color"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	"github.com/YousefHaggyHeroku/pack/internal/config"
	"github.com/YousefHaggyHeroku/pack/internal/inspectimage"
	"github.com/YousefHaggyHeroku/pack/internal/inspectimage/writer"
	ilogging "github.com/YousefHaggyHeroku/pack/internal/logging"
	"github.com/buildpacks/pack"
	h "github.com/buildpacks/pack/testhelpers"
)

func TestJSON(t *testing.T) {
	color.Disable(true)
	defer color.Disable(false)
	spec.Run(t, "JSON Writer", testJSON, spec.Parallel(), spec.Report(report.Terminal{}))
}

func testJSON(t *testing.T, when spec.G, it spec.S) {
	var (
		assert = h.NewAssertionManager(t)
		outBuf bytes.Buffer

		remoteInfo *pack.ImageInfo
		localInfo  *pack.ImageInfo

		expectedLocalOutput = `{
  "local_info": {
    "stack": "test.stack.id.local",
    "base_image": {
      "top_layer": "some-local-top-layer",
      "reference": "some-local-run-image-reference"
    },
    "run_images": [
      {
        "name": "user-configured-mirror-for-local",
        "user_configured": true
      },
      {
        "name": "some-local-run-image"
      },
      {
        "name": "some-local-mirror"
      },
      {
        "name": "other-local-mirror"
      }
    ],
    "buildpacks": [
      {
        "homepage": "https://some-homepage-one",
        "id": "test.bp.one.local",
        "version": "1.0.0"
      },
      {
        "homepage": "https://some-homepage-two",
        "id": "test.bp.two.local",
        "version": "2.0.0"
      }
    ],
    "processes": [
      {
        "type": "some-local-type",
        "shell": "bash",
        "command": "/some/local command",
        "default": true,
        "args": [
          "some",
          "local",
          "args"
        ]
      },
      {
        "type": "other-local-type",
        "shell": "",
        "command": "/other/local/command",
        "default": false,
        "args": [
          "other",
          "local",
          "args"
        ]
      }
    ]
  }
}`
		expectedRemoteOutput = `{  
  "remote_info": {
    "stack": "test.stack.id.remote",
    "base_image": {
      "top_layer": "some-remote-top-layer",
      "reference": "some-remote-run-image-reference"
    },
    "run_images": [
      {
        "name": "user-configured-mirror-for-remote",
        "user_configured": true
      },
      {
        "name": "some-remote-run-image"
      },
      {
        "name": "some-remote-mirror"
      },
      {
        "name": "other-remote-mirror"
      }
    ],
    "buildpacks": [
      {
        "id": "test.bp.one.remote",
        "version": "1.0.0",
        "homepage": "https://some-homepage-one"
      },
      {
        "id": "test.bp.two.remote",
        "version": "2.0.0",
        "homepage": "https://some-homepage-two"
      }
    ],
    "processes": [
      {
        "type": "some-remote-type",
        "shell": "bash",
        "command": "/some/remote command",
        "default": true,
        "args": [
          "some",
          "remote",
          "args"
        ]
      },
      {
        "type": "other-remote-type",
        "shell": "",
        "command": "/other/remote/command",
        "default": false,
        "args": [
          "other",
          "remote",
          "args"
        ]
      }
    ]
  }
}`
	)

	when("Print", func() {
		it.Before(func() {
			type someData struct {
				String string
				Bool   bool
				Int    int
				Nested struct {
					String string
				}
			}

			remoteInfo = &pack.ImageInfo{
				StackID: "test.stack.id.remote",
				Buildpacks: []buildpack.GroupBuildpack{
					{ID: "test.bp.one.remote", Version: "1.0.0", Homepage: "https://some-homepage-one"},
					{ID: "test.bp.two.remote", Version: "2.0.0", Homepage: "https://some-homepage-two"},
				},
				Base: platform.RunImageMetadata{
					TopLayer:  "some-remote-top-layer",
					Reference: "some-remote-run-image-reference",
				},
				Stack: platform.StackMetadata{
					RunImage: platform.StackRunImageMetadata{
						Image:   "some-remote-run-image",
						Mirrors: []string{"some-remote-mirror", "other-remote-mirror"},
					},
				},
				BOM: []buildpack.BOMEntry{{
					Require: buildpack.Require{
						Name:    "name-1",
						Version: "version-1",
						Metadata: map[string]interface{}{
							"RemoteData": someData{
								String: "aString",
								Bool:   true,
								Int:    123,
								Nested: struct {
									String string
								}{
									String: "anotherString",
								},
							},
						},
					},
					Buildpack: buildpack.GroupBuildpack{ID: "test.bp.one.remote", Version: "1.0.0", Homepage: "https://some-homepage-one"},
				}},
				Processes: pack.ProcessDetails{
					DefaultProcess: &launch.Process{
						Type:    "some-remote-type",
						Command: "/some/remote command",
						Args:    []string{"some", "remote", "args"},
						Direct:  false,
					},
					OtherProcesses: []launch.Process{
						{
							Type:    "other-remote-type",
							Command: "/other/remote/command",
							Args:    []string{"other", "remote", "args"},
							Direct:  true,
						},
					},
				},
			}

			localInfo = &pack.ImageInfo{
				StackID: "test.stack.id.local",
				Buildpacks: []buildpack.GroupBuildpack{
					{ID: "test.bp.one.local", Version: "1.0.0", Homepage: "https://some-homepage-one"},
					{ID: "test.bp.two.local", Version: "2.0.0", Homepage: "https://some-homepage-two"},
				},
				Base: platform.RunImageMetadata{
					TopLayer:  "some-local-top-layer",
					Reference: "some-local-run-image-reference",
				},
				Stack: platform.StackMetadata{
					RunImage: platform.StackRunImageMetadata{
						Image:   "some-local-run-image",
						Mirrors: []string{"some-local-mirror", "other-local-mirror"},
					},
				},
				BOM: []buildpack.BOMEntry{{
					Require: buildpack.Require{
						Name:    "name-1",
						Version: "version-1",
						Metadata: map[string]interface{}{
							"LocalData": someData{
								Bool: false,
								Int:  456,
							},
						},
					},
					Buildpack: buildpack.GroupBuildpack{ID: "test.bp.one.remote", Version: "1.0.0", Homepage: "https://some-homepage-one"},
				}},
				Processes: pack.ProcessDetails{
					DefaultProcess: &launch.Process{
						Type:    "some-local-type",
						Command: "/some/local command",
						Args:    []string{"some", "local", "args"},
						Direct:  false,
					},
					OtherProcesses: []launch.Process{
						{
							Type:    "other-local-type",
							Command: "/other/local/command",
							Args:    []string{"other", "local", "args"},
							Direct:  true,
						},
					},
				},
			}

			outBuf = bytes.Buffer{}
		})

		when("local and remote image exits", func() {
			it("prints both local and remote image info in a JSON format", func() {
				runImageMirrors := []config.RunImage{
					{
						Image:   "un-used-run-image",
						Mirrors: []string{"un-used"},
					},
					{
						Image:   "some-local-run-image",
						Mirrors: []string{"user-configured-mirror-for-local"},
					},
					{
						Image:   "some-remote-run-image",
						Mirrors: []string{"user-configured-mirror-for-remote"},
					},
				}
				sharedImageInfo := inspectimage.GeneralInfo{
					Name:            "test-image",
					RunImageMirrors: runImageMirrors,
				}
				jsonWriter := writer.NewJSON()

				logger := ilogging.NewLogWithWriters(&outBuf, &outBuf)
				err := jsonWriter.Print(logger, sharedImageInfo, localInfo, remoteInfo, nil, nil)
				assert.Nil(err)

				assert.ContainsJSON(outBuf.String(), `{ "image_name": "test-image" }`)
				assert.ContainsJSON(outBuf.String(), expectedLocalOutput)
				assert.ContainsJSON(outBuf.String(), expectedRemoteOutput)
			})
		})

		when("only local image exists", func() {
			it("prints local image info in JSON format", func() {
				runImageMirrors := []config.RunImage{
					{
						Image:   "un-used-run-image",
						Mirrors: []string{"un-used"},
					},
					{
						Image:   "some-local-run-image",
						Mirrors: []string{"user-configured-mirror-for-local"},
					},
					{
						Image:   "some-remote-run-image",
						Mirrors: []string{"user-configured-mirror-for-remote"},
					},
				}
				sharedImageInfo := inspectimage.GeneralInfo{
					Name:            "test-image",
					RunImageMirrors: runImageMirrors,
				}
				jsonWriter := writer.NewJSON()

				logger := ilogging.NewLogWithWriters(&outBuf, &outBuf)
				err := jsonWriter.Print(logger, sharedImageInfo, localInfo, nil, nil, nil)
				assert.Nil(err)

				assert.ContainsJSON(outBuf.String(), `{ "image_name": "test-image" }`)
				assert.ContainsJSON(outBuf.String(), expectedLocalOutput)

				assert.NotContains(outBuf.String(), "test.stack.id.remote")
				assert.ContainsJSON(outBuf.String(), expectedLocalOutput)
			})
		})

		when("only remote image exists", func() {
			it("prints remote image info in JSON format", func() {
				runImageMirrors := []config.RunImage{
					{
						Image:   "un-used-run-image",
						Mirrors: []string{"un-used"},
					},
					{
						Image:   "some-local-run-image",
						Mirrors: []string{"user-configured-mirror-for-local"},
					},
					{
						Image:   "some-remote-run-image",
						Mirrors: []string{"user-configured-mirror-for-remote"},
					},
				}
				sharedImageInfo := inspectimage.GeneralInfo{
					Name:            "test-image",
					RunImageMirrors: runImageMirrors,
				}
				jsonWriter := writer.NewJSON()

				logger := ilogging.NewLogWithWriters(&outBuf, &outBuf)
				err := jsonWriter.Print(logger, sharedImageInfo, nil, remoteInfo, nil, nil)
				assert.Nil(err)

				assert.ContainsJSON(outBuf.String(), `{ "image_name": "test-image" }`)
				assert.NotContains(outBuf.String(), "test.stack.id.local")
				assert.ContainsJSON(outBuf.String(), expectedRemoteOutput)
			})
		})
	})
}
