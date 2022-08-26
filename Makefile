GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=civo
BINARY_MAC=$(BINARY_NAME)_mac
BINARY_LINUX=$(BINARY_NAME)_linux
BINARY_WINDOWS=$(BINARY_NAME)_windows
OS=$(shell uname)

all: build
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -rf dest/

build: buildmac buildlinux buildwindows
	@rm -f civo
	# Assuming either Mac or Linux at this point...too many options for running Make on Windows
ifeq ($(OS),Darwin)
	ln -s dest/$(BINARY_MAC) civo
else
	ln -s dest/$(BINARY_LINUX) civo
endif

buildmac: $(BINARY_MAC)
$(BINARY_MAC): buildprep
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o dest/$(BINARY_MAC) -ldflags "-s -X github.com/civo/cli/common.VersionCli=$(VERSION_CLI) -X github.com/civo/cli/common.CommitCli=$(COMMIT_CLI) -X github.com/civo/cli/common.DateCli=$(shell date +%FT%T%Z)" -v

buildlinux: $(BINARY_LINUX)
$(BINARY_LINUX): buildprep
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o dest/$(BINARY_LINUX) -ldflags "-s -X github.com/civo/cli/common.VersionCli=$(VERSION_CLI) -X github.com/civo/cli/common.CommitCli=$(COMMIT_CLI) -X github.com/civo/cli/common.DateCli=$(shell date +%FT%T%Z)" -v

buildwindows: $(BINARY_WINDOWS)
$(BINARY_WINDOWS): buildprep
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o dest/$(BINARY_WINDOWS) -ldflags "-s -X github.com/civo/cli/common.VersionCli=$(VERSION_CLI) -X github.com/civo/cli/common.CommitCli=$(COMMIT_CLI) -X github.com/civo/cli/common.DateCli=$(shell date +%FT%T%Z)" -v

buildprep:
	git fetch --tags -f
	mkdir -p dest
	# $(eval VERSION_CLI=$(shell git describe --tags | cut -d "v" -f 2 | cut -d "-" -f 1))
	$(eval VERSION_CLI=$(shell git describe --tags | cut -d "v" -f 2 | cut -d "-" -f 1))
	$(eval COMMIT_CLI=$(shell git log --format="%H" -n 1))

release: build
	# extract version from app
	VERSION=`dest/civo_mac version -q`
	# Set a github token with git config --global github.token "....."
	# $ ghr \
  #   -t TOKEN \        # Set Github API Token
  #   -u USERNAME \     # Set Github username
  #   -r REPO \         # Set repository name
  #   -c COMMIT \       # Set target commitish, branch or commit SHA
  #   -n TITLE \        # Set release title
  #   -b BODY \         # Set text describing the contents of the release
  #   -p NUM \          # Set amount of parallelism (Default is number of CPU)
  #   -delete \         # Delete release and its git tag in advance if it exists (same as -recreate)
  #   -replace          # Replace artifacts if it is already uploaded
  #   -draft \          # Release as draft (Unpublish)
  #   -soft \           # Stop uploading if the same tag already exists
  #   -prerelease \     # Create prerelease
  #   TAG PATH
	# ghr $VERSION pkg/
	# https://docs.brew.sh/How-to-Create-and-Maintain-a-Tap
