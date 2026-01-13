# Release Process

This document explains how to release a new version of the Civo CLI.

## Overview

The Civo CLI uses [GoReleaser](https://goreleaser.com/) with GitHub Actions to automate the release process. When a Git tag matching the pattern `v*.*.*` is pushed, the release workflow automatically builds and distributes the CLI across multiple platforms and package managers.

## Prerequisites

Before creating a release, ensure:

1. All changes are merged to the `master` branch
2. All tests pass locally: `make test`
3. The CI pipeline is green on the master branch
4. You have push permissions to the repository

## How to Release a New Version

### 1. Prepare the Release

Ensure the master branch is in a releasable state:

```bash
git checkout master
git pull origin master
make test
```

### 2. Create and Push a Git Tag

Create a tag following semantic versioning (`v{MAJOR}.{MINOR}.{PATCH}`):

```bash
git tag v1.2.3
git push origin v1.2.3
```

### 3. Monitor the Release

Once the tag is pushed:

1. Go to the [Actions tab](../../actions) in GitHub
2. Watch the "goreleaser" workflow run
3. The workflow will:
   - Build binaries for all supported platforms
   - Create Docker images and push to Docker Hub
   - Generate Linux packages (DEB/RPM)
   - Create a Homebrew formula PR
   - Create a draft GitHub release

### 4. Publish the Release

1. Go to the [Releases page](../../releases)
2. Find the draft release created by GoReleaser
3. Review the auto-generated changelog
4. Edit release notes if needed
5. Click "Publish release"

### 5. Merge the Homebrew PR

After the release is published:

1. Go to the [civo/homebrew-tools](https://github.com/civo/homebrew-tools) repository
2. Find the PR created by civobot (branch: `civo-{VERSION}`)
3. Review and merge the PR

## What Gets Released

### Binaries

The CLI is built for the following platforms:

| OS | Architectures |
|----|---------------|
| Linux | 386, amd64, arm, arm64 |
| macOS | amd64, arm64 |
| Windows | amd64, arm64 |

Binaries are packaged as:
- `.tar.gz` for Linux and macOS
- `.zip` for Windows

### Docker Images

Docker images are pushed to Docker Hub with two tags:
- `civo/cli:latest` - Always points to the latest release
- `civo/cli:v{VERSION}` - Version-specific tag

The image is based on Alpine Linux and includes `kubectl` pre-installed.

### Linux Packages

- **DEB packages** for Debian/Ubuntu systems
- **RPM packages** for RedHat/CentOS/Fedora systems

### Homebrew

A PR is automatically created to update the Homebrew formula in the [civo/homebrew-tools](https://github.com/civo/homebrew-tools) tap.

### Checksums

A SHA256 checksums file is generated for all artifacts.

## Version Management

The version is derived from Git tags at build time. There is no version file to update manually.

Version information is injected via ldflags during compilation:
- `VersionCli` - The version number (from Git tag)
- `CommitCli` - The Git commit hash
- `DateCli` - The build date

Users can check the installed version with:

```bash
civo version           # Shows: Civo CLI v{VERSION}
civo version -v        # Shows full build info including commit and date
```

## Troubleshooting

### Release Workflow Failed

1. Check the workflow logs in the [Actions tab](../../actions)
2. Common issues:
   - Docker Hub authentication failed - verify `DOCKER_USERNAME` and `DOCKER_PASSWORD` secrets
   - GitHub token expired - regenerate `GORELEASER_GITHUB_TOKEN`
   - GoReleaser config error - run `goreleaser check` locally

### Tag Already Exists

If you need to re-release the same version:

```bash
git tag -d v1.2.3                    # Delete local tag
git push origin :refs/tags/v1.2.3    # Delete remote tag
git tag v1.2.3                       # Recreate tag
git push origin v1.2.3               # Push new tag
```

### Homebrew PR Not Created

Verify the `GORELEASER_GITHUB_TOKEN` has write access to the `civo/homebrew-tools` repository.
