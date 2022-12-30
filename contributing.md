# Contributing guidelines

* [Contributing guidelines](#contributing-guidelines)
    * [Providing Feedback](#providing-feedback)
    * [Opening pull requests (PRs)](#opening-pull-requests-prs)
        * [Choose a good PR title](#choose-a-good-pr-title)
        * [Review your own code](#review-your-own-code)
        * [Do not rebase commits in your branch](#do-not-rebase-commits-in-your-branch)
    * [Contributing to documentation](#contributing-to-documentation)
        * [Ask for help](#ask-for-help)
    * [Prioritizing issues with milestones](#prioritizing-issues-with-milestones)

Before you create a new PR on Ignite CLI, make sure that you read and comply with this document.

Start a new [Discussion](https://github.com/ignite/cli/discussions/new) if you want to propose changes to this document.

To prepare for success, see [Set Up Your Ignite CLI Development Environment](dev-env-setup.md).

To contribute to docs and tutorials, see [Contributing to Ignite CLI Docs](docs/docs/contributing/02-docs.md).

Thank you for your contribution!

## Providing Feedback

* Before you open an issue, do a web search, and check
  for [existing open and closed GitHub Issues](https://github.com/ignite/cli/issues) to see if your question has already
  been asked and answered. If you find a relevant topic, you can comment on that issue.

* To provide feedback or ask a question, create a [GitHub issue](https://github.com/ignite/cli/issues/new/choose). Be
  sure to provide the relevant information, case study, or informative links as suggested by the Pull Request template.

* We recommend using GitHub issues for issues and feedback. However, you can ask quick questions on the **#🛠️
  build-chains** channel in the official [Ignite Discord](https://discord.gg/ignite).

## Opening pull requests (PRs)

Review the issues and discussions before you open a PR.

### Choose a good PR title

Avoid long names in your PR titles. Make sure your title has fewer than 60 characters.

Follow [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0) guidelines and keywords to find the best
title.

Use parentheses to identify the package or feature that you worked on. For example:  `feat(services/chain)`
, `fix(scaffolding)`, `docs(migration)`.

### Review your own code

Make sure that you manually tested the changes you're introducing before creating a PR or pushing another commit.

Monitor your PR to make sure that all CI checks pass and the PR shows **All checks have passed** (the checkmark is
green).

### Do not rebase commits in your branch

Avoid rebasing after you open your PRs to reviews. Instead, add more commits to your PR. It's OK to do force pushes if
the PR is still in draft mode and was never opened to reviews before.

A reviewer likes to see a linear commit history while reviewing. If you tend to force push from an older commit, a
reviewer might lose track in your recent changes and will have to start reviewing from scratch.

Don't worry about adding too many commits. The commits are squashed into a single commit while merging. Your PR title is
used as the commit message.

## Contributing to documentation

When you open a PR for the Ignite CLI codebase, you must also update the relevant documentation. For changes to:

* [Developer Guide](https://docs.ignite.com/guide) tutorials, update content in the `/docs/guide` folder.
* [Knowledge Base](https://docs.ignite.com/kb), update content in the `/docs/kb` folder.
* [Ignite CLI reference](https://docs.ignite.com/cli), navigate to the `./ignite/cmd` package and update the
  documentation of the related command from its `cobra.Command` struct. The CLI docs are automatically generated, so do
  not make changes to  `docs/cli/index.md`.

### Ask for help

If you started a PR but couldn't finish it for whatever reason, don't give up. Instead, just ask for help. Someone else
can take over and assume the ownership.

## Prioritizing issues with milestones

Ignite CLI follows Git Flow for branch
strategy <https://www.atlassian.com/git/tutorials/comparing-workflows/gitflow-workflow>.

* Each Ignite CLI release has a milestone, see <https://github.com/ignite/cli/milestones>.

* Issues in each milestone have a **priority/high**, **priority/medium**, or **priority/low** label.

    * Select issues to work on for the earliest milestone. For example, select to work on an issue labeled as \*\*
      priority/low\*\* in milestone v0.1.0 before you work on an issue labeled as **priority/high** in milestone v0.2.0.

* Milestone **Next** is applied to issues that suggest adding features, docs, and so on.

    * Issues with the **Next** milestone have a higher priority than other **Issues with no milestone** (no milestone is
      assigned).

    * Issues in the **Next** milestone usually have a lower priority than milestones that are associated with a release
      version, like **Milestone v0.1.0**.

* A single project board <https://github.com/ignite/cli/projects/4> shows the issues we are currently working on and
  what issues we plan to work on.

We appreciate your contribution!
