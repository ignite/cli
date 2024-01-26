# Release Process

This document outlines the release process for Ignite. It ensures consistency, quality, and clear communication with users.
Ignite uses [semantic versioning](https://semver.org/) to indicate the stability and compatibility of releases.

## Release Branches

Releases are tagged from release branches. The release branch is typically named after the release version, such as `release/v28.x.y` or `release/v30.x.y`. The release branch is created from the `main` branch and contains all the changes that will be included in the release.

## Development Branch

The `main` branch is the development branch for Ignite. All new features and bug fixes are developed on this branch. The `main` branch is typically updated daily or weekly, depending on the amount of development activity.

## Backporting Features and Bug Fixes

Features or bug fixes that are ready to be included in a release must be backported from the `main` branch to the release branch. This is done using Mergify, a CI/CD tool that automates the process of backporting changes. Add the `backport release/x.y.z` label to a pull request (PR) to indicate that it should be backported to the release branch. Mergify will automatically backport the PR when it is merged into the `main` branch.

## All PRs Target Main

All PRs should target the `main` branch. This ensures that changes are always integrated into the development branch and are ready for backporting.

## Changelog

The changelog lists all of the changes that have been made to Ignite since the last release. The changelog must be up-to-date before a release is made.

## Release Preparation and Testing

Before a release is made, it is important to prepare the release branches and perform thorough testing.
The process differs depending on whether the release is a `MAJOR`, `MINOR`, or `PATCH` release.

### Major Release

* Freeze the `main` branch.
* Create a release branch from the `main` branch.
  * Update `.github/mergify.yml` to allow backporting changes to the release branch.
  * Possibly update CI/CD configuration to support the new release.
* Running all unit tests, integration tests, and manual scenarios.
* Ensuring that the release branch is still compatible with all supported environments.
* Prepare the changelog.

### Minor and Patch Releases

* Verifying that all wanted changes have been backported to the release branch.
* Running all unit tests, integration tests, and manual scenarios.
* Ensuring that the release branch is still compatible with all supported environments.
* Prepare the changelog.

## Release Publication

Once the release is ready, it can be published to the [releases page](https://github.com/ignite/cli/releases) on GitHub. This involves tagging the release branch with the release version and creating a release announcement. The release anouncement should contain the changelog of the release.

```bash
git checkout release/v28.x.y
git tag v28.x.y -m "Release Ignite v28.x.y"
```

## Post-Release Activities

After a release has been made, it is important to monitor feedback and bug reports to inform subsequent releases.

This includes updating the `main` branch changelog with the new release.

In case of a new `MAJOR` release, the release author must also update the `main` branch to the next `MAJOR` version number: rename the `go.mod` and all references to the version number in the codebase.

## Maintenance Policy

Only the latest released version of Ignite is maintained for new features and bug fixes. Older versions may continue to function, but they will not receive any updates. This ensures that Ignite remains stable and reliable for all users.

Users are encouraged to upgrade to the latest release as soon as possible to benefit from the latest features and security updates.
Ignite ensures compatibility for the `chain` and `app` commands between major releases.
Other commands may change between major releases and may require the user to upgrade their codebase to the Cosmos SDK version Ignite is using.
