# Release Process

This document outlines the release process for Ignite, ensuring consistency, quality, and clear communication with users. Ignite follows [semantic versioning](https://semver.org/) to signal the stability and compatibility of each release.

## Development Branch

The `main` branch serves as the development branch for Ignite. All new features, bug fixes, and updates are merged into this branch. The `main` branch is typically updated regularly, depending on development activity.

## Backporting Features & Bug Fixes

Features and bug fixes ready for release are backported from `main` to the release branch. This process is automated using [Mergify](https://mergify.com/), a CI/CD tool. By adding the `backport release/x.y.z` label to a PR, Mergify will automatically backport the PR to the release branch when it is merged into `main`.

## Changelog

Before any release, the changelog must be up-to-date. It lists all changes made to Ignite since the last release and must be carefully reviewed to ensure accuracy.

## Release Cadence: Alpha → Beta → RC → Full Release

To accommodate Ignite’s dependency on the [Cosmos SDK](https://github.com/cosmos/cosmos-sdk) and the addition of new features, a structured release cadence is used that includes Alpha, Beta, and Release Candidate (RC) stages before the Full Release. This ensures features are thoroughly tested, dependencies are compatible, and disruptions are minimized.

### Major & Minor Releases

For **major** releases, the Alpha → Beta → RC progression is **required**. For **minor** releases, these stages are **optional** and used at the discretion of recent updates and their complexity.

These stages help ensure stability, especially given Ignite's integration with the Cosmos SDK and potential for breaking changes. This process allows developers to prepare for compatibility changes:

- **Alpha Releases**: Early, incomplete versions with new features that may be unstable, shared internally or with select testers.
- **Beta Releases**: Feature-complete versions shared with the community for feedback and testing. These releases are more stable than Alpha but may still contain bugs.
- **RC (Release Candidate)**: A near-final version with all intended features and fixes. Multiple RCs may be issued for **MAJOR releases** (e.g., `v30.0.0-rc1`, `v30.0.0-rc2`, etc.) as issues are identified and resolved. The final RC becomes the Full Release if no critical issues remain.

### Patch Releases

**Patch releases** do not go through Alpha, Beta, and RC stages. They are fast-tracked for release after internal testing, as they address specific bug fixes or security vulnerabilities. Patch releases should remain backward-compatible and thoroughly tested for regressions.

## Managing SDK Dependencies & Compatibility

Given Ignite’s reliance on the Cosmos SDK, ensuring compatibility between Ignite releases and the SDK is crucial. When upgrading the SDK version, a transition period may be needed to allow users time to adapt.

- **Backward Compatibility**: Ignite strives to maintain backward compatibility for `chain` and `app` commands between major releases. This allows users to upgrade Ignite without immediate refactoring.
- **Breaking Changes**: If breaking changes are introduced (such as SDK upgrades or major feature revisions), transition periods will be defined and communicated in release notes starting from Alpha versions.
- **Transition Periods**: When significant changes impact downstream applications, transition periods give users time to test and adapt their applications before the final release.

## Release Branches

Releases are tagged from dedicated release branches, named after the release version (e.g., `release/v28.x.y` or `release/v30.x.y`). These branches are created from `main` and contain all changes intended for the release.

- **Alpha, Beta, and RC Branches**: Pre-release branches are named accordingly, such as `release/v28.x.y-alpha`, `release/v28.x.y-beta`, or `release/v28.x.y-rc1`, `release/v28.x.y-rc2`, etc.

## Release Preparation & Testing

The preparation and testing phases vary depending on the type of release:

### Major Releases

- **Freeze `main`**: No new features are merged into `main` during final preparation.
- **Create the release branch**: A new branch (e.g., `release/v30.x.y`) is created from `main`.
- **Backport**: Ensure that all desired features and fixes are backported to the release branch.
- **Testing**: Run unit, integration, and manual tests.
- **Changelog**: Finalize the changelog.

### Minor & Patch Releases

- **Backport**: Ensure that all necessary changes are backported to the release branch.
- **Testing**: Conduct unit, integration, and manual tests.
- **Changelog**: Finalize the changelog.

## Release Publication

When testing is complete, the release is published to the [releases page](https://github.com/ignite/cli/releases) on GitHub. This includes tagging the release branch with the version number and publishing a release announcement with the changelog.

```sh
git checkout release/v28.x.y
git tag v28.x.y -m "Release Ignite v28.x.y"
```

For Alpha, Beta, and RC releases, use the appropriate tags (e.g., `v28.x.y-alpha`, `v28.x.y-beta`, `v28.x.y-rc1`, etc.).

## Post-Release Activities

After a release, monitor feedback and bug reports. These will inform subsequent patch releases or feature additions.

Following a **MAJOR** release, the `main` branch must be updated to the next **MAJOR** version. This includes updating the `go.mod` file and any other version number references in the codebase.

## Maintenance Policy

Only the latest released version of Ignite is actively maintained for new features and fixes. Older versions may continue to function but will not receive updates, ensuring stability and security for users.

Users are encouraged to upgrade to the latest release to benefit from the newest features and fixes.

Ignite ensures compatibility for `chain` and `app` commands between **MAJOR** releases, but other commands may change and may require users to upgrade their codebase to match the Cosmos SDK version used by Ignite.
