## Getting started

* Install Golang, you can download it at [this address](https://golang.org/) for your architecture, we recommend 1.13 or higher
* Clone the repository
* Install `goreleaser` and` chglog`, [goreleser](https://github.com/goreleaser/goreleaser) and [chglog](https://github.com/goreleaser/chglog)
* To compile the changes  run `make` in the directory
* Use the IDE of preference, we like **visual studio code**

## Overall architecture

Currently all the components are integrated using cobra as the base for our cli and civogo which is the interface to our api, in addition to that we use other libraries for colors, tables and spinners. We also use `editorconfig` to make all the encoding standard regardless of the IDE used. To see the help of the civo api you can go [here](https://www.civo.com/api)

### Library used in the cli

* [Cobra](https://github.com/spf13/cobra) (to build the cli)
* [Civogo](https://github.com/civo/civogo) (library for Civo's api)
* [Spinner](github.com/briandowns/spinner) (to do the waiting spinners)
* [Color](https://github.com/gookit/color) (to color the messages, it is also compatible with Windows)
* [Go-homedir](https://github.com/mitchellh/go-homedir) (to get the home of the system user)
* [Tablewriter](https://github.com/olekukonko/tablewriter) (the utility to make the output tables to the terminal)
* [Go-update](https://github.com/tj/go-update) (it is to be able to auto-update the binary)

## Release

The release process is currently done manual to be able to have a better control of the changes that are made, that does not mean we will not be automating it in the future. The versioning system used is [Semantic Versioning](https://semver.org/) and [goreleaser](https://github.com/goreleaser/) is used for the process in conjunction with [chglog](https://github.com/goreleaser/chglog) to generate the change logs, they are only released by the project owners and when a push or merge is made to master, this triggers a GitHub action that generates, the docker images, the binary for each platform and the tab for homebrew.

### Changelog

At the moment it is something that is done manually and we are working to integrate it to our CI, how does `chglog` work, to create a new version we run this command:

```bash
chglog add --version v0.1.1
```

and after that we run another one that will generate our changelog:

```bash
chglog format --template repo> CHANGELOG.md
```

### Notification

In each release, an announcement is made on our slack community channel
