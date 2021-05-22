## Getting started

* Install Go - you can download it from [golang.org](https://golang.org/) for your architecture or use `brew install go` if you use Homebrew. We recommend 1.13 or higher.
* Clone this repository
* Install `goreleaser` from [goreleser](https://github.com/goreleaser/goreleaser) and `chglog` from [chglog](https://github.com/goreleaser/chglog)
* To compile the CLI, run `make` in the directory
* Use the IDE of preference, we like **Visual Studio Code**

## Overall architecture

Currently, all the components are integrated using [Cobra](https://github.com/spf13/cobra) as the base for our CLI and Civogo which is the Go interface to our API, in addition to that we use other libraries for colours, tables and spinners. We also use `editorconfig` to make all the encoding standard regardless of the IDE used. To see the Civo API document you can go [here](https://www.civo.com/api).


### Library used in the CLI

* [Cobra](https://github.com/spf13/cobra) (to build the CLI)
* [Civogo](https://github.com/civo/civogo) (library for Civo's API)
* [Spinner](https://github.com/briandowns/spinner) (to do the waiting spinners)
* [Color](https://github.com/gookit/color) (to color the messages, also compatible with Windows)
* [Go-homedir](https://github.com/mitchellh/go-homedir) (to get the home of the system user)
* [Tablewriter](https://github.com/olekukonko/tablewriter) (the utility to create and generate text based tables in the terminal)
* [Go-update](https://github.com/tj/go-update) (it is to be able to auto-update the binary)
* [Go-pluralize](https://github.com/alejandrojnm/go-pluralize) (to pluralize any string based on the number of elements of the slice passed to it)

## Release

The release process is currently done manually to be able to have a better control of the changes that are made. However, that does not mean we will not be automating it in the future.

The versioning system used is [Semantic Versioning](https://semver.org/) and [goreleaser](https://github.com/goreleaser/) is used for the process in conjunction with [chglog](https://github.com/goreleaser/chglog) to generate the changelogs. They are only released by the project owners and when a push or merge is made to master, this triggers a GitHub action that generates the docker images, the binary for each platform and the tab for Homebrew.

### Changelog

At the moment these are updated manually, and we are working to integrate it to our CI. To use `chglog` to create a new version we run this command:

```bash
chglog add --version v0.1.1
```

Followed by another command that will generate our changelog:

```bash
chglog format --template repo> CHANGELOG.md
```

### Notification

For each release, an announcement is made on our Slack community channel.
