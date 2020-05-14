# Civo CLI (2020+)

## The plan

The plan is very simple, to write a new version of the Ruby Civo CLI in Golang and deploy it using a simple `curl https://get... | sh` type script or Homebrew for Macs.

We want it to stay interface compatible where possible, with all the aliases currently configured - but want to enhance it with things like custom formatting of output in to JSON or in to custom string formats (e.g. `-o "Hostname|Size"`).

## External libraries

Golang has no shortage of external libraries for various parts of this, but the ones currently planned to be used are:

### Cobra CLI Library

* https://ordina-jworks.github.io/development/2018/10/20/make-your-own-cli-with-golang-and-cobra.html
* https://www.bradcypert.com/testing-a-cobra-cli-in-go/
* https://www.linode.com/docs/development/go/using-cobra/

### Other libraries

* https://github.com/briandowns/spinner
* go get github.com/fatih/color
* https://github.com/olekukonko/tablewriter
* https://github.com/spf13/viper

## Progress

- ✅ ~~Makefile for cross-platform builds~~
- ✅ ~~API Key management~~
- ✅ ~~Regions~~
- ✅ ~~Quotas~~
- ✅ ~~Sizes~~
- ✅ ~~Instances~~
- ✅ ~~Domain names~~
- ✅ ~~Domain records~~
- ✅ ~~Firewalls~~
- ✅ ~~Load balancers~~
- ✅ ~~SSH keys~~
- ✅ ~~Networks~~
- ✅ ~~Snapshots~~
- ✅ ~~Volumes~~
- ✅ ~~Templates~~
- ✅ ~~Kubernetes Clusters~~
- ✅ ~~Kubernetes Applications~~

- `curl | bash` installation mechanism
- Homebrew
