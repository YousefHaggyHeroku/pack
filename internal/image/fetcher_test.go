package image_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/buildpacks/imgutil/local"
	"github.com/buildpacks/imgutil/remote"
	"github.com/docker/docker/client"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/heroku/color"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	pubcfg "github.com/buildpacks/pack/config"
	"github.com/buildpacks/pack/internal/image"
	"github.com/buildpacks/pack/internal/logging"
	h "github.com/buildpacks/pack/testhelpers"
)

var docker client.CommonAPIClient
var dockerRegistry *h.DockerRegistry

func TestFetcher(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())

	color.Disable(true)
	defer color.Disable(false)

	dockerConfigDir, err := ioutil.TempDir("", "test.docker.config.dir")
	h.AssertNil(t, err)
	defer os.RemoveAll(dockerConfigDir)

	dockerRegistry = h.NewDockerRegistry(h.WithAuth(dockerConfigDir))
	dockerRegistry.Start(t)
	defer dockerRegistry.Stop(t)

	os.Setenv("DOCKER_CONFIG", dockerRegistry.DockerDirectory)
	defer os.Unsetenv("DOCKER_CONFIG")


	// h.RequireDocker(t)
	//
	// registryConfig = h.RunRegistry(t)
	// defer registryConfig.StopRegistry(t)
	//
	// // TODO: is there a better solution to the auth problem?
	// os.Setenv("DOCKER_CONFIG", registryConfig.DockerConfigDir)
	//
	// var err error
	// docker, err = client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.38"))
	// h.AssertNil(t, err)
	spec.Run(t, "Fetcher", testFetcher, spec.Report(report.Terminal{}))
}

func testFetcher(t *testing.T, when spec.G, it spec.S) {
	var (
		fetcher  *image.Fetcher
		repoName string
		repo     string
		outBuf   bytes.Buffer
	)

	it.Before(func() {
		docker = h.DockerCli(t)
		repo = "some-org/" + h.RandString(10)
		repoName = dockerRegistry.RepoName(repo)
		fetcher = image.NewFetcher(logging.NewLogWithWriters(&outBuf, &outBuf), docker)
	})

	when("#Fetch", func() {
		when("daemon is false", func() {
			when("PullAlways", func() {
				when("there is a remote image", func() {
					it.Before(func() {
						img, err := remote.NewImage(repoName, authn.DefaultKeychain)
						h.AssertNil(t, err)

						h.AssertNil(t, img.Save())
					})

					it("returns the remote image", func() {
						_, err := fetcher.Fetch(context.TODO(), repoName, false, pubcfg.PullAlways)
						h.AssertNil(t, err)
					})
				})

				when("there is no remote image", func() {
					it("returns an error", func() {
						_, err := fetcher.Fetch(context.TODO(), repoName, false, pubcfg.PullAlways)
						h.AssertError(t, err, fmt.Sprintf("image '%s' does not exist in registry", repoName))
					})
				})
			})

			when("PullIfNotPresent", func() {
				when("there is a remote image", func() {
					it.Before(func() {
						img, err := remote.NewImage(repoName, authn.DefaultKeychain)
						h.AssertNil(t, err)

						h.AssertNil(t, img.Save())
					})

					it("returns the remote image", func() {
						_, err := fetcher.Fetch(context.TODO(), repoName, false, pubcfg.PullIfNotPresent)
						h.AssertNil(t, err)
					})
				})

				when("there is no remote image", func() {
					it("returns an error", func() {
						_, err := fetcher.Fetch(context.TODO(), repoName, false, pubcfg.PullIfNotPresent)
						h.AssertError(t, err, fmt.Sprintf("image '%s' does not exist in registry", repoName))
					})
				})
			})
		})

		when("daemon is true", func() {
			when("PullNever", func() {
				when("there is a local image", func() {
					it.Before(func() {
						// Make sure the repoName is not a valid remote repo.
						// This is to verify that no remote check is made
						// when there's a valid local image.
						repoName = "invalidhost" + repoName

						img, err := local.NewImage(repoName, docker)
						h.AssertNil(t, err)

						h.AssertNil(t, img.Save())
					})

					it.After(func() {
						h.DockerRmi(docker, repoName)
					})

					it("returns the local image", func() {
						_, err := fetcher.Fetch(context.TODO(), repoName, true, pubcfg.PullNever)
						h.AssertNil(t, err)
					})
				})

				when("there is no local image", func() {
					it("returns an error", func() {
						_, err := fetcher.Fetch(context.TODO(), repoName, true, pubcfg.PullNever)
						h.AssertError(t, err, fmt.Sprintf("image '%s' does not exist on the daemon", repoName))
					})
				})
			})

			when("PullAlways", func() {
				when("there is a remote image", func() {
					var (
						logger *logging.LogWithWriters
						output func() string
					)

					it.Before(func() {
						// Instantiate a pull-able local image
						// as opposed to a remote image so that the image
						// is created with the OS of the docker daemon
						img, err := local.NewImage(repoName, docker)
						h.AssertNil(t, err)
						defer h.DockerRmi(docker, repoName)

						h.AssertNil(t, img.Save())

						h.PushImage(t, docker, repoName)

						var outCons *color.Console
						outCons, output = h.MockWriterAndOutput()
						logger = logging.NewLogWithWriters(outCons, outCons)
						fetcher = image.NewFetcher(logger, docker)
					})

					it.After(func() {
						h.DockerRmi(docker, repoName)
					})

					it("pull the image and return the local copy", func() {
						_, err := fetcher.Fetch(context.TODO(), repoName, true, pubcfg.PullAlways)
						h.AssertNil(t, err)
						h.AssertNotEq(t, output(), "")
					})

					it("doesn't log anything in quiet mode", func() {
						logger.WantQuiet(true)
						_, err := fetcher.Fetch(context.TODO(), repoName, true, pubcfg.PullAlways)
						h.AssertNil(t, err)
						h.AssertEq(t, output(), "")
					})
				})

				when("there is no remote image", func() {
					when("there is a local image", func() {
						it.Before(func() {
							img, err := local.NewImage(repoName, docker)
							h.AssertNil(t, err)

							h.AssertNil(t, img.Save())
						})

						it.After(func() {
							h.DockerRmi(docker, repoName)
						})

						it("returns the local image", func() {
							_, err := fetcher.Fetch(context.TODO(), repoName, true, pubcfg.PullAlways)
							h.AssertNil(t, err)
						})
					})

					when("there is no local image", func() {
						it("returns an error", func() {
							_, err := fetcher.Fetch(context.TODO(), repoName, true, pubcfg.PullAlways)
							h.AssertError(t, err, fmt.Sprintf("image '%s' does not exist on the daemon", repoName))
						})
					})
				})
			})

			when("PullIfNotPresent", func() {
				when("there is a remote image", func() {
					var (
						label          = "label"
						remoteImgLabel string
					)

					it.Before(func() {
						// Instantiate a pull-able local image
						// as opposed to a remote image so that the image
						// is created with the OS of the docker daemon
						remoteImg, err := local.NewImage(repoName, docker)
						h.AssertNil(t, err)
						defer h.DockerRmi(docker, repoName)

						h.AssertNil(t, remoteImg.SetLabel(label, "1"))
						h.AssertNil(t, remoteImg.Save())

						h.PushImage(t, docker, remoteImg.Name())

						remoteImgLabel, err = remoteImg.Label(label)
						h.AssertNil(t, err)
					})

					it.After(func() {
						h.DockerRmi(docker, repoName)
					})

					when("there is a local image", func() {
						var localImgLabel string

						it.Before(func() {
							localImg, err := local.NewImage(repoName, docker)
							h.AssertNil(t, localImg.SetLabel(label, "2"))
							h.AssertNil(t, err)

							h.AssertNil(t, localImg.Save())

							localImgLabel, err = localImg.Label(label)
							h.AssertNil(t, err)
						})

						it.After(func() {
							h.DockerRmi(docker, repoName)
						})

						it("returns the local image", func() {
							fetchedImg, err := fetcher.Fetch(context.TODO(), repoName, true, pubcfg.PullIfNotPresent)
							h.AssertNil(t, err)
							h.AssertNotContains(t, outBuf.String(), "Pulling image")

							fetchedImgLabel, err := fetchedImg.Label(label)
							h.AssertNil(t, err)

							h.AssertEq(t, fetchedImgLabel, localImgLabel)
							h.AssertNotEq(t, fetchedImgLabel, remoteImgLabel)
						})
					})

					when("there is no local image", func() {
						it("returns the remote image", func() {
							fetchedImg, err := fetcher.Fetch(context.TODO(), repoName, true, pubcfg.PullIfNotPresent)
							h.AssertNil(t, err)

							fetchedImgLabel, err := fetchedImg.Label(label)
							h.AssertNil(t, err)
							h.AssertEq(t, fetchedImgLabel, remoteImgLabel)
						})
					})
				})

				when("there is no remote image", func() {
					when("there is a local image", func() {
						it.Before(func() {
							img, err := local.NewImage(repoName, docker)
							h.AssertNil(t, err)

							h.AssertNil(t, img.Save())
						})

						it.After(func() {
							h.DockerRmi(docker, repoName)
						})

						it("returns the local image", func() {
							_, err := fetcher.Fetch(context.TODO(), repoName, true, pubcfg.PullIfNotPresent)
							h.AssertNil(t, err)
						})
					})

					when("there is no local image", func() {
						it("returns an error", func() {
							_, err := fetcher.Fetch(context.TODO(), repoName, true, pubcfg.PullIfNotPresent)
							h.AssertError(t, err, fmt.Sprintf("image '%s' does not exist on the daemon", repoName))
						})
					})
				})
			})
		})
	})
}
