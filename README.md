# pack - Buildpack CLI

[![Build results](https://github.com/YousefHaggyHeroku/pack/workflows/build/badge.svg)](https://github.com/YousefHaggyHeroku/pack/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/YousefHaggyHeroku/pack)](https://goreportcard.com/report/github.com/YousefHaggyHeroku/pack)
[![codecov](https://codecov.io/gh/YousefHaggyHeroku/pack/branch/main/graph/badge.svg)](https://codecov.io/gh/YousefHaggyHeroku/pack)
[![GoDoc](https://godoc.org/github.com/YousefHaggyHeroku/pack?status.svg)](https://godoc.org/github.com/YousefHaggyHeroku/pack)
[![GitHub license](https://img.shields.io/github/license/YousefHaggyHeroku/pack)](https://github.com/YousefHaggyHeroku/pack/blob/main/LICENSE)
[![Slack](https://slack.buildpacks.io/badge.svg)](https://slack.buildpacks.io/)

`pack` makes it easy for...
- [**App Developers**][app-dev] to use buildpacks to convert code into runnable images.
- [**Buildpack Authors**][bp-author] to develop and package buildpacks for distribution.
- [**Operators**][operator] to package buildpacks for distribution and maintain applications.

## Usage

<img src="resources/pack-build.gif" width="600px" />

## Getting Started
Get started by running through our tutorial: [An App’s Brief Journey from Source to Image][getting-started]

## Contributing
- [CONTRIBUTING](CONTRIBUTING.md) - Information on how to contribute, including the pull request process.
- [DEVELOPMENT](DEVELOPMENT.md) - Further detail to help you during the development process.
- [RELEASE](RELEASE.md) - Further details about our release process.

## Documentation
Check out the command line documentation [here][pack-docs]

## Specifications
`pack` is a CLI implementation of the [Platform Interface Specification][platform-spec] for [Cloud Native Buildpacks][buildpacks.io].

To learn more about the details, check out the [specs repository][specs].

[app-dev]: https://buildpacks.io/docs/app-developer-guide/
[bp-author]: https://buildpacks.io/docs/buildpack-author-guide/
[operator]: https://buildpacks.io/docs/operator-guide/
[buildpacks.io]: https://buildpacks.io/
[install-pack]: https://buildpacks.io/docs/install-pack/
[getting-started]: https://buildpacks.io/docs/app-journey
[specs]: https://github.com/buildpacks/spec/
[platform-spec]: https://github.com/buildpacks/spec/blob/main/platform.md
[pack-docs]: https://buildpacks.io/docs/tools/pack/cli/pack/
