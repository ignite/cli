# Contributing guidelines

If you're willing to create a new PR on Starport, make sure that you read and comply with this document.

Start a new [Discussion](https://github.com/tendermint/starport/discussions/new) if you want to propose changes to this document.

Thank you for your contribution!

## Providing Feedback

* Before you open an issue, do a web search, and check for [existing open and closed GitHub Issues](https://github.com/tendermint/starport/issues) to see if your question has already been asked and answered. If you find a relevant topic, you can comment on that issue.
* To provide feedback or ask a question, create a [GitHub issue](https://github.com/tendermint/starport/issues/new/choose). Be sure to provide the relevant information, case study, or informative links as suggested by the Pull Request template.
* We recommend using GitHub issues for issues and feedback. However, you can ask quick questions on the [#🔨cosmos-sdk-starport](https://discord.com/channels/669268347736686612/737461683588431924) channel in Discord.

## Opening pull requests

### Choose a good PR title

Avoid long names in your PR titles. Make sure your title has fewer than 60 characters.

Follow [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) guidelines and keywords to find the best title.

Use parentheses to identify the package or feature that you worked on. For example:  `feat(services/chain)`, `fix(scaffolding)`.

### Review your own code

Make sure that you manually tested the changes you're introducing before creating a PR or pushing another commit.

Monitor your PR to make sure that all CI checks pass and the PR shows **All checks have passed** (the checkmark is green).

### Avoid rebasing commits in your branch

Avoid rebasing after you open your PRs to reviews. Instead, add more commits to your PR. It's OK to do force pushes if the PR is still in draft mode and was never opened to reviews before.

A reviewer likes to see a linear commit history while reviewing. If you tend to force push from an older commit, reviewer might lose track in your recent changes and will have to start reviewing from scratch.

Don't worry about adding too many commits. The commits are squashed into a single commit while merging. Your PR title is used as the commit message.

### Ask for help

If you started a PR but couldn't finish it for whatever reason, don't give up. Instead, just ask for help. Someone else can take over and assume the ownership.

We appreciate every bit of your work!

## Contributing to documentation

[CLI reference](https://github.com/tendermint/starport/blob/f668bba58c04318f98db8cac0c9e154fa7e7ea34/docs/cli/index.md) is generated automatically with a [GitHub action](https://github.com/tendermint/starport/blob/f668bba58c04318f98db8cac0c9e154fa7e7ea34/.github/workflows/gen-docs-cli.yml). Instead of editing `docs/cli/index.md` directly (which will be overwritten by the GH action), please, make changes to the relevant parts of the source code that describe CLI commands.
